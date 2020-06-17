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

func (t *FabricClient) TableSet(key string, value string, fcn string) (string, error) {
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: fcn, Args: [][]byte{[]byte(key), []byte(value)}}
	resp, err := t.Exec(req)
	if err != nil {
		return "", err
	}
	return string(resp.TransactionID), nil
}

func (t *FabricClient) TableGet(key string, fcn string) ([]byte, error) {
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: fcn,
		Args: [][]byte{[]byte(key)},
	}
	resp, err := t.Exec(req)
	if err != nil {
		return nil, err
	}
	//fmt.Println("resp: ", resp)
	return resp.Payload, nil
}

func flushIncomeToFabric(s *FabricClient, income CharityIncome) error {
	return flushValueToFabric(s, income.Id, income, "addCharityIncome", func(txid string) {
		AddTxidToIncome(income, txid)
	})
}

func flushOutcomeToFabric(s *FabricClient, outcome CharityOutcome) error {
	return flushValueToFabric(s, outcome.Id, outcome, "addCharityOutcome", func(txid string) {
		AddTxidToOutcome(outcome, txid)
	})
}

func flushValueToFabric(s *FabricClient, key string, val interface{}, fcn string, cb func(txid string)) error {
	outcomeString, err := json.Marshal(val)
	if err != nil {
		fmt.Println("cannot marshal ", err)
		return err
	}
	fabricFlushLock.Lock()
	defer fabricFlushLock.Unlock()
	txid, err := s.TableSet(key, string(outcomeString), fcn)
	if err != nil {
		return err
	}
	go cb(txid)
	return nil
}
