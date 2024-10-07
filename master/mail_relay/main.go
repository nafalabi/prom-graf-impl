package main

import (
	"log"
	"mailgun-relay/sender"
	templt "mailgun-relay/template"
	"mailgun-relay/utils"
	"net/http"

	_ "embed"
)

var logger = utils.NewLogger(utils.INFO)

func alertHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("Starting sending mail alerts...")

	sendError := func(err error) {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	payload, err := utils.GetPayloadData(r.Body)
	if err != nil {
		sendError(err)
		return
	}

	subject := "Prometheus Server Alerts"

	var config sender.Config
	err = utils.ReadAndUnmarshal("config.json", &config)
	if err != nil {
		sendError(err)
		return
	}

	err = sender.SendMail(config, subject, payload)
	if err != nil {
		sendError(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	logger.Info("Alert successfully processed")
}

func webtestHandler(w http.ResponseWriter, r *http.Request) {
	sendError := func(err error) {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	payload := utils.DataPayload{}
	err := utils.ReadAndUnmarshal("webhook_req.json", &payload)
	if err != nil {
		sendError(err)
		return
	}

	tmpl := templt.GetTemplate()
	err = tmpl.Execute(w, payload)
	if err != nil {
		sendError(err)
		return
	}
}

func main() {
	http.HandleFunc("/alert", alertHandler)
	http.HandleFunc("/webtest", webtestHandler)

	logger.Info("Starting webhook server on :5001")
	if err := http.ListenAndServe(":5001", nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
