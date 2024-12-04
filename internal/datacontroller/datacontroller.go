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

	var cartItems []users.ItemCount
	if err = json.NewDecoder(resUsers.Body).Decode(&cartItems); err != nil {
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
	if err = json.NewDecoder(resStatic.Body).Decode(&items); err != nil {
		httpErr.InternalError(w, fmt.Errorf("decode static: %w", err))
		return
	}

	var res dto.Cart
	for _, cartItem := range cartItems {
		for _, item := range items {
			if cartItem.ItemID == item.ID {
				res.Items = append(res.Items, dto.ItemCard{
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

func (dc *DataController) UpdateCart(w http.ResponseWriter, r *http.Request) {
	dc.proxyUsersRequestResponse(
		r.Context(),
		w,
		&proxyParams{
			destHTTPMethod: http.MethodPatch,
			destPath:       CartPath,
			sourceReq:      r,
		},
	)
}

func (dc *DataController) CreateOrder(w http.ResponseWriter, r *http.Request) {
	dc.proxyUsersRequestResponse(
		r.Context(),
		w,
		&proxyParams{
			destHTTPMethod: http.MethodPost,
			destPath:       OrdersPath,
			sourceReq:      r,
		},
	)
}

func (dc *DataController) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	dc.proxyUsersRequestResponse(
		r.Context(),
		w,
		&proxyParams{
			destHTTPMethod: http.MethodPatch,
			destPath:       OrdersPath,
			sourceReq:      r,
		},
	)
}

func (dc *DataController) GetOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqUsers, err := http.NewRequestWithContext(ctx, http.MethodGet, dc.usersAddr+OrdersPath, http.NoBody)
	if err != nil {
		httpErr.InternalError(w, fmt.Errorf("http new request: %w", err))
		return
	}

	resUsers, err := dc.cli.Do(reqUsers)
	if err != nil {
		httpErr.InternalError(w, fmt.Errorf("cli do users: %w", err))
		return
	}

	var orders []users.Order
	if err = json.NewDecoder(resUsers.Body).Decode(&orders); err != nil {
		httpErr.InternalError(w, fmt.Errorf("decode users: %w", err))
		return
	}

	reqStaticItems, err := http.NewRequestWithContext(ctx, http.MethodGet, dc.staticAddr+ItemsPath, http.NoBody)
	if err != nil {
		httpErr.InternalError(w, fmt.Errorf("http new request: %w", err))
		return
	}

	resStaticItems, err := dc.cli.Do(reqStaticItems)
	if err != nil {
		httpErr.InternalError(w, fmt.Errorf("cli do static: %w", err))
		return
	}

	var items []static.Item
	if err = json.NewDecoder(resStaticItems.Body).Decode(&items); err != nil {
		httpErr.InternalError(w, fmt.Errorf("decode static: %w", err))
		return
	}

	reqStaticPayments, err := http.NewRequestWithContext(ctx, http.MethodGet, dc.staticAddr+PaymentsPath, http.NoBody)
	if err != nil {
		httpErr.InternalError(w, fmt.Errorf("http new request: %w", err))
		return
	}

	resStaticPayments, err := dc.cli.Do(reqStaticPayments)
	if err != nil {
		httpErr.InternalError(w, fmt.Errorf("cli do static: %w", err))
		return
	}

	var payments []static.Payment
	if err = json.NewDecoder(resStaticPayments.Body).Decode(&payments); err != nil {
		httpErr.InternalError(w, fmt.Errorf("decode static: %w", err))
		return
	}

	reqStaticPickupPoints, err := http.NewRequestWithContext(ctx, http.MethodGet, dc.staticAddr+PickupPointsPath, http.NoBody)
	if err != nil {
		httpErr.InternalError(w, fmt.Errorf("http new request: %w", err))
		return
	}

	resStaticPickupPoints, err := dc.cli.Do(reqStaticPickupPoints)
	if err != nil {
		httpErr.InternalError(w, fmt.Errorf("cli do static: %w", err))
		return
	}

	var pickupPoints []static.PickupPoint
	if err = json.NewDecoder(resStaticPickupPoints.Body).Decode(&pickupPoints); err != nil {
		httpErr.InternalError(w, fmt.Errorf("decode static: %w", err))
		return
	}

	res := make([]dto.Order, len(orders))
	for i, order := range orders {
		res[i].ID = order.ID
		res[i].Status = order.Status

		// items + price total
		for _, orderItem := range order.Items {
			for _, item := range items {
				if orderItem.ItemID == item.ID {
					res[i].Items = append(res[i].Items, dto.ItemCard{
						Item: dto.Item{
							ID:    item.ID,
							Name:  item.Name,
							Price: item.Price,
						},
						Count: orderItem.Count,
					})
					res[i].PriceTotal += item.Price * orderItem.Count
					break
				}
			}
		}

		// payment
		for _, payment := range payments {
			if payment.ID == order.PaymentID {
				res[i].Payment = dto.Payment{
					ID:   order.PaymentID,
					Name: payment.Name,
				}
			}
		}

		// postomat
		for _, pickupPoint := range pickupPoints {
			if pickupPoint.ID == order.PostomatID {
				res[i].Postomat = dto.PickupPoint{
					ID:      order.PostomatID,
					Address: pickupPoint.Address,
				}
			}
		}
	}

	writer.WriteJson(w, res)
}
