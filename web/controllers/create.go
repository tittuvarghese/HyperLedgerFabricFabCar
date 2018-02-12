package controllers

import (
	"encoding/json"
	"net/http"
)

type Car struct {
	Make   string `json:"make"`
	Model  string `json:"model"`
	Colour string `json:"colour"`
	Owner  string `json:"owner"`
}

func (app *Application) CreateHandler(w http.ResponseWriter, r *http.Request) {
	data := &struct {
		TransactionId string
		Success       bool
		Response      bool
	}{
		TransactionId: "",
		Success:       false,
		Response:      false,
	}
	if r.FormValue("submitted") == "true" {
		/* Form Data */
		carData := Car{}
		carKey := r.FormValue("carKey")
		carData.Make = r.FormValue("carMake")
		carData.Model = r.FormValue("carModel")
		carData.Colour = r.FormValue("carColor")
		carData.Owner = r.FormValue("carOwner")

		RequestData, _ := json.Marshal(carData)
		txid, err := app.Fabric.CreateCar(carKey, string(RequestData))

		if err != nil {
			http.Error(w, "Unable to create record in the blockchain", 500)
		}
		data.TransactionId = txid
		data.Success = true
		data.Response = true
	}
	renderTemplate(w, r, "create.html", data)
}
