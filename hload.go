package main

import (
	"fmt"
	"net/http"
)

func getAllShortURLs() {

}

func main() {

	/*
		for i := 0; i < 10000; i++ {
			resp, err := http.PostForm("http://example.com/form",
				url.Values{"longurl": {"Value"}})
		}*/

	for i := 0; i < 100000; i++ {
		_, err := http.Get("http://127.0.0.1:8000/2qNfHK5")
		if err != nil {
			fmt.Println("Could not do get request", err)
		}
	}

}
