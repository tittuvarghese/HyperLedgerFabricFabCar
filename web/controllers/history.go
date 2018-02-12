package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (app *Application) HistoryHandler(w http.ResponseWriter, r *http.Request) {
	blockData, err := app.Fabric.QueryAll()
	if err != nil {
		http.Error(w, "Unable to query the blockchain", 500)
	}

	type Car struct {
		Make   string `json:"make"`
		Model  string `json:"model"`
		Colour string `json:"colour"`
		Owner  string `json:"owner"`
	}

	type CarData struct {
		Key    string `json:"key"`
		Record Car    `json:"record"`
	}

	type RecordHistory struct {
		TxId      string `json:"TxId"`
		Value     Car    `json:"Value"`
		Timestamp string `json:"Timestamp"`
		IsDelete  string `json:"IsDelete"`
	}

	var data []CarData
	json.Unmarshal([]byte(blockData), &data)

	returnData := &struct {
		ResponseData         []CarData
		TransactionRequested string
		TransactionUpdated   string
		RecordHistory        []RecordHistory
	}{
		ResponseData:         data,
		TransactionRequested: "true",
	}
	// Query History Using Key
	if r.FormValue("requested") == "true" {
		// Retrieving Single Query
		QueryValue := r.FormValue("carKeySearch")
		blockHistory, _ := app.Fabric.GetHistoryofCar(QueryValue)
		var queryResponse []RecordHistory
		json.Unmarshal([]byte(blockHistory), &queryResponse)
		returnData.RecordHistory = queryResponse
		returnData.TransactionRequested = "true"
		fmt.Println("### Response History ###")
		fmt.Printf("%s", blockHistory)
	}
	renderTemplate(w, r, "history.html", returnData)
}
