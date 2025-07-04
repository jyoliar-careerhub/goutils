package apiactor

import (
	"fmt"
	"io"
	"net/http"
)

func CallApi(httpReq *http.Request) (io.ReadCloser, error) {
	client := &http.Client{}
	res, err := client.Do(httpReq)

	return GetBody(res, err)
}

func GetBody(res *http.Response, err error) (io.ReadCloser, error) {
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

func IsHttpError(err error) bool {
	_, ok := err.(*HttpError)
	return ok
}

func IsHttpErrorWithStatusCode(err error, statusCode int) bool {
	httpErr, ok := err.(*HttpError)
	return ok && httpErr.StatusCode == statusCode
}
