package main

import (
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

//善款捐赠
func (t *Chaincode) addCharityIncome(stub shim.ChaincodeStubInterface, incomeId string, value string) peer.Response {
	key := "charityIncome" + "." + incomeId
	err := stub.PutState(key, []byte(value))
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("add charity income done", key)
	return shim.Success([]byte{})
}

//善款支出
func (t *Chaincode) addCharityOutcome(stub shim.ChaincodeStubInterface, outcomeId string, value string) peer.Response {
	key := "charityOutcome" + "." + outcomeId
	err := stub.PutState(key, []byte(value))
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("add charity outcome done", key)
	return shim.Success([]byte{})
}

func (t *Chaincode) getCharityIncome(stub shim.ChaincodeStubInterface, incomeId string) peer.Response {
	key := "charityIncome" + "." + incomeId
	value, err := stub.GetState(key)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(value)
}

func (t *Chaincode) getCharityOutcome(stub shim.ChaincodeStubInterface, outcomeId string) peer.Response {
	key := "charityOutcome" + "." + outcomeId
	value, err := stub.GetState(key)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(value)
}
