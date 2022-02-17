package gateway

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"

	"github.cbhq.net/pete/coinbase-ens-gateway/gateway/abi"
	"github.cbhq.net/pete/coinbase-ens-gateway/pkg/namehash"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	dnsname "github.com/petejkim/ens-dnsname"
	"github.com/pkg/errors"
)

const (
	TTL_SECONDS = 300 // 5 minutes
)

func (gw *Gateway) Resolve(sender string, dataHex string) ([]byte, error) {
	// TODO: validate contract address
	if !common.IsHexAddress(sender) {
		return nil, errors.New("invalid sender address")
	}
	senderAddress := common.HexToAddress(sender)

	requestData, err := hexutil.Decode(dataHex)
	if err != nil {
		return nil, errors.Wrap(err, "data is not a valid hex string")
	}

	// check the first four-bytes to ensure that it's calling resolve(bytes,bytes)
	if !bytes.Equal(requestData[0:4], abi.SelectorIResolverServiceResolve) {
		return nil, errors.New("data is not a resolve call")
	}

	// decode resolve(bytes,bytes)
	decoded, err := abi.IResolverService.Methods["resolve"].Inputs.Unpack(requestData[4:])
	if err != nil {
		return nil, errors.Wrap(err, "resolve call data could not be decoded")
	}

	nameBytes := decoded[0].([]byte)
	callData := decoded[1].([]byte)

	// decode dns-encoded name
	name, err := dnsname.Decode(nameBytes)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse dns-encoded name")
	}

	selector := callData[0:4]
	data := callData[4:]
	var resultData []byte

	// use the right resolver lookup function based on the selector
	if bytes.Equal(selector, abi.SelectorIAddrResolverAddr) {
		resultData, err = gw.ResolveAddr(name, data)
	} else {
		return nil, errors.Errorf("unsupported resolver lookup: %s", hexutil.Encode(selector))
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to resolve")
	}

	expires := uint64(time.Now().Unix() + TTL_SECONDS)
	hash := HashResult(&senderAddress, expires, requestData, resultData)
	fmt.Println(hexutil.Encode(hash))

	sig, err := crypto.Sign(hash, gw.privKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to sign")
	}
	sig[len(sig)-1] += 27
	fmt.Println(hexutil.Encode(sig))

	// ABI-encode the result
	packed, err := abi.IResolverService.Methods["resolve"].Outputs.Pack(
		resultData, expires, sig,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to ABI-encode the result")
	}

	return packed, nil
}

func (gw *Gateway) ResolveAddr(name string, data []byte) ([]byte, error) {
	nh, err := namehash.NameHash(name)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get namehash")
	}

	// verify that data matches the name hash
	if !bytes.Equal(data, nh[:]) {
		return nil, errors.New("data does not match the name hash")
	}

	result, err := abi.IAddrResolver.Methods["addr"].Outputs.Pack(
		common.HexToAddress("0x3333333333333333333333333333333333333333"),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to ABI-encode the result")
	}
	return result, nil
}

func HashResult(target *common.Address, expires uint64, request []byte, result []byte) []byte {
	expiresBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(expiresBytes, expires)

	// https://github.com/ensdomains/offchain-resolver/blob/main/packages/contracts/contracts/SignatureVerifier.sol#L15
	// keccak256(0x1900 . target . expires . keccak256(request) . keccak256(result))
	return crypto.Keccak256(
		[]byte{0x19, 0x00},
		target.Bytes(),
		expiresBytes,
		crypto.Keccak256(request),
		crypto.Keccak256(result),
	)
}
