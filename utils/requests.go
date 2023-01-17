package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

var data urlData

var Client *Kafka

type urlData interface {
	store(key string, url string)
	loadKey(url string) (key string, ok bool)
	loadURL(key string) (url string, ok bool)
}

type GetURL struct {
	URL string `json:"longurl"`
}

type PingResponse struct {
	Topic string `json:"topic"`
	Type  string `json:"type"`
}

func InitData(db urlData) {
	data = db
}

func GetURLFromKey(key string) (string, bool) {
	url, ok := data.loadURL(key)
	if !ok {
		return "", ok
	}

	return url, ok
}

func askMasterForURL(key string) (string, bool) {
	kafkaMaster := Kafka{
		Topic: "bmadzhuga-events",
		Type:  "master",
	}

	err := kafkaMaster.Connect()

	if err != nil {
		return "", false
	}

	err = kafkaMaster.Send(Client.Topic, key, true)

	if err != nil {
		return "", false
	}

	response, err := Client.readFromTopic()

	if err != nil {
		return "", false
	}

	responseArray := strings.Split(response, "::")

	if len(responseArray) != 3 {
		return "", false
	}

	status := responseArray[2]
	url := responseArray[1]
	keyOut := responseArray[0]

	if status == "failed" {
		return "", false
	}

	if key != keyOut {
		return "", false
	}

	return url, true
}

func HandleGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		return
	}
	start := time.Now()

	promReceivedLinkCount.Inc()

	key := r.URL.Path[1:]

	var url string
	var ok bool

	if Client.Type == "master" {
		url, ok = GetURLFromKey(key)
	} else {
		url, ok = askMasterForURL(key)
	}

	if !ok {
		http.Error(w, "key not found", http.StatusNotFound)
		return
	}
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "http://")

	http.Redirect(w, r, "https://"+url, 301)

	duration := time.Since(start)

	requestProcessingTimeSummaryMs.Observe(duration.Seconds())
	requestProcessingTimeHistogramMs.Observe(duration.Seconds())

	PrometheusPush()
}

func HandlePing(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonResp := PingResponse{Topic: Client.Topic, Type: Client.Type}
	resp, err := json.Marshal(jsonResp)

	if err != nil {
		http.Error(w, "JSON is invalid", 400)
		return
	}

	w.Write(resp)
}

func createResp(w http.ResponseWriter, key string, url string) {
	jsonResp := UrlDB{Key: key, URL: url}
	resp, err := json.Marshal(jsonResp)
	if err != nil {
		http.Error(w, "JSON is invalid", 400)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func HandlePut(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPut && r.Method != http.MethodPost {
		return
	}

	fmt.Printf("PUT")

	start := time.Now()

	promRegisteredLinkCount.Inc()

	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	var jsonURL GetURL
	decoder.Decode(&jsonURL)

	url := jsonURL.URL

	if url == "" {
		w.WriteHeader(http.StatusOK)
		duration := time.Since(start)

		requestProcessingTimeSummaryMs.Observe(duration.Seconds())
		requestProcessingTimeHistogramMs.Observe(duration.Seconds())

		PrometheusPush()

		return
	}

	if key, ok := data.loadKey(url); ok {
		createResp(w, key, url)

		duration := time.Since(start)

		requestProcessingTimeSummaryMs.Observe(duration.Seconds())
		requestProcessingTimeHistogramMs.Observe(duration.Seconds())

		PrometheusPush()

		return
	}

	key, ok := "", true
	for ok {
		key = RandKey()
		_, ok = data.loadURL(key)
	}
	data.store(key, url)
	createResp(w, key, url)

	duration := time.Since(start)

	requestProcessingTimeSummaryMs.Observe(duration.Seconds())
	requestProcessingTimeHistogramMs.Observe(duration.Seconds())

	PrometheusPush()
}
