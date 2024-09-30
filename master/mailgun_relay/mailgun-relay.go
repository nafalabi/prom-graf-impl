package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	mailgun "github.com/mailgun/mailgun-go/v4"
)

type EnvVar struct {
	domain string
	apiKey string
	sender string
}

func (env *EnvVar) loadEnv() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	env.domain = os.Getenv("MAILGUN_DOMAIN")
	env.apiKey = os.Getenv("MAILGUN_API_KEY")
	env.sender = os.Getenv("SENDER_MAIL_ADDRESS")
}

type Alert struct {
	Status      string            `json:"status"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	StartsAt    string            `json:"startsAt"`
	EndsAt      string            `json:"endsAt"`
}

type AlertmanagerPayload struct {
	Alerts []Alert `json:"alerts"`
}

type Receiver = map[string][]string

func getReceivers() (Receiver, error) {
	var receivers Receiver

	file, err := os.Open("receivers.json")
	if err != nil {
		return receivers, fmt.Errorf("Error opening file: %v", err.Error())
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return receivers, fmt.Errorf("Error reading file: %v", err.Error())
	}

	if err := json.Unmarshal(bytes, &receivers); err != nil {
		return receivers, fmt.Errorf("Error decoding JSON: %v", err.Error())
	}

	return receivers, nil
}

func formatAlertBody(alerts []Alert) string {
	var body strings.Builder
	for _, alert := range alerts {
		body.WriteString(fmt.Sprintf("**Alert Name:** %s\n", alert.Labels["alertname"]))
		body.WriteString(fmt.Sprintf("**Status:** %s\n", alert.Status))
		body.WriteString(fmt.Sprintf("**Starts At:** %s\n", alert.StartsAt))
		body.WriteString(fmt.Sprintf("**Ends At:** %s\n", alert.EndsAt))

		body.WriteString("\n**Labels:**\n")
		for k, v := range alert.Labels {
			body.WriteString(fmt.Sprintf("  - %s: %s\n", k, v))
		}

		body.WriteString("\n**Annotations:**\n")
		for k, v := range alert.Annotations {
			body.WriteString(fmt.Sprintf("  - %s: %s\n", k, v))
		}

		body.WriteString("\n----------------------------------------\n")
	}
	return body.String()
}

func sendMail(subject string, recipients []string, alerts []Alert) error {
	envVar := EnvVar{}
	envVar.loadEnv()

	mg := mailgun.NewMailgun(envVar.domain, envVar.apiKey)

	body := formatAlertBody(alerts)

	message := mg.NewMessage(envVar.sender, subject, body)
	for _, recipient := range recipients {
		err := message.AddRecipient(recipient)
		if err != nil {
			continue
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, _, err := mg.Send(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to send mail: %v", err)
	}

	log.Printf("Mail sent successfully to %v", recipients)
	return nil
}

func alertHandler(w http.ResponseWriter, r *http.Request) {
	var payload AlertmanagerPayload

	body, err := io.ReadAll(r.Body)
	if err != nil {
    log.Printf("Could not read request body: %v", err)
		http.Error(w, "Could not read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &payload)
	if err != nil {
    log.Printf("Could not read request body: %v", err)
		http.Error(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
		return
	}

	subject := "Prometheus AlertManager - "
	if len(payload.Alerts) > 0 {
		subject += payload.Alerts[0].Labels["alertname"] + " - " + payload.Alerts[0].Status
	} else {
		subject += "No Alerts"
	}

	receivers, error := getReceivers()
	if error != nil {
    log.Printf("Error getting list of receivers: %v", error)
		return
	}

	var recipients []string
	for groupName, addresses := range receivers {
		log.Printf("Adding receivers from group '%s' ...", groupName)
		recipients = append(recipients, addresses...)
	}

	err = sendMail(subject, recipients, payload.Alerts)
	if err != nil {
		log.Printf("Error sending mail: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Alert processed")
}

func main() {
	http.HandleFunc("/alert", alertHandler)

	log.Println("Starting webhook server on :5001")
	if err := http.ListenAndServe(":5001", nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
