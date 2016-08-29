package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"obcsdk/chaincode"
	"obcsdk/peernetwork"
	"obcsdk/util"
)

/********** Test Objective : Ledger Stress with 4 Clients, 4 Peers ************
*
*   Setup: 4 node peer network with security enabled
*   1. Deploy chaincode https://goo.gl/TysS79
*   2. Invoke 5K txns from each client simultaneously on 4 different peers
*   3. Check if the counter value(20000) matches with query on "counter"
*
* USAGE: NETWORK="LOCAL" go run LedgerStressOneCliOnePeer.go 
*  This NETWORK env value could be LOCAL or Z
*********************************************************************/
var peerNetworkSetup peernetwork.PeerNetwork
var AVal, BVal, curAVal, curBVal, invokeValue int64
var argA = []string{"a"}
var argB = []string{"counter"}
var counter int64
var wg sync.WaitGroup

const (
	TRX_COUNT = 20000
	CLIENTS   = 4
)

func initNetwork() {
	util.Logger("========= Init Network =========")
	//peernetwork.GetNC_Local()
	peerNetworkSetup = chaincode.InitNetwork()
	chaincode.InitChainCodes()
	util.Logger("========= Register Users =========")
	chaincode.RegisterCustomUsers()
}

func invokeChaincode(peer string) {
	counter++
	arg1Construct := []string{util.CHAINCODE_NAME, util.INVOKE, peer}
	arg2Construct := []string{"a" + strconv.FormatInt(counter, 10), util.DATA, "counter"}

	_, _ = chaincode.InvokeOnPeer(arg1Construct, arg2Construct)
}

func Init() {
	//initialize
	done := make(chan bool, 1)
	counter = 0
	wg.Add(CLIENTS)
	// Setup the network based on the NetworkCredentials.json provided
	initNetwork()

	//Deploy chaincode
	util.DeployChaincode(done)
}

func InvokeLoop() {
	curTime := time.Now()
	go func() {
		for i := 1; i <= TRX_COUNT/CLIENTS; i++ {
			if counter > 0 && counter%1000 == 0 {
				elapsed := time.Since(curTime)
				util.Logger(fmt.Sprintf("=========>>>>>> Iteration# %d Time: %s CLIENT-1", counter, elapsed))
				curTime = time.Now()
			}
			//invokeChaincode("PEER0") //For Local testing
			//invokeChaincode("vp0") //For Z -Testing
			invokeChaincode(util.GetPeer(0))
		}
		wg.Done()
	}()
	go func() {
		for i := 1; i <= TRX_COUNT/CLIENTS; i++ {
			if counter > 0 && counter%1000 == 0 {
				elapsed := time.Since(curTime)
				util.Logger(fmt.Sprintf("=========>>>>>> Iteration# %d Time: %s CLIENT-2", counter, elapsed))
				curTime = time.Now()
			}
			//invokeChaincode("PEER1") //For Local testing
			//invokeChaincode("vp1") //For Z -Testing
			invokeChaincode(util.GetPeer(1))
		}
		wg.Done()
	}()
	go func() {
		for i := 1; i <= TRX_COUNT/CLIENTS; i++ {
			if counter > 0 && counter%1000 == 0 {
				elapsed := time.Since(curTime)
				util.Logger(fmt.Sprintf("=========>>>>>> Iteration# %d Time: %s CLIENT-3", counter, elapsed))
				curTime = time.Now()
			}
			//invokeChaincode("PEER2") //For Local testing
			//invokeChaincode("vp2") //For Z -Testing
			invokeChaincode(util.GetPeer(2))
		}
		wg.Done()
	}()
	go func() {
		for i := 1; i <= TRX_COUNT/CLIENTS; i++ {
			if counter > 0 && counter%1000 == 0 {
				elapsed := time.Since(curTime)
				util.Logger(fmt.Sprintf("=========>>>>>> Iteration# %d Time: %s CLIENT-4", counter, elapsed))
				curTime = time.Now()
			}
			//invokeChaincode("PEER3") //For Local testing
			//invokeChaincode("vp3") //For Z -Testing
			invokeChaincode(util.GetPeer(3))

		}
		wg.Done()
	}()
}

//Execution starts here ...
func main() {
	util.InitLogger("LedgerStressFourCliFourPeer")
	//TODO:Add support similar to GNU getopts, http://goo.gl/Cp6cIg
	if len(os.Args) < 1 {
		util.Logger("Usage: go run LedgerStressFourCliFourPeer.go ")
		return
	}
	//TODO: Have a regular expression to check if the give argument is correct format
	/*if !strings.Contains(os.Args[1], "http://") {
		fmt.Println("Error: Argument submitted is not right format ex: http://127.0.0.1:5000 ")
		return;
	}
	//Get the URL
	url := os.Args[1]*/

	// time to messure overall execution of the testcase
	defer util.TimeTracker(time.Now(), "Total execution time for LedgerStressFourCliFourPeer.go ")

	Init()
	util.Logger("========= Transacations execution stated  =========")
	InvokeLoop()
	wg.Wait()
	util.Logger("========= Transacations execution ended  =========")
	util.TearDown(counter)
}