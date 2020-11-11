package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

var (
	states = []string{"start", "stop"}
)

func validateState(s string) bool {
	for _, st := range states {
		if s != st {
			continue
		}
		return true
	}
	return false
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("incorrect number of args")
	}

	if os.Getenv("DISCORD_WEBHOOK") == "" {
		log.Fatal("missing DISCORD_WEBHOOK please set it to send payloads")
	}

	state := os.Args[1]

	if !validateState(state) {
		log.WithField("state", state).Fatalf("invalid state must be one of: %v", states)
	}

	pb, err := json.Marshal(generatePayload(state))
	if err != nil {
		log.WithField("err", err).Fatal("encountered err encoding payload")
	}

	resp, err := http.Post(os.Getenv("DISCORD_WEBHOOK"), "application/json", bytes.NewReader(pb))
	if err != nil {
		log.WithField("err", err).Fatal("encountered err sending payload")
	}

	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithField("err", err).Errorf("failed to read response")
	}
	log.WithField("statuscode", resp.StatusCode).WithField("body", string(r)).Infof("response came back")
}

type payload struct {
	Embeds []embed `json:"embeds,omitempty"`
}

type embed struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Color       int    `json:"color,omitempty"`
	Footer      struct {
		Text string `json:"text,omitempty"`
	} `json:"footer,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}

func generatePayload(state string) payload {
	p := payload{
		Embeds: []embed{
			{
				Title:       "Minecraft Server - %s",
				Description: "",
				Color:       0,
				Footer: struct {
					Text string "json:\"text,omitempty\""
				}{
					"Systemd Service",
				},
				Timestamp: time.Now(),
			},
		},
	}
	switch state {
	case "start":
		p.Embeds[0].Title = fmt.Sprintf(p.Embeds[0].Title, "Starting")
		p.Embeds[0].Color = 4437377
		p.Embeds[0].Description = "The minecraft server is starting up, please wait 5 minutes before trying to login."
	case "stop":
		p.Embeds[0].Title = fmt.Sprintf(p.Embeds[0].Title, "Stopping")
		p.Embeds[0].Color = 15730953
		p.Embeds[0].Description = "The minecraft server is shutting down, either by hand or by Amazon."
	}
	return p
}
