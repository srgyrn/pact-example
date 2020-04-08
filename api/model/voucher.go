package model

import (
	"errors"
	"strings"
)

// DefaultCurrency is the default currency of this project
const DefaultCurrency = "USD"

// Voucher holds data about the voucher account of user
type Voucher struct {
	Balance  float32
	Currency string
	userKey  string
}

// VoucherDB holds the needed data for every DB operation to run
type VoucherDB struct {
	Account Voucher // Voucher account information
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

// NewVoucherDB creates a VoucherDB struct with empty initial values and returns it.
func NewVoucherDB() *VoucherDB {
	return &VoucherDB{Voucher{}, make(map[string]*Voucher)}
}

// AddToDB adds given voucher account to DB
func (v *VoucherDB) AddToDB() error {
	if v.Account.Currency != DefaultCurrency {
		return errors.New("wrong currency given")
	}

	key := generateKeyForVoucher(v.Account.userKey)
	v.db[key] = &v.Account

	return nil
}

// Find looks for the given key in DB and returns it if it exists.
// If key is not provided or not found, function returns an error.
func (v *VoucherDB) Find(key string) (interface{}, error) {
	if len(strings.TrimSpace(key)) == 0 {
		return nil, errors.New("key cannot be empty")
	}

	if account, ok := v.db[key]; ok {
		return account, nil
	}

	return nil, errors.New("account not found")
}

// Delete removes the given key from DB
func (v *VoucherDB) Delete(key string) bool {
	if _, ok := v.db[key]; ok {
		delete(v.db, key)
		return true
	}

	return false
}

// UpdateBalance updates the balance of the given voucher account
func (v *VoucherDB) UpdateBalance(amount float32) (float32, error) {
	key := generateKeyForVoucher(v.Account.userKey)

	if _, ok := v.db[key]; !ok {
		return 0, errors.New("account not found in DB")
	}

	va := v.db[key]
	va.Balance += amount
	return va.Balance, nil
}

func generateKeyForVoucher(userKey string) string {
	return userKey + "-" + strings.ToLower(DefaultCurrency)
}
