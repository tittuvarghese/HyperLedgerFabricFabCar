package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/servntire/servntire-demo/blockchain"
)

// Fix empty GOPATH with golang 1.8 (see https://github.com/golang/go/blob/1363eeba6589fca217e155c829b2a7c00bc32a92/src/go/build/build.go#L260-L277)
func defaultGOPATH() string {
	env := "HOME"
	if runtime.GOOS == "windows" {
		env = "USERPROFILE"
	} else if runtime.GOOS == "plan9" {
		env = "home"
	}
	if home := os.Getenv(env); home != "" {
		def := filepath.Join(home, "go")
		if filepath.Clean(def) == filepath.Clean(runtime.GOROOT()) {
			// Don't set the default GOPATH to GOROOT,
			// as that will trigger warnings from the go tool.
			return ""
		}
		return def
	}
	return ""
}

func main() {
	// Setup correctly the GOPATH in the environment
	if goPath := os.Getenv("GOPATH"); goPath == "" {
		os.Setenv("GOPATH", defaultGOPATH())
	}

	// Initialize the Fabric SDK
	fabricSdk, err := blockchain.Initialize()
	if err != nil {
		fmt.Printf("Unable to initialize the Fabric SDK: %v", err)
	}

	// Install and instantiate the chaincode
	err = fabricSdk.InstallAndInstantiateCC()
	if err != nil {
		fmt.Printf("Unable to install and instantiate the chaincode: %v\n", err)
	}

	// Query the chaincode
	fmt.Println("###### Query All ######")
	response, err := fabricSdk.QueryAll()
	if err != nil {
		fmt.Printf("Unable to query all records from the chaincode: %v\n", err)
	} else {
		fmt.Printf("Response from the queryAll: %v\n", response)
	}

	fmt.Println("###### Query One ######")
	response, _, err = fabricSdk.QueryOne("CAR4")
	if err != nil {
		fmt.Printf("Unable to query one from chaincode: %v\n", err)
	} else {
		fmt.Printf("Response from the queryOne: %v\n", response)
	}

	// Create New Car
	type Car struct {
		Make   string `json:"make"`
		Model  string `json:"model"`
		Colour string `json:"colour"`
		Owner  string `json:"owner"`
	}
	carData := Car{}
	carKey := "CAR10"
	carData.Make = "Volkswagen"
	carData.Model = "Vento"
	carData.Colour = "grey"
	carData.Owner = "Mohan"

	RequestData, _ := json.Marshal(carData)
	txId, err := fabricSdk.CreateCar(carKey, string(RequestData))

	if err != nil {
		fmt.Printf("Unable to create record on the chaincode: %v\n", err)
	} else {
		fmt.Printf("Successfully created record, transaction ID: %s\n", txId)
	}

	// Query the chaincode Again to retrieve updated data
	fmt.Println("###### Query All ######")
	response, err = fabricSdk.QueryAll()
	if err != nil {
		fmt.Printf("Unable to query all from the chaincode: %v\n", err)
	} else {
		fmt.Printf("Response from the queryAll: %v\n", response)
	}

	// Changing Car Owner by Passing Key and Value
	fmt.Println("###### Change Ownership ######")
	txId, err = fabricSdk.ChangeCarOwner("CAR10", "Keyana")
	if err != nil {
		fmt.Printf("Unable to invoke - Change Car Owner on the chaincode: %v\n", err)
	} else {
		fmt.Printf("Successfully invoke - Change Car Owner, transaction ID: %s\n", txId)
	}

}
