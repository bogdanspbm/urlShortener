package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"urlShortener/utils"
	_ "urlShortener/utils"
)

func getAllShortURLs() {

}

func getLinks() []string {
	var res []string

	file, err := os.Open("links.txt")

	if err != nil {
		return res
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		res = append(res, scanner.Text())
	}

	return res
}

func main() {

	links := getLinks()

	for i := 0; i < 10000; i++ {
		link := links[rand.Intn(len(links))]
		url := utils.GetURL{URL: link}
		jsonResp, _ := json.Marshal(url)
		req, _ := http.NewRequest("PUT", "http://127.0.0.1:8000/create", bytes.NewBuffer(jsonResp))

		client := &http.Client{}
		_, err := client.Do(req)

		if err != nil {
			fmt.Println("Could not do put request", err)
			return
		}
	}

	for i := 0; i < 100000; i++ {
		_, err := http.Get("http://127.0.0.1:8000/2qNfHK5")
		if err != nil {
			fmt.Println("Could not do get request", err)
			return
		}
	}

}
