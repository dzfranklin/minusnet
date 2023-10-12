package hubs

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go.uber.org/zap"
)

func (h *Hub) RequestStatus() (Status, error) {
	client := &http.Client{}
	var statusText string
	var err error
	for retries := 0; retries < 3; retries++ {
		statusText, err = h.requestStatusOnce(client)
		if err != nil {
			zap.S().Errorf("getStatusOnce failed: %v", err)
			time.Sleep(10 * time.Second)
			continue
		}
		break
	}
	if err != nil {
		return Status{}, err
	}
	return parseStatus(statusText)
}

func (h *Hub) requestStatusOnce(client *http.Client) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	ts := time.Now().Format("02_01_2006")
	name := "/status_log/" + h.Serial + "_" + ts + "_status.csv"
	reqBody := strings.NewReader("url=" + url.QueryEscape(name))
	req, err := http.NewRequestWithContext(ctx, "POST", h.endpoint()+name, reqBody)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/115.0")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", h.endpoint()+"/basic_-_status.htm")
	req.Header.Set("Origin", h.endpoint())
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(respBody), nil
}
