package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// NewClient - init new storage
func NewClient() *ClientInstance {
	tr := &http.Transport{
		MaxIdleConnsPerHost: 1024,
		TLSHandshakeTimeout: 50 * time.Second,
	}
	return &ClientInstance{client: &http.Client{Timeout: 50 * time.Second, Transport: tr}}
}

// ClientInstance - Инстанс клиента
type ClientInstance struct {
	client *http.Client
}

// GetJSON - get JSON
func (ci *ClientInstance) GetJSON(url string, target interface{}) error {
	r, err := ci.client.Get(url)
	if err != nil {
		return err
	}

	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("Couldn't parse response body. %+v", err)
	}

	return json.NewDecoder(bytes.NewReader(body)).Decode(target)
}

// IsExist - Проверка существования ссылки
func (ci *ClientInstance) IsExist(url string) error {
	r, err := ci.client.Head(url)
	if err != nil {
		return err
	}

	if r.StatusCode == 404 {
		return errors.New("URL not exist enimore")
	}

	return nil
}
