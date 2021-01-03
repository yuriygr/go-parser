package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// NewClient - init new storage
func NewClient() *Client {
	tr := &http.Transport{
		MaxIdleConnsPerHost: 1024,
		TLSHandshakeTimeout: 50 * time.Second,
	}
	return &Client{client: &http.Client{Timeout: 50 * time.Second, Transport: tr}}
}

// Client - Инстанс клиента
type Client struct {
	client *http.Client
}

// GetJSON - get JSON
func (c *Client) GetJSON(url string, target interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 Safari/537.36")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("Referer", url)

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Couldn't parse response body. %+v", err)
	}

	return json.NewDecoder(bytes.NewReader(body)).Decode(target)
}

// IsExist - Проверка существования ссылки
func (c *Client) IsExist(url string) error {
	r, err := c.client.Head(url)
	if err != nil {
		return err
	}

	if r.StatusCode == 404 {
		return errors.New("URL not exist enimore")
	}

	if r.StatusCode == 502 {
		return errors.New("Сервер не отдает нам видео, зачем оно нам?")
	}

	return nil
}
