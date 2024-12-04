package datacontroller

import (
	"context"
	"fmt"
	"io"
	"net/http"

	httpErr "oreshnik/pkg/http/error"
)

type proxyParams struct {
	destHTTPMethod string
	destPath       string
	sourceReq      *http.Request
}

func (dc *DataController) proxyStaticRequestResponse(
	ctx context.Context, w http.ResponseWriter, params *proxyParams,
) {
	if err := dc.proxyRequestResponse(ctx, w, dc.staticAddr, params); err != nil {
		httpErr.InternalError(w, fmt.Errorf("proxy static request response: %w", err))
		return
	}
}

func (dc *DataController) proxyUsersRequestResponse(
	ctx context.Context, w http.ResponseWriter, params *proxyParams,
) {
	if err := dc.proxyRequestResponse(ctx, w, dc.usersAddr, params); err != nil {
		httpErr.InternalError(w, fmt.Errorf("proxy static request response: %w", err))
		return
	}
}

func (dc *DataController) proxyRequestResponse(
	ctx context.Context, w http.ResponseWriter, destURL string, params *proxyParams,
) error {
	resp, err := dc.proxyRequest(ctx, destURL, params)
	if err != nil {
		return fmt.Errorf("proxy request: %w", err)
	}
	defer resp.Body.Close()

	if err = dc.proxyResponse(w, resp); err != nil {
		return fmt.Errorf("proxy response: %w", err)
	}

	return nil
}

func (dc *DataController) proxyRequest(
	ctx context.Context, destURL string, params *proxyParams,
) (*http.Response, error) {
	req, err := http.NewRequestWithContext(
		ctx, params.destHTTPMethod, destURL+params.destPath, http.NoBody,
	)
	if err != nil {
		return nil, fmt.Errorf("new request with context: %w", err)
	}

	req.Header = params.sourceReq.Header
	req.Body = params.sourceReq.Body
	// in case if we create source request ourselves without URL
	if params.sourceReq.URL != nil {
		req.URL.RawQuery = params.sourceReq.URL.RawQuery
	}

	resp, err := dc.cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client do: %w", err)
	}

	return resp, nil
}

func (dc *DataController) proxyResponse(w http.ResponseWriter, resp *http.Response) error {
	for key, values := range resp.Header {
		for _, val := range values {
			// use Set to don't duplicate headers
			w.Header().Set(key, val)
		}
	}

	w.WriteHeader(resp.StatusCode)

	if _, err := io.Copy(w, resp.Body); err != nil {
		return fmt.Errorf("io copy response body: %w", err)
	}

	return nil
}
