package httpx

import (
	"bytes"
	"context"
	"github.com/okcredit/go-common/encoding/json"
	"github.com/okcredit/go-common/errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

type ApiStub struct {
	BaseUrl    string
	HttpClient *http.Client
	Logger     *log.Logger

	once sync.Once
}

func (s *ApiStub) CallApi(ctx context.Context, httpMethod string, path string, reqBody interface{}, resBody interface{}) error {
	s.once.Do(func() {
		if s.HttpClient == nil {
			s.HttpClient = http.DefaultClient
		}
	})

	// request
	httpReqBody, err := s.buildRequestBody(reqBody)
	if err != nil {
		s.Logger.Printf("failed to json encode request body (method=%s, path=%s): %v", httpMethod, path, err)
		return err
	}

	httpReq, err := http.NewRequest(httpMethod, s.BaseUrl+path, httpReqBody)
	if err != nil {
		s.Logger.Printf("failed to create http request (method=%s, path=%s): %v", httpMethod, path, err)
		return err
	}
	httpReq = httpReq.WithContext(ctx)

	// call api
	httpRes, err := s.HttpClient.Do(httpReq)
	if err != nil {
		s.Logger.Printf("http request failed (method=%s, path=%s): %v", httpMethod, path, err)
		return err
	}

	// read res body
	defer httpRes.Body.Close()
	resBodyData, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		s.Logger.Printf("failed to read response body: %v", err)
		resBodyData = make([]byte, 0)
	}

	if len(resBodyData) > 0 {
		s.Logger.Printf("(%s %s) response body = %s", httpMethod, path, resBodyData)
	} else {
		s.Logger.Printf("(%s %s) response body = nil", httpMethod, path)
	}

	// response
	if (httpRes.StatusCode / 100) == 2 {
		// success response

		if len(resBodyData) > 0 {

			// try to json decode body
			if err := json.Unmarshal(resBodyData, resBody); err != nil {
				s.Logger.Printf("failed to json decode response body (method=%s, path=%s): %v", httpMethod, path, err)
				return err
			}

		}
		return nil

	} else {
		// failure response
		e := errors.From(httpRes.StatusCode, "no_body")

		if len(resBodyData) > 0 {

			// try to json decode error body
			if err := json.Unmarshal(resBodyData, &e); err != nil {
				s.Logger.Printf("failed to json decode error body: %v", err)
				e = errors.From(httpRes.StatusCode, string(resBodyData))
			}

		}

		return e
	}
}

func (s *ApiStub) buildRequestBody(reqBody interface{}) (io.Reader, error) {
	if reqBody == nil {
		// no request body; don't encode
		return http.NoBody, nil
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	if len(data) > 0 {
		s.Logger.Printf("request body = %s", data)
		return bytes.NewReader(data), nil
	} else {
		return http.NoBody, nil
	}
}
