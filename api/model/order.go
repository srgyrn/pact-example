package model

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Payment type IDs
const (
	CreditCard = iota + 1
	CashOnDelivery
	Paypal
)

// Zone IDs
const (
	ZoneEurope = iota + 1
	ZoneMena
	ZoneAmerica
)

// Order holds every detail related to an order
type Order struct {
	ID                  int     `json:ID`
	Total               float32 `json:Total`
	PaymentWay          int     `json:PaymentWay`
	ShippingCountryZone int     `json:CountryZone`
	IsDeleted           bool    `json:IsDeleted`
}

// OrderHandler holds the needed data for every DB operation to run
type OrderHandler struct {
	Ord *Order
	db  map[string]*Order
}

// NewOrderDB creates and returns OrderHandler struct
func NewOrderDB() *OrderHandler {
	return &OrderHandler{
		Ord: nil,
		db:  make(map[string]*Order),
	}
}

func (o *OrderHandler) BulkInsert(b []byte) error {
	err := json.Unmarshal(b, &o.db)
	if err != nil {
		return fmt.Errorf("failed to unmarshal\n%s", err)
	}

	return nil
}

// AddToDB adds the order given to the DB.
// An error is thrown in the following circumstances:
//		- the payment way has not been set
//		- the country zone has not been set
func (o *OrderHandler) AddToDB() error {
	if o.Ord.PaymentWay == 0 {
		return errors.New("payment way is missing")
	}

	if o.Ord.ShippingCountryZone == 0 {
		return errors.New("zone is missing")
	}

	key := string(len(o.db) + 1)
	o.db[key] = o.Ord

	return nil
}

// Find function finds the order from db and returns it.
// An error is returned if key does not exist in DB map.
func (o *OrderHandler) Find(key string) (interface{}, error) {
	if _, ok := o.db[key]; !ok {
		return nil, errors.New("order not found")
	}

	return o.db[key], nil
}

// Delete function removes the order associated with the key in parameter from the DB.
func (o *OrderHandler) Delete(key string) bool {
	if _, ok := o.db[key]; ok {
		ord := o.db[key]
		ord.IsDeleted = true

		return true
	}

	return false
}
