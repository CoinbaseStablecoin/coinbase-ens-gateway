package abi

import (
	"log"
	"strings"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
)

var (
	IExtendedResolver = mustParseABI(`[
		{
			"inputs": [
				{
					"internalType": "bytes",
					"name": "name",
					"type": "bytes"
				},
				{
					"internalType": "bytes",
					"name": "data",
					"type": "bytes"
				}
			],
			"name": "resolve",
			"outputs": [
				{
					"internalType": "bytes",
					"name": "",
					"type": "bytes"
				}
			],
			"stateMutability": "view",
			"type": "function"
		}
	]`)
)

var (
	SelectorResolve = mustGetSelector(IExtendedResolver, "resolve")
)

func mustParseABI(json string) *ethabi.ABI {
	a, err := ethabi.JSON(strings.NewReader(json))
	if err != nil {
		log.Fatalln("could not parse ABI")
	}
	return &a
}

func mustGetSelector(parsedABI *ethabi.ABI, methodName string) []byte {
	method, ok := parsedABI.Methods[methodName]
	if !ok {
		log.Fatalln("could not find method:", methodName)
	}
	return method.ID
}
