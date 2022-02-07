package yalohttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	DefaultURL  = "https://api-staging2.yalochat.com/awesome-bank/v1/messages"
	BearerToken = "token-here"
)

type Client struct {
	http.Client
	Request *http.Request
	//retrialLimit int
}

func NewClient() (*Client, error) {
	client := &Client{}

	req, err := http.NewRequest("", DefaultURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", BearerToken))
	req.Header.Add("Content-type", fmt.Sprintf("application/json"))

	client.Request = req
	return client, nil
}

func (c *Client) Post(body interface{}) (*http.Response, error) {
	retrialLimit := 0
	c.Request.Method = http.MethodPost
	c.Request.Body = serializeBodyToJson(body)

	t := time.Now()
	resp, err := c.Client.Do(c.Request)
	_ = time.Since(t)
	//fmt.Println(ta)

	for (err != nil || checkRetries(resp.StatusCode)) && retrialLimit < 5 {
		resp, err = c.Client.Do(c.Request)
		retrialLimit++
		fmt.Println("retrialLimit ->", retrialLimit)
		xb, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("retrialCause ->", string(xb), resp.Status)
	}

	retrialLimit = 0
	fmt.Println(resp.StatusCode)
	return resp, nil
}

func check2xx(statusCode int) bool {
	if statusCode >= 200 && statusCode < 300 {
		return true
	}

	return false
}

func checkRetries(statusCode int) bool {
	return statusCode == 429
}

func serializeBodyToJson(item interface{}) io.ReadCloser {
	xb, err := json.Marshal(&item)
	if err != nil {
		fmt.Printf("error marshling the item: %v\n", err)
	}

	return ioutil.NopCloser(bytes.NewReader(xb))
}
