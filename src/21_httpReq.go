package main

import (
	"fmt"
	"net/http"
	"strings"
)

func main() {

	for i := 19999; i < 29999; i++ {
		targetUrl := fmt.Sprintf("http://127.0.0.1:12380/my-key%d", i)


		payload := strings.NewReader(fmt.Sprintf("my-value%d", i))

		req, _ := http.NewRequest("PUT", targetUrl, payload)

		req.Header.Add("Content-Type", "application/json")

		response, err := http.DefaultClient.Do(req)

		if err != nil {
			panic(err)
		}
		response.Body.Close()
	}




}
