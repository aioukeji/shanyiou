package main

import (
	"fmt"
	"io"
	"net/http"
	"net/rpc"
	"os"

	"github.com/aioukeji/shanyiou/chainsdk"
	"github.com/aioukeji/shanyiou/server"
)

func main() {
	channelClient, err := chainsdk.GetChannelClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	svc := server.FabricClient{
		ChaincodeID: chainsdk.ChaincodeID,
		Client:      channelClient,
		PeerTarget:  chainsdk.PeerTarget,
	}

	defer svc.Close()
	stub := &server.RpcStub{FabricClient: svc}
	_ = server.CharityApi(stub) // assert interface
	rpc.Register(stub)
	http.HandleFunc("/jsonrpc", func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			req.Body.Close()
		}()
		w.Header().Set("Content-Type", "application/json")
		res := server.NewRPCRequest(req.Body).Call()
		io.Copy(w, res)
	})
	fmt.Println("HTTP Server setup done")
	http.ListenAndServe(":8910", nil)
}
