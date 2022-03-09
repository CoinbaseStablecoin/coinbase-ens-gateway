package server

import (
	"log"
	"math/big"
	"time"

	"github.cbhq.net/pete/coinbase-ens-gateway/resolver"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

const (
	TTL_SECONDS = 300 // 5 minutes
)

var ZeroAddress = common.BigToAddress(new(big.Int))

type resolveParams struct {
	Sender string `uri:"sender" binding:"required"`
	Data   string `uri:"data" binding:"required"`
}

func (s *Server) GetResolve(c *gin.Context) {
	var params resolveParams

	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(400, gin.H{"message": "invalid params"})
		return
	}

	result, err := s.Resolve(params.Sender, params.Data)
	if err != nil {
		log.Println("failed to resolve:", err)
		c.JSON(400, gin.H{"message": "invalid params"})
		return
	}

	c.JSON(200, gin.H{"data": hexutil.Encode(result)})
}

func (s *Server) Resolve(sender string, callDataHex string) ([]byte, error) {
	lookup, err := resolver.DecodeRequest(sender, callDataHex)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode request")
	}

	expires := uint64(time.Now().Unix() + TTL_SECONDS)

	var (
		resultData []byte
		hash       []byte
	)

	if ethLookup, ok := lookup.(*resolver.EthLookup); ok {
		if ethLookup.Name() == "pete.cbdev.eth" {
			resultData, hash, err = ethLookup.EncodeResult(common.HexToAddress("0x1111111111111111111111111111111111111111").Bytes(), expires)
		} else {
			resultData, hash, err = ethLookup.EncodeResult(ZeroAddress.Bytes(), expires)
		}
	} else if multicoinLookup, ok := lookup.(*resolver.MulticoinLookup); ok {
		if multicoinLookup.Name() == "pete.cbdev.eth" && multicoinLookup.CoinType().String() == "60" {
			resultData, hash, err = multicoinLookup.EncodeResult(common.HexToAddress("0x1111111111111111111111111111111111111111").Bytes(), expires)
		} else {
			resultData, hash, err = multicoinLookup.EncodeResult([]byte{}, expires)
		}
	} else if textLookup, ok := lookup.(*resolver.TextLookup); ok {
		if textLookup.Name() == "pete.cbdev.eth" && textLookup.Key() == "com.twitter" {
			resultData, hash, err = textLookup.EncodeResult([]byte("petejkim"), expires)
		} else {
			resultData, hash, err = textLookup.EncodeResult([]byte{}, expires)
		}
	}
	if err != nil {
		return nil, errors.New("unsupported lookup type")
	}

	sig, err := crypto.Sign(hash, s.privKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to sign")
	}
	sig[len(sig)-1] += 27

	responseData, err := resolver.EncodeResponse(resultData, expires, sig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to encode the response")
	}

	return responseData, nil
}
