package stateV0

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/pkg/errors"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	"github.com/prysmaticlabs/prysm/beacon-chain/state/stateutil"
	pb "github.com/prysmaticlabs/prysm/proto/beacon/p2p/v1"
	"github.com/prysmaticlabs/prysm/shared/bytesutil"
	"github.com/prysmaticlabs/prysm/shared/hashutil"
)

// FieldTrie is the representation of the representative
// trie of the particular field.
type FieldTrie struct {
	*sync.RWMutex
	reference   *stateutil.Reference
	fieldLayers [][]*[32]byte
	field       fieldIndex
}

// NewFieldTrie is the constructor for the field trie data structure. It creates the corresponding
// trie according to the given parameters. Depending on whether the field is a basic/composite array
// which is either fixed/variable length, it will appropriately determine the trie.
func NewFieldTrie(field fieldIndex, elements interface{}, length uint64) (*FieldTrie, error) {
	if elements == nil {
		return &FieldTrie{
			field:     field,
			reference: stateutil.NewRef(1),
			RWMutex:   new(sync.RWMutex),
		}, nil
	}
	datType, ok := fieldMap[field]
	if !ok {
		return nil, errors.Errorf("unrecognized field in trie")
	}
	fieldRoots, err := fieldConverters(field, []uint64{}, elements, true)
	if err != nil {
		return nil, err
	}
	switch datType {
	case basicArray:
		return &FieldTrie{
			fieldLayers: stateutil.ReturnTrieLayer(fieldRoots, length),
			field:       field,
			reference:   stateutil.NewRef(1),
			RWMutex:     new(sync.RWMutex),
		}, nil
	case compositeArray:
		return &FieldTrie{
			fieldLayers: stateutil.ReturnTrieLayerVariable(fieldRoots, length),
			field:       field,
			reference:   stateutil.NewRef(1),
			RWMutex:     new(sync.RWMutex),
		}, nil
	default:
		return nil, errors.Errorf("unrecognized data type in field map: %v", reflect.TypeOf(datType).Name())
	}

}

// RecomputeTrie rebuilds the affected branches in the trie according to the provided
// changed indices and elements. This recomputes the trie according to the particular
// field the trie is based on.
func (f *FieldTrie) RecomputeTrie(indices []uint64, elements interface{}) ([32]byte, error) {
	f.Lock()
	defer f.Unlock()
	var fieldRoot [32]byte
	if len(indices) == 0 {
		return f.TrieRoot()
	}
	datType, ok := fieldMap[f.field]
	if !ok {
		return [32]byte{}, errors.Errorf("unrecognized field in trie")
	}
	fieldRoots, err := fieldConverters(f.field, indices, elements, false)
	if err != nil {
		return [32]byte{}, err
	}
	switch datType {
	case basicArray:
		fieldRoot, f.fieldLayers, err = stateutil.RecomputeFromLayer(fieldRoots, indices, f.fieldLayers)
		if err != nil {
			return [32]byte{}, err
		}
		return fieldRoot, nil
	case compositeArray:
		fieldRoot, f.fieldLayers, err = stateutil.RecomputeFromLayerVariable(fieldRoots, indices, f.fieldLayers)
		if err != nil {
			return [32]byte{}, err
		}
		return stateutil.AddInMixin(fieldRoot, uint64(len(f.fieldLayers[0])))
	default:
		return [32]byte{}, errors.Errorf("unrecognized data type in field map: %v", reflect.TypeOf(datType).Name())
	}

}

// CopyTrie copies the references to the elements the trie
// is built on.
func (f *FieldTrie) CopyTrie() *FieldTrie {
	if f.fieldLayers == nil {
		return &FieldTrie{
			field:     f.field,
			reference: stateutil.NewRef(1),
			RWMutex:   new(sync.RWMutex),
		}
	}
	dstFieldTrie := make([][]*[32]byte, len(f.fieldLayers))
	for i, layer := range f.fieldLayers {
		dstFieldTrie[i] = make([]*[32]byte, len(layer))
		copy(dstFieldTrie[i], layer)
	}
	return &FieldTrie{
		fieldLayers: dstFieldTrie,
		field:       f.field,
		reference:   stateutil.NewRef(1),
		RWMutex:     new(sync.RWMutex),
	}
}

// TrieRoot returns the corresponding root of the trie.
func (f *FieldTrie) TrieRoot() ([32]byte, error) {
	datType, ok := fieldMap[f.field]
	if !ok {
		return [32]byte{}, errors.Errorf("unrecognized field in trie")
	}
	switch datType {
	case basicArray:
		return *f.fieldLayers[len(f.fieldLayers)-1][0], nil
	case compositeArray:
		trieRoot := *f.fieldLayers[len(f.fieldLayers)-1][0]
		return stateutil.AddInMixin(trieRoot, uint64(len(f.fieldLayers[0])))
	default:
		return [32]byte{}, errors.Errorf("unrecognized data type in field map: %v", reflect.TypeOf(datType).Name())
	}
}

// this converts the corresponding field and the provided elements to the appropriate roots.
func fieldConverters(field fieldIndex, indices []uint64, elements interface{}, convertAll bool) ([][32]byte, error) {
	switch field {
	case blockRoots, stateRoots, randaoMixes:
		val, ok := elements.([][]byte)
		if !ok {
			return nil, errors.Errorf("Wanted type of %v but got %v",
				reflect.TypeOf([][]byte{}).Name(), reflect.TypeOf(elements).Name())
		}
		return handleByteArrays(val, indices, convertAll)
	case eth1DataVotes:
		val, ok := elements.([]*ethpb.Eth1Data)
		if !ok {
			return nil, errors.Errorf("Wanted type of %v but got %v",
				reflect.TypeOf([]*ethpb.Eth1Data{}).Name(), reflect.TypeOf(elements).Name())
		}
		return handleEth1DataSlice(val, indices, convertAll)
	case validators:
		val, ok := elements.([]*ethpb.Validator)
		if !ok {
			return nil, errors.Errorf("Wanted type of %v but got %v",
				reflect.TypeOf([]*ethpb.Validator{}).Name(), reflect.TypeOf(elements).Name())
		}
		return handleValidatorSlice(val, indices, convertAll)
	case previousEpochAttestations, currentEpochAttestations:
		val, ok := elements.([]*pb.PendingAttestation)
		if !ok {
			return nil, errors.Errorf("Wanted type of %v but got %v",
				reflect.TypeOf([]*pb.PendingAttestation{}).Name(), reflect.TypeOf(elements).Name())
		}
		return handlePendingAttestation(val, indices, convertAll)
	default:
		return [][32]byte{}, errors.Errorf("got unsupported type of %v", reflect.TypeOf(elements).Name())
	}
}

func handleByteArrays(val [][]byte, indices []uint64, convertAll bool) ([][32]byte, error) {
	length := len(indices)
	if convertAll {
		length = len(val)
	}
	roots := make([][32]byte, 0, length)
	rootCreator := func(input []byte) {
		newRoot := bytesutil.ToBytes32(input)
		roots = append(roots, newRoot)
	}
	if convertAll {
		for i := range val {
			rootCreator(val[i])
		}
		return roots, nil
	}
	if len(val) > 0 {
		for _, idx := range indices {
			if idx > uint64(len(val))-1 {
				return nil, fmt.Errorf("index %d greater than number of byte arrays %d", idx, len(val))
			}
			rootCreator(val[idx])
		}
	}
	return roots, nil
}

func handleEth1DataSlice(val []*ethpb.Eth1Data, indices []uint64, convertAll bool) ([][32]byte, error) {
	length := len(indices)
	if convertAll {
		length = len(val)
	}
	roots := make([][32]byte, 0, length)
	hasher := hashutil.CustomSHA256Hasher()
	rootCreator := func(input *ethpb.Eth1Data) error {
		newRoot, err := eth1Root(hasher, input)
		if err != nil {
			return err
		}
		roots = append(roots, newRoot)
		return nil
	}
	if convertAll {
		for i := range val {
			err := rootCreator(val[i])
			if err != nil {
				return nil, err
			}
		}
		return roots, nil
	}
	if len(val) > 0 {
		for _, idx := range indices {
			if idx > uint64(len(val))-1 {
				return nil, fmt.Errorf("index %d greater than number of items in eth1 data slice %d", idx, len(val))
			}
			err := rootCreator(val[idx])
			if err != nil {
				return nil, err
			}
		}
	}
	return roots, nil
}

func handleValidatorSlice(val []*ethpb.Validator, indices []uint64, convertAll bool) ([][32]byte, error) {
	length := len(indices)
	if convertAll {
		length = len(val)
	}
	roots := make([][32]byte, 0, length)
	hasher := hashutil.CustomSHA256Hasher()
	rootCreator := func(input *ethpb.Validator) error {
		newRoot, err := stateutil.ValidatorRootWithHasher(hasher, input)
		if err != nil {
			return err
		}
		roots = append(roots, newRoot)
		return nil
	}
	if convertAll {
		for i := range val {
			err := rootCreator(val[i])
			if err != nil {
				return nil, err
			}
		}
		return roots, nil
	}
	if len(val) > 0 {
		for _, idx := range indices {
			if idx > uint64(len(val))-1 {
				return nil, fmt.Errorf("index %d greater than number of validators %d", idx, len(val))
			}
			err := rootCreator(val[idx])
			if err != nil {
				return nil, err
			}
		}
	}
	return roots, nil
}

func handlePendingAttestation(val []*pb.PendingAttestation, indices []uint64, convertAll bool) ([][32]byte, error) {
	length := len(indices)
	if convertAll {
		length = len(val)
	}
	roots := make([][32]byte, 0, length)
	hasher := hashutil.CustomSHA256Hasher()
	rootCreator := func(input *pb.PendingAttestation) error {
		newRoot, err := stateutil.PendingAttRootWithHasher(hasher, input)
		if err != nil {
			return err
		}
		roots = append(roots, newRoot)
		return nil
	}
	if convertAll {
		for i := range val {
			err := rootCreator(val[i])
			if err != nil {
				return nil, err
			}
		}
		return roots, nil
	}
	if len(val) > 0 {
		for _, idx := range indices {
			if idx > uint64(len(val))-1 {
				return nil, fmt.Errorf("index %d greater than number of pending attestations %d", idx, len(val))
			}
			err := rootCreator(val[idx])
			if err != nil {
				return nil, err
			}
		}
	}
	return roots, nil
}
