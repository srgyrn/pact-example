package main

import (
	"bytes"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/srgyrn/pact-example/api/model"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func Test_refundHandler(t *testing.T) {
	tests := []struct {
		name    string
		orderId int
		userKey string
		wantErr bool
	}{
		{
			name:    "returns error when customer not found",
			orderId: 1,
			userKey: "barbara-streisand",
			wantErr: true,
		},
		{
			name:    "returns error when order not found",
			orderId: 987,
			userKey: "john-doe",
			wantErr: true,
		},
		{
			name:    "creates voucher and adds funds",
			orderId: 1,
			userKey: "john-doe",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := httprouter.New()
			router.GET("/order/:orderId/refund/", refundHandler)

			req := httptest.NewRequest(
				http.MethodPost,
				fmt.Sprintf("/order/%d/refund/", tt.orderId),
				bytes.NewBufferString(fmt.Sprintf("{\"user_key\": \"%s\"}", tt.userKey)))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			if tt.wantErr && http.StatusBadRequest != rr.Code {
				t.Errorf("")
			}

			if !tt.wantErr && http.StatusOK != rr.Code {
				t.Errorf("")
			}
		})
	}
}

func Test_makeRefund(t *testing.T) {
	type args struct {
		userKey string
		orderId string
	}
	tests := []struct {
		name          string
		args          args
		wantErr       bool
		createVoucher bool
	}{
		{
			name: "returns error when user not found",
			args: args{
				userKey: "barbara-streisand",
				orderId: "1",
			},
			wantErr:       true,
			createVoucher: false,
		},
		{
			name: "returns error when order not found",
			args: args{
				userKey: "john-doe",
				orderId: "987",
			},
			wantErr:       true,
			createVoucher: false,
		},
		{
			name: "returns error when order is already refunded",
			args: args{
				userKey: "john-doe",
				orderId: "2",
			},
			wantErr:       true,
			createVoucher: false,
		},
		{
			name: "creates voucher and adds funds",
			args: args{
				userKey: "john-doe",
				orderId: "3",
			},
			wantErr:       false,
			createVoucher: true,
		},
		{
			name: "adds funds but does not create voucher account",
			args: args{
				userKey: "john-doe",
				orderId: "1",
			},
			wantErr:       false,
			createVoucher: false,
		},
	}
	for _, tt := range tests {
		initTestDBs()
		t.Run(tt.name, func(t *testing.T) {
			if err := makeRefund(tt.args.userKey, tt.args.orderId); (err != nil) != tt.wantErr {
				t.Errorf("makeRefund() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				want := model.Voucher{}
				if tt.createVoucher {
					dbs.ord.Find(string(tt.args.orderId))
					want, _ = model.NewVoucher(dbs.ord.Ord.Total, tt.args.userKey)
				}

				dbs.vch.Find(tt.args.userKey + "-usd")
				got := dbs.vch.Account
				if !reflect.DeepEqual(want, got) {
					t.Errorf("makeRefund() failed. want: %v, got: %v", want, got)
				}
			}
		})
	}
}

func Test_makeRefund_existingVoucherAccount(t *testing.T) {
	initTestDBs()
	if err := makeRefund("jane-doe", string(4)); err != nil {
		t.Errorf("makeRefund() error = %v", err)
	}

	type resultSet struct {
		voucher model.Voucher
		user    model.User
	}

	voucherExpected, _ := model.NewVoucher(150, "jane-doe")
	userExpected := model.User{
		Name:     "Jane",
		LastName: "Doe",
		Balance:  150,
		Orders:   []int{},
	}

	want := resultSet{
		voucherExpected,
		userExpected,
	}

	dbs.usr.Find("jane-doe")
	dbs.vch.Find("jane-doe-usd")

	got := resultSet{
		*dbs.vch.Account,
		*dbs.usr.Usr,
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("makeRefund() failed. want = %v, got = %v", want, got)
	}
}

func initTestDBs() {
	orders := []model.Order{
		{
			ID:                  1,
			Total:               100,
			PaymentWay:          model.CreditCard,
			ShippingCountryZone: model.ZoneEurope,
			IsDeleted:           false,
		},
		{
			ID:                  2,
			Total:               200,
			PaymentWay:          model.CashOnDelivery,
			ShippingCountryZone: model.ZoneMena,
			IsDeleted:           true,
		},
		{
			ID:                  3,
			Total:               300,
			PaymentWay:          model.CashOnDelivery,
			ShippingCountryZone: model.ZoneMena,
			IsDeleted:           false,
		},
		{
			ID:                  4,
			Total:               150,
			PaymentWay:          model.CashOnDelivery,
			ShippingCountryZone: model.ZoneMena,
			IsDeleted:           false,
		},
	}

	users := []model.User{
		{
			Name:     "John",
			LastName: "Doe",
			Balance:  100,
			Orders:   []int{1, 2, 3},
		},
		{
			Name:     "Jane",
			LastName: "Doe",
			Balance:  150,
			Orders:   []int{},
		},
	}

	dbs.usr = model.NewUserHandler()
	for _, u := range users {
		dbs.usr.Usr = &u
		dbs.usr.AddToDB()
	}

	dbs.ord = model.NewOrderHandler()
	for _, o := range orders {
		dbs.ord.Ord = &o
		dbs.ord.AddToDB()
	}

	dbs.vch = model.NewVoucherHandler()
}
