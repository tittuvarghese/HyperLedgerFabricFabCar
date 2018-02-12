package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// ServntireDemoChaincode implementation of Chaincode
type ServntireDemoChaincode struct {
}
type Car struct {
	Make   string `json:"make"`
	Model  string `json:"model"`
	Colour string `json:"colour"`
	Owner  string `json:"owner"`
}

// Init of the chaincode
// This function is called only one when the chaincode is instantiated.
// So the goal is to prepare the ledger to handle future requests.
func (t *ServntireDemoChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("########### ServntireDemoChaincode Init ###########")

	// Get the function and arguments from the request
	function, _ := stub.GetFunctionAndParameters()

	// Check if the request is the init function
	if function != "init" {
		return shim.Error("Unknown function call")
	}

	// Adding to the ledger.
	cars := []Car{
		Car{Make: "Toyota", Model: "Prius", Colour: "blue", Owner: "Tomoko"},
		Car{Make: "Ford", Model: "Mustang", Colour: "red", Owner: "Brad"},
		Car{Make: "Hyundai", Model: "Tucson", Colour: "green", Owner: "Jin Soo"},
		Car{Make: "Volkswagen", Model: "Passat", Colour: "yellow", Owner: "Max"},
		Car{Make: "Tesla", Model: "S", Colour: "black", Owner: "Adriana"},
		Car{Make: "Peugeot", Model: "205", Colour: "purple", Owner: "Michel"},
		Car{Make: "Chery", Model: "S22L", Colour: "white", Owner: "Aarav"},
		Car{Make: "Fiat", Model: "Punto", Colour: "violet", Owner: "Pari"},
		Car{Make: "Tata", Model: "Nano", Colour: "indigo", Owner: "Valeria"},
		Car{Make: "Holden", Model: "Barina", Colour: "brown", Owner: "Shotaro"},
	}

	i := 0
	for i < len(cars) {
		fmt.Println("i is ", i)
		carAsBytes, _ := json.Marshal(cars[i])
		stub.PutState("CAR"+strconv.Itoa(i), carAsBytes)
		fmt.Println("Added", cars[i])
		i = i + 1
	}

	// Return a successful message
	return shim.Success(nil)
}

// Invoke
// All future requests named invoke will arrive here.
func (t *ServntireDemoChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("########### ServntireDemoChaincode Invoke ###########")

	// Get the function and arguments from the request
	function, args := stub.GetFunctionAndParameters()

	// Check whether it is an invoke request
	if function != "invoke" {
		return shim.Error("Unknown function call")
	}

	// Check whether the number of arguments is sufficient
	if len(args) < 1 {
		return shim.Error("The number of arguments is insufficient.")
	}

	// In order to manage multiple type of request, we will check the first argument.
	// Here we have one possible argument: query (every query request will read in the ledger without modification)
	if args[0] == "query" {
		return t.query(stub, args)
	}

	// The update argument will manage all update in the ledger
	if args[0] == "invoke" {
		return t.invoke(stub, args)
	}

	// Querying Single Record by Passing CAR ID => Key as parameter
	if args[0] == "queryone" {
		return t.queryone(stub, args)
	}

	// If the arguments given don’t match any function, we return an error
	return shim.Error("Unknown action, check the first argument")
}

// query
// Every readonly functions in the ledger will be here
func (t *ServntireDemoChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	// Check whether the number of arguments is sufficient
	if len(args) < 2 {
		return shim.Error("The number of arguments is insufficient.")
	}

	// Retrieves All the records here.
	if args[1] == "all" {

		//startKey := "CAR0"
		//endKey := "CAR999"

		// GetState by passing lower and upper limits
		resultsIterator, err := stub.GetStateByRange("", "")
		if err != nil {
			return shim.Error(err.Error())
		}
		defer resultsIterator.Close()

		// buffer is a JSON array containing QueryResults
		var buffer bytes.Buffer
		buffer.WriteString("[")

		bArrayMemberAlreadyWritten := false
		for resultsIterator.HasNext() {
			queryResponse, err := resultsIterator.Next()
			if err != nil {
				return shim.Error(err.Error())
			}
			// Add a comma before array members, suppress it for the first array member
			if bArrayMemberAlreadyWritten == true {
				buffer.WriteString(",")
			}
			buffer.WriteString("{\"Key\":")
			buffer.WriteString("\"")
			buffer.WriteString(queryResponse.Key)
			buffer.WriteString("\"")

			buffer.WriteString(", \"Record\":")
			// Record is a JSON object, so we write as-is
			buffer.WriteString(string(queryResponse.Value))
			buffer.WriteString("}")
			bArrayMemberAlreadyWritten = true
		}
		buffer.WriteString("]")

		fmt.Printf("- queryAllCars:\n%s\n", buffer.String())

		return shim.Success(buffer.Bytes())
	}

	// If the arguments given don’t match any function, we return an error
	return shim.Error("Unknown query action, check the second argument.")
}

//  Retrieves a single record from the ledger by accepting Key value
func (t *ServntireDemoChaincode) queryone(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 2 {
		return shim.Error("The number of arguments is insufficient.")
	}

	// GetState retrieves the data from ledger using the Key
	carAsBytes, _ := stub.GetState(args[1])

	// Transaction Response
	return shim.Success(carAsBytes)

}

func main() {
	// Start the chaincode and make it ready for futures requests
	err := shim.Start(new(ServntireDemoChaincode))
	if err != nil {
		fmt.Printf("Error starting Servntire Demo chaincode: %s", err)
	}
}
