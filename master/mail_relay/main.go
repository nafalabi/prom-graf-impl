package main

import (
	"fmt"
	"log"
	"mailgun-relay/sender"
	templt "mailgun-relay/template"
	"mailgun-relay/utils"
	"net/http"

	_ "embed"

	"github.com/go-chi/chi/v5"
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

	var configMap sender.ConfigMap
	err = utils.ReadAndUnmarshal("config.json", &configMap)
	if err != nil {
		sendError(err)
		return
	}

	configKey := chi.URLParam(r, "configKey")
	senderConfig, ok := configMap[configKey]
	if !ok {
		sendError(fmt.Errorf("Can't get the sender config for '%v'", configKey))
		return
	}

	subject := "Prometheus Server Alerts"

	err = sender.SendMail(senderConfig, subject, payload)
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
	router := chi.NewRouter()

	router.Post("/alert/{configKey}", alertHandler)
	router.Get("/webtest", webtestHandler)

	logger.Info("Starting webhook server on :5001")
	if err := http.ListenAndServe(":5001", router); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
