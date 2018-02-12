package controllers

import (
	"encoding/json"
	"net/http"
)

func (app *Application) HomeHandler(w http.ResponseWriter, r *http.Request) {
	blockData, err := app.Fabric.QueryAll()
	if err != nil {
		http.Error(w, "Unable to query the blockchain", 500)
	}

	type CarData struct {
		Key    string `json:"Key"`
		Record struct {
			Make   string `json:"make"`
			Model  string `json:"model"`
			Colour string `json:"colour"`
			Owner  string `json:"owner"`
		} `json:"Record"`
	}

	var data []CarData
	json.Unmarshal([]byte(blockData), &data)

	returnData := &struct {
		ResponseData []CarData
	}{
		ResponseData: data,
	}

	renderTemplate(w, r, "home.html", returnData)
}
