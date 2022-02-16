package gateway

import (
	"bytes"
	"fmt"

	"github.cbhq.net/pete/coinbase-ens-gateway/gateway/abi"
	"github.cbhq.net/pete/coinbase-ens-gateway/pkg/namehash"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	dnsname "github.com/petejkim/ens-dnsname"
	"github.com/pkg/errors"
)

func (gw *Gateway) Resolve(sender string, dataHex string) (string, error) {
	if !common.IsHexAddress(sender) {
		return "", errors.New("invalid sender address")
	}

	encodedData, err := hexutil.Decode(dataHex)
	if err != nil {
		return "", errors.Wrap(err, "data is not a valid hex string")
	}

	// check the first four-bytes to ensure that it's calling resolve(bytes,bytes)
	if !bytes.Equal(encodedData[0:4], abi.SelectorResolve) {
		return "", errors.New("data is not a resolve call")
	}

	// decode resolve(bytes,bytes)
	decoded, err := abi.IExtendedResolver.Methods["resolve"].Inputs.Unpack(encodedData[4:])
	if err != nil {
		return "", errors.Wrap(err, "resolve call data could not be decoded")
	}

	name, ok := decoded[0].([]byte)
	if !ok {
		return "", errors.New("invalid decoded data")
	}
	data, ok := decoded[1].([]byte)
	if !ok {
		return "", errors.New("invalid decoded data")
	}

	fmt.Println("name:", hexutil.Encode(name))
	fmt.Println("data:", hexutil.Encode(data))

	decodedName, err := dnsname.Decode(name)
	if err != nil {
		return "", errors.New("failed to parse dns-encoded name")
	}
	fmt.Println("decodedName:", decodedName)

	nh, err := namehash.NameHash(decodedName)
	if err != nil {
		return "", errors.New("failed to get namehash")
	}
	fmt.Println("namehash:", hexutil.Encode(nh[:]))

	return "", nil
}
