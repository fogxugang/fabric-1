/*
 Copyright IBM Corp. 2016-2017 All Rights Reserved.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package chaincode

import (
	"fmt"
	"os"
	"testing"

	"github.com/hyperledger/fabric/peer/common"
	pb "github.com/hyperledger/fabric/protos/peer"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func initInstallTest(fsPath string, t *testing.T) *cobra.Command {
	viper.Set("peer.fileSystemPath", fsPath)
	finitInstallTest(fsPath)

	//if mkdir fails everthing will fail... but it should not
	if err := os.Mkdir(fsPath, 0755); err != nil {
		t.Fatalf("could not create install env")
	}

	InitMSP()

	signer, err := common.GetDefaultSigner()
	if err != nil {
		t.Fatalf("Get default signer error: %v", err)
	}

	mockCF := &ChaincodeCmdFactory{
		Signer: signer,
	}

	cmd := installCmd(mockCF)
	AddFlags(cmd)

	return cmd
}

func finitInstallTest(fsPath string) {
	os.RemoveAll(fsPath)
}

// TestBadVersion tests generation of install command
func TestBadVersion(t *testing.T) {
	fsPath := "/tmp/installtest"

	cmd := initInstallTest(fsPath, t)
	defer finitInstallTest(fsPath)

	args := []string{"-n", "example02", "-p", "github.com/hyperledger/fabric/examples/chaincode/go/chaincode_example02"}
	cmd.SetArgs(args)

	if err := cmd.Execute(); err == nil {
		t.Fatalf("Expected error executing install command for version not specified")
	}
}

// TestNonExistentCC non existent chaincode should fail as expected
func TestNonExistentCC(t *testing.T) {
	fsPath := "/tmp/installtest"

	cmd := initInstallTest(fsPath, t)
	defer finitInstallTest(fsPath)

	args := []string{"-n", "badexample02", "-p", "github.com/hyperledger/fabric/examples/chaincode/go/bad_example02", "-v", "testversion"}
	cmd.SetArgs(args)

	if err := cmd.Execute(); err == nil {
		t.Fatalf("Expected error executing install command for bad chaincode")
	}

	if _, err := os.Stat(fsPath + "/chaincodes/badexample02.testversion"); err == nil {
		t.Fatalf("chaincode example02.testversion should not exist")
	}
}

func installEx02() error {

	signer, err := common.GetDefaultSigner()
	if err != nil {
		return fmt.Errorf("Get default signer error: %v", err)
	}

	mockResponse := &pb.ProposalResponse{
		Response:    &pb.Response{Status: 200},
		Endorsement: &pb.Endorsement{},
	}

	mockEndorerClient := common.GetMockEndorserClient(mockResponse, nil)

	mockCF := &ChaincodeCmdFactory{
		EndorserClient: mockEndorerClient,
		Signer:         signer,
	}

	cmd := installCmd(mockCF)
	AddFlags(cmd)

	args := []string{"-n", "example02", "-p", "github.com/hyperledger/fabric/examples/chaincode/go/chaincode_example02", "-v", "anotherversion"}
	cmd.SetArgs(args)

	if err := cmd.Execute(); err != nil {
		return fmt.Errorf("Run chaincode upgrade cmd error:%v", err)
	}

	return nil
}

func TestInstall(t *testing.T) {
	InitMSP()
	if err := installEx02(); err != nil {
		t.Fatalf("Install failed with error: %v", err)
	}
}
