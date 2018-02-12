package controllers

import (
	"encoding/json"
	"net/http"
)

func (app *Application) UpdateHandler(w http.ResponseWriter, r *http.Request) {
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

	var data []CarData
	json.Unmarshal([]byte(blockData), &data)

	returnData := &struct {
		TransactionId        string
		Success              bool
		Response             bool
		ResponseData         []CarData
		TransactionRequested string
		TransactionUpdated   string
		QueryData            Car
		SearchKey            string
	}{
		TransactionId:        "",
		Success:              false,
		Response:             false,
		ResponseData:         data,
		TransactionRequested: "true",
		TransactionUpdated:   "false",
	}
	// Query Single Record
	if r.FormValue("requested") == "true" {
		// Retrieving Single Query
		QueryValue := r.FormValue("carKeySearch")
		blockData, _, _ := app.Fabric.QueryOne(QueryValue)
		var queryResponse Car
		json.Unmarshal([]byte(blockData), &queryResponse)
		returnData.TransactionRequested = "false"
		returnData.TransactionUpdated = "true"
		returnData.SearchKey = QueryValue
		returnData.QueryData = queryResponse
	}
	// Update Single Record
	if r.FormValue("updated") == "true" {
		/* Form Data */
		carData := Car{}
		carKey := r.FormValue("carKey")
		carData.Make = r.FormValue("carMake")
		carData.Model = r.FormValue("carModel")
		carData.Colour = r.FormValue("carColor")
		carData.Owner = r.FormValue("carOwner")

		RequestData, _ := json.Marshal(carData)
		txid, err := app.Fabric.UpdateCarRecord(carKey, string(RequestData))

		if err != nil {
			http.Error(w, "Unable to update record in the blockchain", 500)
		}
		returnData.TransactionId = txid
		returnData.Success = true
		returnData.Response = true
		returnData.TransactionRequested = "true"
		returnData.TransactionUpdated = "false"
	}

	renderTemplate(w, r, "update.html", returnData)
}
