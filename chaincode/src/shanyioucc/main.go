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
	if len(args) == 2 {
		if fun == "addCharityIncome" {
			return t.addCharityIncome(stub, args[0], args[1])
		} else if fun == "addCharityOutcome" {
			return t.addCharityOutcome(stub, args[0], args[1])
		}
	} else if len(args) == 1 {
		if fun == "getCharityIncome" {
			return t.getCharityIncome(stub, args[0])
		} else if fun == "getCharityOutcome" {
			return t.getCharityOutcome(stub, args[0])
		}
	}

	return shim.Error(fmt.Sprintf("invalid function name: %v", fun))

}

func main() {
	err := shim.Start(new(Chaincode))
	if err != nil {
		fmt.Printf("launch chaincode error: %s", err)
	}
}
