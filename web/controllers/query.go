package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (app *Application) QueryHandler(w http.ResponseWriter, r *http.Request) {

	QueryValue := r.FormValue("car")
	fmt.Println(QueryValue)
	blockData, txnID, err := app.Fabric.QueryOne(QueryValue)

	fmt.Println("#### Query One ###")
	fmt.Printf("%v", blockData)

	if err != nil {
		http.Error(w, "Unable to query the blockchain", 500)
	}

	type CarData struct {
		Make   string `json:"make"`
		Model  string `json:"model"`
		Colour string `json:"colour"`
		Owner  string `json:"owner"`
	}

	var data CarData
	json.Unmarshal([]byte(blockData), &data)

	returnData := &struct {
		ResponseData  CarData
		TransactionID string
	}{
		ResponseData:  data,
		TransactionID: txnID,
	}
	returnData.TransactionID = txnID

	fmt.Println("######## ResponseData")
	fmt.Printf("%v", returnData)

	renderTemplate(w, r, "query.html", returnData)
}
