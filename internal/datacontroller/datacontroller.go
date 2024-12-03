package datacontroller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"oreshnik/pkg/dto"
	"oreshnik/pkg/dto/static"
	"oreshnik/pkg/dto/users"
	httpErr "oreshnik/pkg/http/error"
	"oreshnik/pkg/http/writer"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type DataController struct {
	staticAddr string
	usersAddr  string

	cli HTTPClient
}

func New(
	staticAddr string,
	usersAddr string,
	cli HTTPClient,
) *DataController {
	return &DataController{
		staticAddr: staticAddr,
		usersAddr:  usersAddr,
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

func (dc *DataController) GetCart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqUsers, err := http.NewRequestWithContext(ctx, http.MethodGet, dc.usersAddr+CartPath, http.NoBody)
	if err != nil {
		httpErr.InternalError(w, fmt.Errorf("http new request: %w", err))
		return
	}

	resUsers, err := dc.cli.Do(reqUsers)
	if err != nil {
		httpErr.InternalError(w, fmt.Errorf("cli do users: %w", err))
		return
	}

	var cartItems []users.CartItem
	if err := json.NewDecoder(resUsers.Body).Decode(&cartItems); err != nil {
		httpErr.InternalError(w, fmt.Errorf("decode users: %w", err))
		return
	}

	reqStatic, err := http.NewRequestWithContext(ctx, http.MethodGet, dc.staticAddr+ItemsPath, http.NoBody)
	if err != nil {
		httpErr.InternalError(w, fmt.Errorf("http new request: %w", err))
		return
	}

	resStatic, err := dc.cli.Do(reqStatic)
	if err != nil {
		httpErr.InternalError(w, fmt.Errorf("cli do static: %w", err))
		return
	}

	var items []static.Item
	if err := json.NewDecoder(resStatic.Body).Decode(&items); err != nil {
		httpErr.InternalError(w, fmt.Errorf("decode static: %w", err))
		return
	}

	var res dto.Cart
	for _, cartItem := range cartItems {
		for _, item := range items {
			if cartItem.ItemID == item.ID {
				res.Items = append(res.Items, dto.ItemCart{
					Item: dto.Item{
						ID:    item.ID,
						Name:  item.Name,
						Price: item.Price,
					},
					Count: cartItem.Count,
				})
				res.PriceTotal += item.Price * cartItem.Count
				break
			}
		}
	}

	writer.WriteJson(w, res)
}
