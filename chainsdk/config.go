package chainsdk

import "github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"

type InitInfo struct {
	ChannelID      string
	ChannelConfig  string
	OrgAdmin       string
	OrgName        string
	OrdererOrgName string

	ChaincodeID     string
	ChaincodeGoPath string
	ChaincodePath   string
	UserName        string
}

var orgResMgmt *resmgmt.Client

const (
	ConfigFile       = "config.yaml"
	ChaincodeID      = "shanyiou"
	ChaincodeVersion = "1.0"
	PeerTarget       = "peer0.org1.example.com"
)

var GInitInfo = &InitInfo{
	ChannelID:     "mychannel",
	ChannelConfig: "fixtures/artifacts/mychannel.tx",

	OrgAdmin:       "Admin",
	OrgName:        "Org1",
	OrdererOrgName: "orderer.example.com",

	ChaincodeID:     ChaincodeID,
	ChaincodeGoPath: "chaincode",
	ChaincodePath:   "shanyioucc/",
	UserName:        "User1",
}
