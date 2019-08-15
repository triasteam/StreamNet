package main

import (
	v "./vue"
	auth "./auth"
	"crypto"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/AddNode", AddNode)
	http.HandleFunc("/QueryNodes", QueryNodes)
	http.HandleFunc("/QueryNodeDetail", QueryNodeDetail)
	err := http.ListenAndServe("0.0.0.0:8000", nil)

	if err != nil {
		fmt.Println(err)
	}
}

func AddNode(writer http.ResponseWriter, request *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	var addNodeRequest *v.AddNodeRequest
	if err := json.NewDecoder(request.Body).Decode(&addNodeRequest); err != nil {
		fmt.Println(err)
		request.Body.Close()
	}

	if validPrivilege(addNodeRequest.Address, addNodeRequest.Sign) == false{
		return;
	}

	var o v.OCli
	response := o.AddAttestationInfoFunction(addNodeRequest)

	if err := json.NewEncoder(writer).Encode(response); err != nil {
		fmt.Println(err)
	}
}

func QueryNodes(writer http.ResponseWriter, request *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	var queryNodesRequest *v.QueryNodesRequest
	if err := json.NewDecoder(request.Body).Decode(&queryNodesRequest); err != nil {
		fmt.Println(err)
		request.Body.Close()
	}

	if validPrivilege(queryNodesRequest.Address, queryNodesRequest.Sign) == false{
		return;
	}
	var o v.OCli
	response := o.GetRankFunction(queryNodesRequest)

	if err := json.NewEncoder(writer).Encode(response); err != nil {
		fmt.Println(err)
	}
}

func validPrivilege(addr string, sign string) bool{
	if sign != "" {
		address := addr
		base64Sig := sign
		sig,_ := base64.StdEncoding.DecodeString(base64Sig);
		var r auth.RSAUtil
		b := r.Verify([]byte(address), sig, crypto.SHA256, "./auth/public_key.pem")
		if b != nil {
			fmt.Println("has no privilege. address:", address)
			return false;
		}
		fmt.Println("privilege valid success !" )
	}
	return true;
}

func QueryNodeDetail(writer http.ResponseWriter, request *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	var detailRequest *v.NodeDetailRequest
	if err := json.NewDecoder(request.Body).Decode(&detailRequest); err != nil {
		fmt.Println(err)
		request.Body.Close()
	}

	var o v.OCli
	response := o.QueryNodeDetail(detailRequest)
	if err := json.NewEncoder(writer).Encode(response); err != nil {
		fmt.Println(err)
	}
}
