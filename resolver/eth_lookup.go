package resolver

import (
	"bytes"

	"github.cbhq.net/pete/coinbase-ens-gateway/pkg/namehash"
	"github.cbhq.net/pete/coinbase-ens-gateway/resolver/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

var _ Lookup = (*EthLookup)(nil)

type EthLookup struct {
	name          string
	senderAddress *common.Address
	requestData   []byte
}

func NewEthLookup(name string, lookupInputs []byte, senderAddress *common.Address, requestData []byte) (*EthLookup, error) {
	nh, err := namehash.NameHash(name)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get namehash")
	}

	decoded, err := abi.IAddrResolver.Methods["addr"].Inputs.Unpack(lookupInputs)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode lookup inputs")
	}

	node, ok := decoded[0].([32]byte) // bytes32
	if !ok {
		return nil, errors.New(`failed to decode "node" in lookup inputs`)
	}

	if !bytes.Equal(node[:], nh[:]) {
		return nil, errors.New("name hash does not match the lookup input")
	}

	return &EthLookup{name, senderAddress, requestData}, nil
}

func (l *EthLookup) Name() string {
	return l.name
}

func (l *EthLookup) EncodeResult(result []byte, expires uint64) (resultData []byte, hash []byte, err error) {
	if len(result) != 20 {
		return nil, nil, errors.New("address must be 20 bytes long")
	}

	if resultData, err = abi.IAddrResolver.Methods["addr"].Outputs.Pack(
		common.BytesToAddress(result), // address
	); err != nil {
		return nil, nil, errors.Wrap(err, "failed to ABI-encode the result")
	}

	hash = HashResult(l.senderAddress, expires, l.requestData, resultData)

	return resultData, hash, nil
}
