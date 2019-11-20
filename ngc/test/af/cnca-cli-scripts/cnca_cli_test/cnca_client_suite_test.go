package main_test

import (
	"testing"
	//"context"
	//"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	//af "github.com/otcshare/epcforedge/ngc/pkg/af/lib"
)

func TestGoApiTests(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CncaClient Suite")
}

/*const(
	cfgPath = "/root/go_projects/src/github.com/otcshare/epcforedge/ngc"
)

var (
	srvCtx		context.Context
	srvCancel	context.CancelFunc
	afIsRunning	bool
)

func runAf(stopIndication chan bool) error {

	By("Starting AF server")

	srvCtx, srvCancel = context.WithCancel(context.Background())
	_ = srvCancel
	afRunFail := make(chan bool)
	go func() {
		err := af.Run(srvCtx, cfgPath + "/configs/af.json")
		if err != nil {
			fmt.Printf("Run() exited with error: %#v", err)
			afIsRunning = false
			afRunFail <- true
		}
		stopIndication <- true
	}()

	return nil
}

func stopAf(stopIndication chan bool) int {
	By("Stopping AF server")
	srvCancel()
	<-stopIndication
	if afIsRunning == true {
		return 0
	}

	return 1
}*/
