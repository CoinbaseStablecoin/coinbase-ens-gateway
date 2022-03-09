package resolver

import (
	"bytes"

	"github.cbhq.net/pete/coinbase-ens-gateway/resolver/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	dnsname "github.com/petejkim/ens-dnsname"
	"github.com/pkg/errors"
)

func DecodeRequest(sender string, requestCallDataHex string) (Lookup, error) {
	if !common.IsHexAddress(sender) {
		return nil, errors.New("sender is not a valid address")
	}
	senderAddress := common.HexToAddress(sender)

	requestCallData, err := hexutil.Decode(requestCallDataHex)
	if err != nil {
		return nil, errors.Wrap(err, "data is not a valid hex string")
	}

	// check the first four-bytes to ensure that it's calling resolve(bytes,bytes)
	if !bytes.Equal(requestCallData[0:4], abi.SelectorIResolverServiceResolve) {
		return nil, errors.New("data is not a resolve call")
	}

	// decode resolve(bytes,bytes)
	decoded, err := abi.IResolverService.Methods["resolve"].Inputs.Unpack(requestCallData[4:])
	if err != nil {
		return nil, errors.Wrap(err, "resolve call data could not be decoded")
	}

	dnsNameBytes := decoded[0].([]byte)
	lookupCallData := decoded[1].([]byte)

	// decode dns-encoded name
	name, err := dnsname.Decode(dnsNameBytes)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse dns-encoded name")
	}

	lookupSelector := lookupCallData[0:4]
	lookupInputs := lookupCallData[4:]

	// use the right resolver lookup function based on the selector
	if bytes.Equal(lookupSelector, abi.SelectorIAddrResolverAddr) {
		// addr(bytes32)
		return NewEthLookup(name, lookupInputs, &senderAddress, requestCallData)
	} else if bytes.Equal(lookupSelector, abi.SelectorIMulticoinResolverAddr) {
		// addr(bytes32,uint256)
		return NewMulticoinLookup(name, lookupInputs, &senderAddress, requestCallData)
	} else if bytes.Equal(lookupSelector, abi.SelectorITextResolverText) {
		// text(bytes32,string)
		return NewTextLookup(name, lookupInputs, &senderAddress, requestCallData)
	}

	return nil, errors.Errorf("unsupported lookup: %s", hexutil.Encode(lookupSelector))
}

func EncodeResponse(resultData []byte, expires uint64, sigRSV []byte) (responseData []byte, err error) {
	if responseData, err = abi.IResolverService.Methods["resolve"].Outputs.Pack(
		resultData, expires, sigRSV,
	); err != nil {
		return nil, errors.Wrap(err, "failed to ABI-encode the result")
	}

	return responseData, nil
}
