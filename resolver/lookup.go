package resolver

import (
	"encoding/binary"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type Lookup interface {
	Name() string
	EncodeResult(result []byte, expires uint64) (resultData []byte, hash []byte, err error)
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
