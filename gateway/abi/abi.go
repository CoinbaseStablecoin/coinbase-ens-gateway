package abi

import (
	"log"
	"strings"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
)

var (
	// https://github.com/ensdomains/offchain-resolver/blob/main/packages/contracts/contracts/OffchainResolver.sol
	// resolve(bytes,bytes)
	IResolverService = mustParseABI(`[
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
					"name": "result",
					"type": "bytes"
				},
				{
					"internalType": "uint64",
					"name": "expires",
					"type": "uint64"
				},
				{
					"internalType": "bytes",
					"name": "sig",
					"type": "bytes"
				}
			],
			"stateMutability": "view",
			"type": "function"
		}
	]`)

	// EIP-137
	// https://github.com/ensdomains/ens-contracts/blob/v0.0.8/contracts/resolvers/profiles/IAddrResolver.sol
	// addr(bytes32)
	IAddrResolver = mustParseABI(`[
		{
			"inputs": [
				{
					"internalType": "bytes32",
					"name": "node",
					"type": "bytes32"
				}
			],
			"name": "addr",
			"outputs": [
				{
					"internalType": "address payable",
					"name": "",
					"type": "address"
				}
			],
			"stateMutability": "view",
			"type": "function"
		}
	]`)

	// EIP-2304
	// https://github.com/ensdomains/ens-contracts/blob/v0.0.8/contracts/resolvers/profiles/IAddressResolver.sol
	// addr(bytes32,uint256)
	// It's unfortunate that this is named very similar to IAddrResolver, but
	// that's the way it is in the official ENS repo.
	IAddressResolver = mustParseABI(`[
		{
			"inputs": [
				{
					"internalType": "bytes32",
					"name": "node",
					"type": "bytes32"
				},
				{
					"internalType": "uint256",
					"name": "coinType",
					"type": "uint256"
				}
			],
			"name": "addr",
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

	// EIP-634
	// https://github.com/ensdomains/ens-contracts/blob/v0.0.8/contracts/resolvers/profiles/ITextResolver.sol
	// text(bytes32,string)
	ITextResolver = mustParseABI(`[
		{
			"inputs": [
				{
					"internalType": "bytes32",
					"name": "node",
					"type": "bytes32"
				},
				{
					"internalType": "string",
					"name": "key",
					"type": "string"
				}
			],
			"name": "text",
			"outputs": [
				{
					"internalType": "string",
					"name": "",
					"type": "string"
				}
			],
			"stateMutability": "view",
			"type": "function"
		}
	]`)
)

var (
	SelectorIResolverServiceResolve = mustGetSelector(IResolverService, "resolve")
	SelectorIAddrResolverAddr       = mustGetSelector(IAddrResolver, "addr")
	SelectorIAddressResolverAddr    = mustGetSelector(IAddressResolver, "addr")
	SelectorITextResolverText       = mustGetSelector(ITextResolver, "text")
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
