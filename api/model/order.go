package model

import (
	"errors"
)

const (
	CreditCard = iota + 1
	CashOnDelivery
	Paypal
)

const (
	ZoneEurope = iota + 1
	ZoneMena
	ZoneAmerica
)

type Order struct {
	ID                  int
	Total               float32
	PaymentWay          int
	ShippingCountryZone int
	IsDeleted           bool
}

type OrderDB struct {
	Ord Order
	db  map[string]*Order
}

func (o *OrderDB) AddToDB() error {
	if o.Ord.PaymentWay == 0 {
		return errors.New("payment way is missing")
	}

	if o.Ord.ShippingCountryZone == 0 {
		return errors.New("zone is missing")
	}

	key := string(len(o.db) + 1)
	o.db[key] = &o.Ord

	return nil
}

func (o *OrderDB) Find(key string) (interface{}, error) {
	if _, ok := o.db[key]; !ok {
		return nil, errors.New("order not found")
	}

	return *o.db[key], nil
}

func (o *OrderDB) Delete(key string) bool {
	if _, ok := o.db[key]; ok {
		ord := o.db[key]
		ord.IsDeleted = true

		return true
	}

	return false
}
