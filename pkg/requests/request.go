package requests

import (
	"errors"
	"io"
	"net/http"
	"time"
)

type (
	httpMethod interface {
		Get(url string) ([]byte, error)
	}

	reqCli struct {
		client *http.Client
	}
)

func NewRequest() httpMethod {
	return &reqCli{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (h *reqCli) Get(url string) ([]byte, error) {

	response, err := h.client.Get(url)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, errors.New("error: server down")
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
