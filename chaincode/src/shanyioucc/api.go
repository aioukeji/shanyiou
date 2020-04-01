package main

import (
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"

	"strings"
)

func (t *Chaincode) tableGet(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("invalid param len")
	}
	tableName := args[0]
	key := args[1]
	if strings.Contains(tableName, ".") || strings.Contains(key, ".") {
		return shim.Error("表名或键名含有 '.' 无效字符")
	}
	value, err := stub.GetState(tableName + "." + key)
	if err != nil {
		return shim.Error(err.Error())
	}
	//fmt.Println("set key done", key)
	return shim.Success(value)
}

func (t *Chaincode) tableSet(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 3 {
		return shim.Error("invalid param len")
	}
	tableName := args[0]
	key := args[1]
	value := []byte(args[2])
	if strings.Contains(tableName, ".") || strings.Contains(key, ".") {
		return shim.Error("表名或键名含有 '.' 无效字符")
	}
	err := stub.PutState(tableName+"."+key, value)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("set key done", key)
	return shim.Success([]byte{})
}
