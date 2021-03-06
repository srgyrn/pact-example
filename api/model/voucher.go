package model

import (
	"errors"
	"strings"
)

// DefaultCurrency is the default currency of this project
const DefaultCurrency = "USD"

// Voucher holds data about the voucher account of user
type Voucher struct {
	Balance  float32 `json:"Balance"`
	Currency string  `json:"Currency"`
	userKey  string
}

// VoucherHandler holds the needed data for every DB operation to run
type VoucherHandler struct {
	Account *Voucher // Voucher account information
	db      map[string]*Voucher
}

// NewVoucher creates a Voucher object. If user key is not provided, it returns an error.
func NewVoucher(balance float32, userKey string) (Voucher, error) {
	if len(strings.TrimSpace(userKey)) == 0 {
		return Voucher{}, errors.New("user key cannot be empty")
	}

	return Voucher{
		Balance:  balance,
		Currency: DefaultCurrency,
		userKey:  userKey,
	}, nil
}

// NewVoucherHandler creates a VoucherHandler struct with empty initial values and returns it.
func NewVoucherHandler() *VoucherHandler {
	return &VoucherHandler{nil, make(map[string]*Voucher)}
}

// AddToDB adds given voucher account to DB
func (v *VoucherHandler) AddToDB() error {
	if v.Account.Currency != DefaultCurrency {
		return errors.New("wrong currency given")
	}

	key := GenerateKeyForVoucher(v.Account.userKey)
	v.db[key] = v.Account

	return nil
}

// Find looks for the given key in DB and returns it if it exists.
// If key is not provided or not found, function returns an error.
func (v *VoucherHandler) Find(key string) error {
	if len(strings.TrimSpace(key)) == 0 {
		return errors.New("key cannot be empty")
	}

	if account, ok := v.db[key]; ok {
		v.Account = account
		return nil
	}

	v.Account = nil
	return errors.New("account not found")
}

// Delete removes the given key from DB
func (v *VoucherHandler) Delete(key string) bool {
	if _, ok := v.db[key]; ok {
		delete(v.db, key)
		return true
	}

	return false
}

// UpdateBalance updates the balance of the given voucher account
func (v *VoucherHandler) UpdateBalance(amount float32) (float32, error) {
	key := GenerateKeyForVoucher(v.Account.userKey)

	if _, ok := v.db[key]; !ok {
		return 0, errors.New("account not found in DB")
	}

	va := v.db[key]
	va.Balance += amount
	return va.Balance, nil
}

// GenerateKeyForVoucher is a helper function that creates the key for the voucher account
func GenerateKeyForVoucher(userKey string) string {
	return userKey + "-" + strings.ToLower(DefaultCurrency)
}
