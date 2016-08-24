package main

/******************** Testing Objective consensu:STATE TRANSFER ********
*   Setup: 4 node local docker peer network with security
*   0. Deploy chaincodeexample02 with 100000, 90000 as initial args
*   1. Send Invoke Requests on multiple peers using go routines.
*   2. Verify query results match on PEER0 and PEER1 after invoke
*********************************************************************/

import (
	"fmt"
	//"strconv"
	"time"

	"obcsdk/chaincode"
	"obcsdk/peernetwork"
	"sync"
)

func main() {

	fmt.Println("Creating a local docker network")
	peernetwork.SetupLocalNetwork(4, true)

	_ = chaincode.InitNetwork()
	chaincode.InitChainCodes()
	chaincode.RegisterUsers()

	time.Sleep(30000 * time.Millisecond)

	fmt.Println("\nPOST/Chaincode: Deploying chaincode at the beginning ....")
	dAPIArgs0 := []string{"example02", "init"}
	depArgs0 := []string{"a", "100000", "b", "90000"}
	chaincode.Deploy(dAPIArgs0, depArgs0)

	//var resa, resb string
	var inita, initb, curra, currb int
	inita = 100000
	initb = 90000
	curra = inita
	currb = initb

	time.Sleep(60000 * time.Millisecond)
	fmt.Println("\nPOST/Chaincode: Querying a and b after deploy >>>>>>>>>>> ")
	qAPIArgs0 := []string{"example02", "query"}
	qArgsa := []string{"a"}
	qArgsb := []string{"b"}
	A, _ := chaincode.Query(qAPIArgs0, qArgsa)
	B, _ := chaincode.Query(qAPIArgs0, qArgsb)
	myStr := fmt.Sprintf("\nA = %s B= %s", A, B)
	fmt.Println(myStr)

	numReq := 250 
	InvokeLoop(numReq)

	time.Sleep(120000 * time.Millisecond)
	curra = curra - 20
	currb = currb + 20


	fmt.Println("\nPOST/Chaincode: Querying a and b after invoke >>>>>>>>>>> ")
	qAPIArgs00 := []string{"example02", "query", "PEER0"}
	qAPIArgs01 := []string{"example02", "query", "PEER1"}
	qAPIArgs02 := []string{"example02", "query", "PEER2"}
	qAPIArgs03 := []string{"example02", "query", "PEER3"}

	res0A, _ := chaincode.QueryOnHost(qAPIArgs00, qArgsa)
	res0B, _ := chaincode.QueryOnHost(qAPIArgs00, qArgsb)

	res1A, _ := chaincode.QueryOnHost(qAPIArgs01, qArgsa)
	res1B, _ := chaincode.QueryOnHost(qAPIArgs01, qArgsb)

	res2A, _ := chaincode.QueryOnHost(qAPIArgs02, qArgsa)
	res2B, _ := chaincode.QueryOnHost(qAPIArgs02, qArgsb)

	res3A, _ := chaincode.QueryOnHost(qAPIArgs03, qArgsa)
	res3B, _ := chaincode.QueryOnHost(qAPIArgs03, qArgsb)

	fmt.Println("Results in a and b PEER0 : ", res0A, res0B)
	fmt.Println("Results in a and b PEER1 : ", res1A, res1B)
	fmt.Println("Results in a and b PEER2 : ", res2A, res2B)
	fmt.Println("Results in a and b PEER3 : ", res3A, res3B)

	ht0, _ := chaincode.GetChainHeight("PEER0")
	ht1, _ := chaincode.GetChainHeight("PEER1")
	ht2, _ := chaincode.GetChainHeight("PEER2")
	ht3, _ := chaincode.GetChainHeight("PEER3")

	fmt.Printf("ht0: %d, ht1: %d, ht2: %d, ht3: %d ", ht0, ht1, ht2, ht3)

}

func InvokeLoop(numReq int) {

	var wg sync.WaitGroup

	invArgs0 := []string{"a", "b", "1"}
	iAPIArgsCurrPeer1 := []string{"example02", "invoke", "PEER1"}
	wg.Add(2)
	go func() {

		defer wg.Done()
		k := 1
		for k <= numReq {
		   go chaincode.InvokeOnPeer(iAPIArgsCurrPeer1, invArgs0)
		   k++
		}
		fmt.Println("# of Req Invoked on PEER1 ", k)
	}()

	go func() {

		defer wg.Done()
		iAPIArgsCurrPeer3 := []string{"example02", "invoke", "PEER3"}
		k := 1
		for k <= numReq {
		    go chaincode.InvokeOnPeer(iAPIArgsCurrPeer3, invArgs0)
		    k++
		}
		fmt.Println("# of Req Invoked  on PEER3", k)
	}()

	wg.Wait()
}
