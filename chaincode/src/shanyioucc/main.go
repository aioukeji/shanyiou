package main

import (
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

type Chaincode struct {
}

func (t *Chaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (t *Chaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fun, args := stub.GetFunctionAndParameters()

	if fun == "tableGet" {
		return t.tableGet(stub, args)
	} else if fun == "tableSet" {
		return t.tableSet(stub, args)
	}

	return shim.Error(fmt.Sprintf("invalid function name: %v", fun))

}

func main() {
	err := shim.Start(new(Chaincode))
	if err != nil {
		fmt.Printf("launch chaincode error: %s", err)
	}
}
