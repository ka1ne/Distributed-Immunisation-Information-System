package main_test

import (
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	//"github.com/ka1ne/fabric-samples/asset-transfer-basic/chaincode-go/chaincode/mocks"
	//"github.com/ka1ne/Distributed-Immunisation-Information-System/app/chaincode/mocks"
	//"github.com/ka1ne/fabric-samples/token-utxo/chaincode-go/chaincode"
	//"github.com/stretchr/testify/require"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate -o mocks/transaction.go -fake-name TransactionContext . transactionContext
type transactionContext interface {
	contractapi.TransactionContextInterface
}

//counterfeiter:generate -o mocks/chaincodestub.go -fake-name ChaincodeStub . chaincodeStub
type chaincodeStub interface {
	shim.ChaincodeStubInterface
}

//counterfeiter:generate -o mocks/statequeryiterator.go -fake-name StateQueryIterator . stateQueryIterator
type stateQueryIterator interface {
	shim.StateQueryIteratorInterface
}

func TestInitLedger(t *testing.T) {
	/*
		chaincodeStub := &mocks.ChaincodeStub{}
		transactionContext := &mocks.TransactionContext{}
		transactionContext.GetStubReturns(chaincodeStub)
	*/

	//assetTransfer := chaincode.SmartContract{}
	//err := assetTransfer.InitLedger(transactionContext)

}
