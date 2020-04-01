package chainsdk

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"

	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"

	"github.com/pkg/errors"
)

func NewSDK(ConfigFile string) (*fabsdk.FabricSDK, error) {
	sdk, err := fabsdk.New(config.FromFile(ConfigFile))
	if err != nil {
		return nil, fmt.Errorf("failed to create SDK: %v", err)
	}
	fmt.Println("SDK created")
	return sdk, nil
}

func CreateChannel(sdk *fabsdk.FabricSDK, info *InitInfo) error {
	// The resource management client is responsible for managing channels (create/update channel)
	resourceManagerClientContext := sdk.Context(fabsdk.WithUser(info.OrgAdmin), fabsdk.WithOrg(info.OrgName))
	var err error
	resMgmtClient, err := resmgmt.New(resourceManagerClientContext)
	orgResMgmt = resMgmtClient
	if err != nil {
		return errors.WithMessage(err, "failed to create channel management client from Admin identity")
	}
	fmt.Println("Resource management client created")
	//return nil

	existed := false

	allChannels, err := resMgmtClient.QueryChannels(resmgmt.WithTargetEndpoints(PeerTarget))
	if err != nil {
		return fmt.Errorf("query channel failed %v", err)
	}
	fmt.Println("existed channels", allChannels.Channels)
	for _, item := range allChannels.Channels {
		if item.ChannelId == info.ChannelID {
			existed = true
			break
		}
	}

	if existed {
		fmt.Println("channel already existed, skip")
	} else {

		// The MSP client allow us to retrieve user information from their identity, like its signing identity which we will need to save the channel
		mspClient, err := mspclient.New(sdk.Context(), mspclient.WithOrg(info.OrgName))
		if err != nil {
			return errors.WithMessage(err, "failed to create MSP client")
		}
		adminIdentity, err := mspClient.GetSigningIdentity(info.OrgAdmin)
		if err != nil {
			return errors.WithMessage(err, "failed to get admin signing identity")
		}
		req := resmgmt.SaveChannelRequest{ChannelID: info.ChannelID, ChannelConfigPath: info.ChannelConfig, SigningIdentities: []msp.SigningIdentity{adminIdentity}}
		txID, err := resMgmtClient.SaveChannel(req, resmgmt.WithOrdererEndpoint(info.OrdererOrgName))
		if err != nil || txID.TransactionID == "" {
			return errors.WithMessage(err, "failed to save channel")
		}
		fmt.Println("Channel created")

		// Make admin user join the previously created channel
		if err = resMgmtClient.JoinChannel(info.ChannelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(info.OrdererOrgName)); err != nil {
			return errors.WithMessage(err, "failed to make admin join channel")
		}
		fmt.Println("Channel joined")
	}

	fmt.Println("Initialization Successful")
	return nil
}

func InstallAndInstantiateCC(info *InitInfo) error {
	l, err := orgResMgmt.QueryInstalledChaincodes(resmgmt.WithTargetEndpoints(PeerTarget))
	if err != nil {
		return errors.WithMessage(err, "cannot list chaincodes")
	}
	installed := false
	fmt.Println("installed chaincodes", l.Chaincodes)

	for _, item := range l.Chaincodes {
		if item.Name == info.ChaincodeID {
			installed = true
			break
		}
	}

	if installed {
		fmt.Println("Chaincode already installed, skip")
	} else {
		// Create the chaincode package that will be sent to the peers
		ccPkg, err := gopackager.NewCCPackage(info.ChaincodePath, info.ChaincodeGoPath)
		if err != nil {
			return errors.WithMessage(err, "failed to create chaincode package")
		}
		fmt.Println("ccPkg created")

		// Install example cc to org peers
		installCCReq := resmgmt.InstallCCRequest{Name: info.ChaincodeID, Path: info.ChaincodePath, Version: ChaincodeVersion, Package: ccPkg}
		resp, err := orgResMgmt.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
		if err != nil {
			return errors.WithMessage(err, "failed to install chaincode")
		}
		fmt.Println("install response:", resp)
		fmt.Println("Chaincode installed")
	}

	l, err = orgResMgmt.QueryInstantiatedChaincodes(info.ChannelID, resmgmt.WithTargetEndpoints(PeerTarget))
	if err != nil {
		return errors.WithMessage(err, "cannot list chaincodes")
	}

	fmt.Println("instantiated chaincodes", l.Chaincodes)
	instantiated := false

	for _, item := range l.Chaincodes {
		if item.Name == info.ChaincodeID {
			instantiated = true
			break
		}
	}

	if instantiated {
		fmt.Println("Chaincode already instantiated, skip")
	} else {
		// Set up chaincode policy
		ccPolicy := cauthdsl.SignedByAnyMember([]string{"Org1MSP"})
		fmt.Println("cc policy", ccPolicy.String())

		resp, err := orgResMgmt.InstantiateCC(info.ChannelID, resmgmt.InstantiateCCRequest{Name: info.ChaincodeID,
			Path: info.ChaincodeGoPath, Version: ChaincodeVersion, Args: [][]byte{[]byte("init")}, Policy: ccPolicy})
		if err != nil || resp.TransactionID == "" {
			return errors.WithMessage(err, "failed to instantiate the chaincode")
		}
		fmt.Println("Chaincode instantiated")
	}
	return nil
}

func GetChannelContext(sdk *fabsdk.FabricSDK, info *InitInfo) (*channel.Client, error) {
	clientChannelContext := sdk.ChannelContext(info.ChannelID, fabsdk.WithUser(info.UserName),
		fabsdk.WithOrg(info.OrgName))
	// returns a Client instance. Channel client can query chaincode, execute chaincode and register/unregister for chaincode events on specific channel.
	channelClient, err := channel.New(clientChannelContext)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new channel client")
	}
	fmt.Println("Channel client created")

	return channelClient, nil
}

func GetChannelClient() (*channel.Client, error) {
	sdk, err := NewSDK(ConfigFile)
	if err != nil {
		fmt.Printf(err.Error())
		return nil, err
	}

	err = CreateChannel(sdk, GInitInfo)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	err = InstallAndInstantiateCC(GInitInfo)
	if err != nil {
		fmt.Printf("install cc failed: %v\n", err)
		return nil, err
	}

	return GetChannelContext(sdk, GInitInfo)
}
