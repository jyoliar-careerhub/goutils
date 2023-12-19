package apiactor

import (
	"fmt"
	"goutils/cchan"
	"io"
	"log"
	"math"
	"net/http"
	"time"
)

type ApiActor struct {
	rjChan    chan *requestJob
	minDelay  int64
	isStarted bool
}

func NewApiActor(minDelay int64) *ApiActor {
	apiActor := ApiActor{
		rjChan:   make(chan *requestJob),
		minDelay: minDelay,
	}

	return &apiActor
}

type Request struct {
	Method string
	Url    string
	Header http.Header
}

func NewRequest(method string, url string) *Request {
	r := &Request{
		Method: method,
		Url:    url,
		Header: make(http.Header),
	}

	return r
}

type requestJob struct {
	request    *http.Request
	resultChan chan<- *result
}

type result struct {
	reader io.ReadCloser
	err    error
}

func Run[QUIT any](a *ApiActor, quitChan <-chan QUIT) {
	if a.isStarted {
		log.Fatal("ApiActor is already started")
	}

	go run(a, quitChan)
	a.isStarted = true
}

func run[QUIT any](a *ApiActor, quitChan <-chan QUIT) {
	var lastEndedTime int64 = math.MaxInt64

	for {
		timeDiff := time.Now().UnixMilli() - lastEndedTime

		if timeDiff < a.minDelay && timeDiff > 0 {
			time.Sleep(time.Millisecond * time.Duration(a.minDelay-timeDiff))
		}

		lastEndedTime = time.Now().UnixMilli()

		rj, ok := cchan.ReceiveOrQuit(a.rjChan, quitChan)
		if !ok {
			for {
				select {
				case reqJob := <-a.rjChan:
					close((*reqJob).resultChan) //대기하고 있을 다른 goroutine을 위해 resultChan을 닫아준다.
				default:
					return
				}
			}
		}

		rc, err := callApi((*rj).request)

		ok = cchan.SendOrQuit(
			&result{
				reader: rc,
				err:    err,
			},
			(*rj).resultChan,
			quitChan,
		)
		if !ok {
			close((*rj).resultChan)
			return
		}
	}
}

func (a *ApiActor) Call(r *Request) (io.ReadCloser, error) {
	if !a.isStarted {
		panic("ApiActor is not started")
	}

	httpReq, err := converthttpReq(r)
	if err != nil {
		return nil, err
	}

	resultChan := make(chan *result)
	a.rjChan <- &requestJob{
		request:    httpReq,
		resultChan: resultChan,
	}

	result, ok := <-resultChan
	if !ok {
		return nil, fmt.Errorf("resultChan closed")
	}

	close(resultChan)
	return result.reader, result.err
}

func converthttpReq(r *Request) (*http.Request, error) {
	httpReq, err := http.NewRequest(r.Method, r.Url, nil)

	if err != nil {
		return nil, err
	}

	httpReq.Header = r.Header

	return httpReq, nil
}
