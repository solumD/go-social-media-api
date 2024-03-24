package client

import (
	"bytes"
	"log"
	"net/http"
)

func main() {
	client := &http.Client{}
	b := []byte{}
	request, err := http.NewRequest(http.MethodGet, "https://dzen.ru/news", bytes.NewReader(b))
	if err != nil {
		log.Println(err)
	}
	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
	}
	log.Println(response)
}
