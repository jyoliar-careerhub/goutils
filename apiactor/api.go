package apiactor

import (
	"fmt"
	"io"
	"net/http"
)

func callApi(httpReq *http.Request) (io.ReadCloser, error) {
	client := &http.Client{}
	res, err := client.Do(httpReq)

	return getBody(res, err)
}

func getBody(res *http.Response, err error) (io.ReadCloser, error) {
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, &HttpError{res.StatusCode, res.Status}
	}

	if res.Body == nil {
		return nil, fmt.Errorf("response body is nil")
	}

	return res.Body, nil
}

type HttpError struct {
	StatusCode int
	Status     string
}

func (h *HttpError) Error() string {
	return fmt.Sprintf("HttpError: %s", h.Status)
}
