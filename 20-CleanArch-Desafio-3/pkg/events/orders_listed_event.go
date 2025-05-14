package events

import "time"

type OrdersListedEvent struct {
	Name    string
	Payload interface{}
}

func NewOrdersListedEvent() *OrdersListedEvent {
	return &OrdersListedEvent{
		Name: "OrdersListed",
	}
}

func (e *OrdersListedEvent) GetName() string {
	return e.Name
}

func (e *OrdersListedEvent) GetPayload() interface{} {
	return e.Payload
}

func (e *OrdersListedEvent) SetPayload(payload interface{}) {
	e.Payload = payload
}

func (e *OrdersListedEvent) GetDateTime() time.Time {
	return time.Now()
}
