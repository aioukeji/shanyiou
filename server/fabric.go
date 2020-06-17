package server

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

type FabricClient struct {
	ChaincodeID string
	PeerTarget  string
	Client      *channel.Client
	RawSDK      *fabsdk.FabricSDK
}

func (t *FabricClient) Close() {
	println("closing sdk")
	if t.RawSDK != nil {
		t.RawSDK.Close()
	} else {
		fmt.Println("empty FabricSDK, skip close")
	}

}

func (t *FabricClient) Exec(req channel.Request) (*channel.Response, error) {
	resp, err := t.Client.Execute(req, channel.WithTargetEndpoints(t.PeerTarget))
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
func (t *FabricClient) TableSet(table string, key string, value string) (string, error) {
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "tableSet", Args: [][]byte{[]byte(table), []byte(key), []byte(value)}}
	resp, err := t.Exec(req)
	if err != nil {
		return "", err
	}
	return string(resp.TransactionID), nil
}

func (t *FabricClient) TableGet(table string, key string) ([]byte, error) {
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "tableGet",
		Args: [][]byte{[]byte(table), []byte(key)},
	}
	resp, err := t.Exec(req)
	if err != nil {
		return nil, err
	}
	//fmt.Println("resp: ", resp)
	return resp.Payload, nil
}

func flushIncomeToFabric(s *FabricClient, income CharityIncome) error {
	return flushValueToFabric(s, tableNameIncome, income.Id, income, func(txid string) {
		AddTxidToIncome(income, txid)
	})
}

func flushOutcomeToFabric(s *FabricClient, outcome CharityOutcome) error {
	return flushValueToFabric(s, tableNameOutcome, outcome.Id, outcome, func(txid string) {
		AddTxidToOutcome(outcome, txid)
	})
}

func flushValueToFabric(s *FabricClient, tableName string, key string, val interface{}, cb func(txid string)) error {
	outcomeString, err := json.Marshal(val)
	if err != nil {
		fmt.Println("cannot marshal ", err)
		return err
	}
	fabricFlushLock.Lock()
	defer fabricFlushLock.Unlock()
	txid, err := s.TableSet(tableName, key, string(outcomeString))
	if err != nil {
		return err
	}
	go cb(txid)
	return nil
}
