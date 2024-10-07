package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	templt "mailgun-relay/template"
	"mailgun-relay/utils"
	"net/http"
	"time"

	_ "embed"
	mailgun "github.com/mailgun/mailgun-go/v4"
)

type Receivers = map[string][]string

//go:embed bawana.png
var bawanaIcon []byte

var logger = utils.NewLogger(utils.INFO)

func formatMail(data utils.DataPayload) (string, error) {
	var result string

	var buff bytes.Buffer
	tmpl := templt.GetTemplate()
	err := tmpl.Execute(&buff, data)
	if err != nil {
		return result, err
	}

	result = buff.String()

	return result, nil
}

func sendMail(subject string, recipients []string, body string) error {
	envVar := utils.NewEnv()

	mg := mailgun.NewMailgun(envVar.Domain, envVar.ApiKey)

	message := mg.NewMessage(envVar.Sender, subject, "")
	for _, recipient := range recipients {
		err := message.AddRecipient(recipient)
		if err != nil {
			logger.Error(err.Error())
			continue
		}
	}

	message.SetHtml(body)
  message.AddInline("./bawana.png")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, _, err := mg.Send(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to send mail: %v", err)
	}

	logger.Info(fmt.Sprintf("Mail sent successfully to %v", recipients))
	return nil
}

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

	subject := "Prometheus AlertManager - "
	if len(payload.Alerts) > 0 {
		subject += payload.Alerts[0].Labels["alertname"] + " - " + payload.Alerts[0].Status
	} else {
		subject += "No Alerts"
	}

	var receivers Receivers
	err = utils.ReadAndUnmarshal("receivers.json", &receivers)
	if err != nil {
		sendError(err)
		return
	}

	var recipients []string
	for groupName, mailAddresses := range receivers {
		logger.Info(fmt.Sprintf("Adding receivers from group '%s' ...", groupName))
		recipients = append(recipients, mailAddresses...)
	}

	mailBody, err := formatMail(payload)
	if err != nil {
		sendError(err)
		return
	}

	err = sendMail(subject, recipients, mailBody)
	if err != nil {
		sendError(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	logger.Info("Alert successfully processed")
}

func templateTest(w http.ResponseWriter, r *http.Request) {
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
	// err = tmpl.ExecuteTemplate(w, "email.default.html", payload)
	err = tmpl.Execute(w, payload)
	if err != nil {
		sendError(err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/alert", alertHandler)
	http.HandleFunc("/webtest", templateTest)

	log.Println("Starting webhook server on :5001")
	if err := http.ListenAndServe(":5001", nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
