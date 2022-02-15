package gateway

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
)

type Gateway struct {
	privKey *ecdsa.PrivateKey
}

func New(privKeyHex string) (*Gateway, error) {
	privKeyB, err := hexutil.Decode(privKeyHex)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode private key")
	}
	privKey, err := crypto.ToECDSA(privKeyB)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize private key")
	}

	gw := &Gateway{privKey}

	return gw, err
}

func (gw *Gateway) SignerAddress() string {
	pubKey := gw.privKey.Public()
	ecdsaPubKey, _ := pubKey.(*ecdsa.PublicKey)

	return crypto.PubkeyToAddress(*ecdsaPubKey).Hex()
}
