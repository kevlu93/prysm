module github.com/prysmaticlabs/prysm

go 1.16

require (
	contrib.go.opencensus.io/exporter/jaeger v0.2.1
	github.com/allegro/bigcache v1.2.1 // indirect
	github.com/aristanetworks/goarista v0.0.0-20200521140103-6c3304613b30
	github.com/bazelbuild/buildtools v0.0.0-20200528175155-f4e8394f069d
	github.com/bazelbuild/rules_go v0.23.2
	github.com/btcsuite/btcd v0.21.0-beta // indirect
	github.com/cespare/cp v1.1.1 // indirect
	github.com/confluentinc/confluent-kafka-go v1.4.2 // indirect
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/d4l3k/messagediff v1.2.1
	github.com/davidlazar/go-crypto v0.0.0-20200604182044-b73af7476f6c // indirect
	github.com/deckarep/golang-set v1.7.1 // indirect
	github.com/dgraph-io/ristretto v0.0.3
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/emicklei/dot v0.11.0
	github.com/ethereum/go-ethereum v0.0.0-00010101000000-000000000000
	github.com/fatih/color v1.9.0 // indirect
	github.com/ferranbt/fastssz v0.0.0-20210120143747-11b9eff30ea9
	github.com/fjl/memsize v0.0.0-20190710130421-bcb5799ab5e5
	github.com/fsnotify/fsnotify v1.4.9
	github.com/gballet/go-libpcsclite v0.0.0-20191108122812-4678299bea08 // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/go-logr/logr v0.2.1 // indirect
	github.com/go-yaml/yaml v2.1.0+incompatible
	github.com/gogo/protobuf v1.3.2
	github.com/golang/gddo v0.0.0-20200528160355-8d077c1d8f4c
	github.com/golang/mock v1.4.4
	github.com/golang/protobuf v1.4.3
	github.com/golang/snappy v0.0.3
	github.com/google/gofuzz v1.2.0
	github.com/google/gopacket v1.1.19 // indirect
	github.com/google/uuid v1.2.0
	github.com/graph-gophers/graphql-go v0.0.0-20200309224638-dae41bde9ef9 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.2
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/hashicorp/golang-lru v0.5.4
	github.com/herumi/bls-eth-go-binary v0.0.0-20210130185500-57372fb27371
	github.com/ianlancetaylor/cgosymbolizer v0.0.0-20200424224625-be1b05b0b279
	github.com/influxdata/influxdb v1.8.0 // indirect
	github.com/ipfs/go-ipfs-addr v0.0.1
	github.com/ipfs/go-log/v2 v2.1.1
	github.com/joonix/log v0.0.0-20200409080653-9c1d2ceb5f1d
	github.com/json-iterator/go v1.1.10
	github.com/k0kubun/go-ansi v0.0.0-20180517002512-3bf9e2903213
	github.com/karalabe/usb v0.0.0-20191104083709-911d15fe12a9 // indirect
	github.com/kevinms/leakybucket-go v0.0.0-20200115003610-082473db97ca
	github.com/koron/go-ssdp v0.0.2 // indirect
	github.com/kr/pretty v0.2.1
	github.com/kr/text v0.2.0 // indirect
	github.com/libp2p/go-libp2p v0.12.1-0.20201208224947-3155ff3089c0
	github.com/libp2p/go-libp2p-blankhost v0.2.0
	github.com/libp2p/go-libp2p-core v0.7.0
	github.com/libp2p/go-libp2p-noise v0.1.2
	github.com/libp2p/go-libp2p-peerstore v0.2.6
	github.com/libp2p/go-libp2p-pubsub v0.4.0
	github.com/libp2p/go-libp2p-swarm v0.3.1
	github.com/libp2p/go-libp2p-tls v0.1.4-0.20200421131144-8a8ad624a291 // indirect
	github.com/libp2p/go-libp2p-yamux v0.4.1 // indirect
	github.com/libp2p/go-netroute v0.1.4 // indirect
	github.com/libp2p/go-sockaddr v0.1.0 // indirect
	github.com/libp2p/go-tcp-transport v0.2.1
	github.com/logrusorgru/aurora v2.0.3+incompatible
	github.com/lunixbochs/vtclean v1.0.0 // indirect
	github.com/manifoldco/promptui v0.7.0
	github.com/minio/highwayhash v1.0.1
	github.com/minio/sha256-simd v0.1.1
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826
	github.com/multiformats/go-multiaddr v0.3.1
	github.com/nbutton23/zxcvbn-go v0.0.0-20180912185939-ae427f1e4c1d
	github.com/olekukonko/tablewriter v0.0.4 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/paulbellamy/ratecounter v0.2.0
	github.com/pborman/uuid v1.2.1
	github.com/peterh/liner v1.2.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/prestonvanloon/go-recaptcha v0.0.0-20190217191114-0834cef6e8bd
	github.com/prometheus/client_golang v1.9.0
	github.com/prometheus/procfs v0.3.0 // indirect
	github.com/prometheus/tsdb v0.10.0 // indirect
	github.com/prysmaticlabs/eth2-types v0.0.0-20210219172114-1da477c09a06
	github.com/prysmaticlabs/ethereumapis v0.0.0-20210311175904-cf9f64632dd4
	github.com/prysmaticlabs/go-bitfield v0.0.0-20210202205921-7fcea7c45dc8
	github.com/prysmaticlabs/prombbolt v0.0.0-20210126082820-9b7adba6db7c
	github.com/rs/cors v1.7.0
	github.com/schollz/progressbar/v3 v3.3.4
	github.com/sirupsen/logrus v1.6.0
	github.com/status-im/keycard-go v0.0.0-20200402102358-957c09536969 // indirect
	github.com/stretchr/testify v1.6.1
	github.com/supranational/blst v0.3.3
	github.com/trailofbits/go-mutexasserts v0.0.0-20200708152505-19999e7d3cef
	github.com/tyler-smith/go-bip39 v1.0.2
	github.com/urfave/cli/v2 v2.2.0
	github.com/wealdtech/go-eth2-util v1.6.3
	github.com/wealdtech/go-eth2-wallet-encryptor-keystorev4 v1.1.3
	github.com/wercker/journalhook v0.0.0-20180428041537-5d0a5ae867b3
	github.com/x-cray/logrus-prefixed-formatter v0.5.2
	go.etcd.io/bbolt v1.3.5
	go.opencensus.io v0.22.5
	go.uber.org/automaxprocs v1.3.0
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.16.0 // indirect
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad
	golang.org/x/exp v0.0.0-20200513190911-00229845015e
	golang.org/x/net v0.0.0-20201209123823-ac852fbbde11 // indirect
	golang.org/x/time v0.0.0-20200630173020-3af7569d3a1e // indirect
	golang.org/x/tools v0.0.0-20210106214847-113979e3529a
	google.golang.org/api v0.34.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20201026171402-d4b8fe4fd877
	google.golang.org/grpc v1.36.0
	google.golang.org/protobuf v1.25.0
	gopkg.in/confluentinc/confluent-kafka-go.v1 v1.4.2
	gopkg.in/d4l3k/messagediff.v1 v1.2.1
	gopkg.in/errgo.v2 v2.1.0
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.18.3
	k8s.io/apimachinery v0.18.3
	k8s.io/client-go v0.18.3
	k8s.io/klog/v2 v2.3.0 // indirect
	k8s.io/utils v0.0.0-20200520001619-278ece378a50 // indirect
)

replace github.com/ethereum/go-ethereum => github.com/prysmaticlabs/bazel-go-ethereum v0.0.0-20210311182818-b9c0d847ae7b

replace github.com/json-iterator/go => github.com/prestonvanloon/go v1.1.7-0.20190722034630-4f2e55fcf87b
