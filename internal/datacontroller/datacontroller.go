package datacontroller

import (
	"net/http"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type DataController struct {
	staticAddr string

	cli HTTPClient
}

func New(
	staticAddr string,
	cli HTTPClient,
) *DataController {
	return &DataController{
		staticAddr: staticAddr,
		cli:        cli,
	}
}

func (dc *DataController) GetItems(w http.ResponseWriter, r *http.Request) {
	dc.proxyStaticRequestResponse(
		r.Context(),
		w,
		&proxyParams{
			destHTTPMethod: http.MethodGet,
			destPath:       ItemsPath,
			sourceReq:      r,
		},
	)
}

func (dc *DataController) GetPickupPoints(w http.ResponseWriter, r *http.Request) {
	dc.proxyStaticRequestResponse(
		r.Context(),
		w,
		&proxyParams{
			destHTTPMethod: http.MethodGet,
			destPath:       PickupPointsPath,
			sourceReq:      r,
		},
	)
}

func (dc *DataController) GetPayments(w http.ResponseWriter, r *http.Request) {
	dc.proxyStaticRequestResponse(
		r.Context(),
		w,
		&proxyParams{
			destHTTPMethod: http.MethodGet,
			destPath:       PaymentsPath,
			sourceReq:      r,
		},
	)
}
