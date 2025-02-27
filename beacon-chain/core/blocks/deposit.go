package blocks

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/helpers"
	iface "github.com/prysmaticlabs/prysm/beacon-chain/state/interface"
	pb "github.com/prysmaticlabs/prysm/proto/beacon/p2p/v1"
	"github.com/prysmaticlabs/prysm/shared/bls"
	"github.com/prysmaticlabs/prysm/shared/bytesutil"
	"github.com/prysmaticlabs/prysm/shared/depositutil"
	"github.com/prysmaticlabs/prysm/shared/mathutil"
	"github.com/prysmaticlabs/prysm/shared/params"
	"github.com/prysmaticlabs/prysm/shared/trieutil"
)

// ProcessPreGenesisDeposits processes a deposit for the beacon state before chainstart.
func ProcessPreGenesisDeposits(
	ctx context.Context,
	beaconState iface.BeaconState,
	deposits []*ethpb.Deposit,
) (iface.BeaconState, error) {
	var err error
	beaconState, err = ProcessDeposits(ctx, beaconState, &ethpb.SignedBeaconBlock{
		Block: &ethpb.BeaconBlock{Body: &ethpb.BeaconBlockBody{Deposits: deposits}}})
	if err != nil {
		return nil, errors.Wrap(err, "could not process deposit")
	}
	for _, deposit := range deposits {
		pubkey := deposit.Data.PublicKey
		index, ok := beaconState.ValidatorIndexByPubkey(bytesutil.ToBytes48(pubkey))
		// In the event of the pubkey not existing, we continue processing the other
		// deposits.
		if !ok {
			continue
		}
		balance, err := beaconState.BalanceAtIndex(index)
		if err != nil {
			return nil, err
		}
		validator, err := beaconState.ValidatorAtIndex(index)
		if err != nil {
			return nil, err
		}
		validator.EffectiveBalance = mathutil.Min(balance-balance%params.BeaconConfig().EffectiveBalanceIncrement, params.BeaconConfig().MaxEffectiveBalance)
		if validator.EffectiveBalance ==
			params.BeaconConfig().MaxEffectiveBalance {
			validator.ActivationEligibilityEpoch = 0
			validator.ActivationEpoch = 0
		}
		if err := beaconState.UpdateValidatorAtIndex(index, validator); err != nil {
			return nil, err
		}
	}
	return beaconState, nil
}

// ProcessDeposits is one of the operations performed on each processed
// beacon block to verify queued validators from the Ethereum 1.0 Deposit Contract
// into the beacon chain.
//
// Spec pseudocode definition:
//   For each deposit in block.body.deposits:
//     process_deposit(state, deposit)
func ProcessDeposits(
	ctx context.Context,
	beaconState iface.BeaconState,
	b *ethpb.SignedBeaconBlock,
) (iface.BeaconState, error) {
	if err := helpers.VerifyNilBeaconBlock(b); err != nil {
		return nil, err
	}

	deposits := b.Block.Body.Deposits
	var err error
	domain, err := helpers.ComputeDomain(params.BeaconConfig().DomainDeposit, nil, nil)
	if err != nil {
		return nil, err
	}

	// Attempt to verify all deposit signatures at once, if this fails then fall back to processing
	// individual deposits with signature verification enabled.
	var verifySignature bool
	if err := verifyDepositDataWithDomain(ctx, deposits, domain); err != nil {
		log.WithError(err).Debug("Failed to verify deposit data, verifying signatures individually")
		verifySignature = true
	}

	for _, deposit := range deposits {
		if deposit == nil || deposit.Data == nil {
			return nil, errors.New("got a nil deposit in block")
		}
		beaconState, err = ProcessDeposit(beaconState, deposit, verifySignature)
		if err != nil {
			return nil, errors.Wrapf(err, "could not process deposit from %#x", bytesutil.Trunc(deposit.Data.PublicKey))
		}
	}
	return beaconState, nil
}

// ProcessDeposit takes in a deposit object and inserts it
// into the registry as a new validator or balance change.
//
// Spec pseudocode definition:
// def process_deposit(state: BeaconState, deposit: Deposit) -> None:
//    # Verify the Merkle branch
//    assert is_valid_merkle_branch(
//        leaf=hash_tree_root(deposit.data),
//        branch=deposit.proof,
//        depth=DEPOSIT_CONTRACT_TREE_DEPTH + 1,  # Add 1 for the List length mix-in
//        index=state.eth1_deposit_index,
//        root=state.eth1_data.deposit_root,
//    )
//
//    # Deposits must be processed in order
//    state.eth1_deposit_index += 1
//
//    pubkey = deposit.data.pubkey
//    amount = deposit.data.amount
//    validator_pubkeys = [v.pubkey for v in state.validators]
//    if pubkey not in validator_pubkeys:
//        # Verify the deposit signature (proof of possession) which is not checked by the deposit contract
//        deposit_message = DepositMessage(
//            pubkey=deposit.data.pubkey,
//            withdrawal_credentials=deposit.data.withdrawal_credentials,
//            amount=deposit.data.amount,
//        )
//        domain = compute_domain(DOMAIN_DEPOSIT)  # Fork-agnostic domain since deposits are valid across forks
//        signing_root = compute_signing_root(deposit_message, domain)
//        if not bls.Verify(pubkey, signing_root, deposit.data.signature):
//            return
//
//        # Add validator and balance entries
//        state.validators.append(get_validator_from_deposit(state, deposit))
//        state.balances.append(amount)
//    else:
//        # Increase balance by deposit amount
//        index = ValidatorIndex(validator_pubkeys.index(pubkey))
//        increase_balance(state, index, amount)
func ProcessDeposit(beaconState iface.BeaconState, deposit *ethpb.Deposit, verifySignature bool) (iface.BeaconState, error) {
	if err := verifyDeposit(beaconState, deposit); err != nil {
		if deposit == nil || deposit.Data == nil {
			return nil, err
		}
		return nil, errors.Wrapf(err, "could not verify deposit from %#x", bytesutil.Trunc(deposit.Data.PublicKey))
	}
	if err := beaconState.SetEth1DepositIndex(beaconState.Eth1DepositIndex() + 1); err != nil {
		return nil, err
	}
	pubKey := deposit.Data.PublicKey
	amount := deposit.Data.Amount
	index, ok := beaconState.ValidatorIndexByPubkey(bytesutil.ToBytes48(pubKey))
	if !ok {
		if verifySignature {
			domain, err := helpers.ComputeDomain(params.BeaconConfig().DomainDeposit, nil, nil)
			if err != nil {
				return nil, err
			}
			if err := verifyDepositDataSigningRoot(deposit.Data, domain); err != nil {
				// Ignore this error as in the spec pseudo code.
				log.Debugf("Skipping deposit: could not verify deposit data signature: %v", err)
				return beaconState, nil
			}
		}

		effectiveBalance := amount - (amount % params.BeaconConfig().EffectiveBalanceIncrement)
		if params.BeaconConfig().MaxEffectiveBalance < effectiveBalance {
			effectiveBalance = params.BeaconConfig().MaxEffectiveBalance
		}
		if err := beaconState.AppendValidator(&ethpb.Validator{
			PublicKey:                  pubKey,
			WithdrawalCredentials:      deposit.Data.WithdrawalCredentials,
			ActivationEligibilityEpoch: params.BeaconConfig().FarFutureEpoch,
			ActivationEpoch:            params.BeaconConfig().FarFutureEpoch,
			ExitEpoch:                  params.BeaconConfig().FarFutureEpoch,
			WithdrawableEpoch:          params.BeaconConfig().FarFutureEpoch,
			EffectiveBalance:           effectiveBalance,
		}); err != nil {
			return nil, err
		}
		if err := beaconState.AppendBalance(amount); err != nil {
			return nil, err
		}
	} else if err := helpers.IncreaseBalance(beaconState, index, amount); err != nil {
		return nil, err
	}

	return beaconState, nil
}

func verifyDeposit(beaconState iface.ReadOnlyBeaconState, deposit *ethpb.Deposit) error {
	// Verify Merkle proof of deposit and deposit trie root.
	if deposit == nil || deposit.Data == nil {
		return errors.New("received nil deposit or nil deposit data")
	}
	eth1Data := beaconState.Eth1Data()
	if eth1Data == nil {
		return errors.New("received nil eth1data in the beacon state")
	}

	receiptRoot := eth1Data.DepositRoot
	leaf, err := deposit.Data.HashTreeRoot()
	if err != nil {
		return errors.Wrap(err, "could not tree hash deposit data")
	}
	if ok := trieutil.VerifyMerkleBranch(
		receiptRoot,
		leaf[:],
		int(beaconState.Eth1DepositIndex()),
		deposit.Proof,
		params.BeaconConfig().DepositContractTreeDepth,
	); !ok {
		return fmt.Errorf(
			"deposit merkle branch of deposit root did not verify for root: %#x",
			receiptRoot,
		)
	}

	return nil
}

func verifyDepositDataSigningRoot(obj *ethpb.Deposit_Data, domain []byte) error {
	return depositutil.VerifyDepositSignature(obj, domain)
}

func verifyDepositDataWithDomain(ctx context.Context, deps []*ethpb.Deposit, domain []byte) error {
	if len(deps) == 0 {
		return nil
	}
	pks := make([]bls.PublicKey, len(deps))
	sigs := make([][]byte, len(deps))
	msgs := make([][32]byte, len(deps))
	for i, dep := range deps {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if dep == nil || dep.Data == nil {
			return errors.New("nil deposit")
		}
		dpk, err := bls.PublicKeyFromBytes(dep.Data.PublicKey)
		if err != nil {
			return err
		}
		pks[i] = dpk
		sigs[i] = dep.Data.Signature
		depositMessage := &pb.DepositMessage{
			PublicKey:             dep.Data.PublicKey,
			WithdrawalCredentials: dep.Data.WithdrawalCredentials,
			Amount:                dep.Data.Amount,
		}
		sr, err := helpers.ComputeSigningRoot(depositMessage, domain)
		if err != nil {
			return err
		}
		msgs[i] = sr
	}
	verify, err := bls.VerifyMultipleSignatures(sigs, msgs, pks)
	if err != nil {
		return errors.Errorf("could not verify multiple signatures: %v", err)
	}
	if !verify {
		return errors.New("one or more deposit signatures did not verify")
	}
	return nil
}
