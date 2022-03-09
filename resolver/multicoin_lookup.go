package resolver

import (
	"bytes"
	"math/big"

	"github.cbhq.net/pete/coinbase-ens-gateway/pkg/namehash"
	"github.cbhq.net/pete/coinbase-ens-gateway/resolver/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

var _ Lookup = (*MulticoinLookup)(nil)

type MulticoinLookup struct {
	name          string
	senderAddress *common.Address
	requestData   []byte
	coinType      *big.Int
}

func NewMulticoinLookup(name string, lookupInputs []byte, senderAddress *common.Address, requestData []byte) (*MulticoinLookup, error) {
	nh, err := namehash.NameHash(name)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get namehash")
	}

	decoded, err := abi.IMulticoinResolver.Methods["addr"].Inputs.Unpack(lookupInputs)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode lookup inputs")
	}

	node, ok := decoded[0].([32]byte) // bytes32
	if !ok {
		return nil, errors.New(`failed to decode "node" in lookup inputs`)
	}

	coinType, ok := decoded[1].(*big.Int) // uint256
	if !ok {
		return nil, errors.New(`failed to decode "coinType" in lookup inputs`)
	}

	if !bytes.Equal(node[:], nh[:]) {
		return nil, errors.New("name hash does not match the lookup input")
	}

	return &MulticoinLookup{name, senderAddress, requestData, coinType}, nil
}

func (l *MulticoinLookup) Name() string {
	return l.name
}

func (l *MulticoinLookup) CoinType() *big.Int {
	bi := new(big.Int)
	return bi.Add(l.coinType, bi)
}

func (l *MulticoinLookup) EncodeResult(result []byte, expires uint64) (resultData []byte, hash []byte, err error) {
	if resultData, err = abi.IMulticoinResolver.Methods["addr"].Outputs.Pack(
		result, // bytes
	); err != nil {
		return nil, nil, errors.Wrap(err, "failed to ABI-encode the result")
	}

	hash = HashResult(l.senderAddress, expires, l.requestData, resultData)

	return resultData, hash, nil
}
