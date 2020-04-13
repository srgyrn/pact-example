package main

import (
	"bytes"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/srgyrn/pact-example/api/model"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
)

func Test_refundHandler(t *testing.T) {
	tests := []struct {
		name    string
		orderID int
		userKey string
		wantErr bool
	}{
		{
			name:    "returns bad request status when customer not found",
			orderID: 1,
			userKey: "barbara-streisand",
			wantErr: true,
		},
		{
			name:    "returns bad request status when order not found",
			orderID: 987,
			userKey: "john-doe",
			wantErr: true,
		},
		{
			name:    "returns successful status",
			orderID: 1,
			userKey: "john-doe",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initTestDBs()

			router := httprouter.New()
			router.POST("/order/:orderID/refund/", refundHandler)

			data := fmt.Sprintf("{\"user_key\": \"%s\"}", tt.userKey)

			req := httptest.NewRequest(
				http.MethodPost,
				fmt.Sprintf("/order/%d/refund/", tt.orderID),
				bytes.NewBufferString(data))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			if tt.wantErr && http.StatusBadRequest != rr.Code {

				t.Errorf("refundHandler(), want = %v, got = %v", http.StatusBadRequest, rr.Code)
			}

			if !tt.wantErr && http.StatusOK != rr.Code {
				r, _ := ioutil.ReadAll(rr.Result().Body)
				t.Errorf("%v \n %v", string(r), dbs.usr)
				t.Errorf("refundHandler(), want = %v, got = %v \n %v", http.StatusOK, rr.Code, string(r))
			}
		})
	}
}

func Test_makeRefund(t *testing.T) {
	type args struct {
		userKey string
		orderID string
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
				orderID: "1",
			},
			wantErr:       true,
			createVoucher: false,
		},
		{
			name: "returns error when order not found",
			args: args{
				userKey: "john-doe",
				orderID: "987",
			},
			wantErr:       true,
			createVoucher: false,
		},
		{
			name: "returns error when order is already refunded",
			args: args{
				userKey: "john-doe",
				orderID: "2",
			},
			wantErr:       true,
			createVoucher: false,
		},
		{
			name: "creates voucher and adds funds",
			args: args{
				userKey: "john-doe",
				orderID: "3",
			},
			wantErr:       false,
			createVoucher: true,
		},
		{
			name: "adds funds but does not create voucher account",
			args: args{
				userKey: "john-doe",
				orderID: "1",
			},
			wantErr:       false,
			createVoucher: false,
		},
	}
	for _, tt := range tests {
		initTestDBs()
		t.Run(tt.name, func(t *testing.T) {
			if err := makeRefund(tt.args.userKey, tt.args.orderID); (err != nil) != tt.wantErr {
				t.Errorf("makeRefund() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				var want model.Voucher
				var got model.Voucher

				if tt.createVoucher {
					dbs.ord.Find(tt.args.orderID)
					want, _ = model.NewVoucher(dbs.ord.Ord.Total, tt.args.userKey)
					dbs.vch.Find(model.GenerateKeyForVoucher(tt.args.userKey))
					got = *dbs.vch.Account
				}

				if !reflect.DeepEqual(want, got) {
					t.Errorf("makeRefund() failed. want: %v, got: %v", want, got)
				}
			}
		})
	}
}

func Test_makeRefund_existingVoucherAccount(t *testing.T) {
	initTestDBs()

	userKey := "jane-doe"
	va, _ := model.NewVoucher(100, userKey)
	dbs.vch.Account = &va
	dbs.vch.AddToDB()

	if err := makeRefund(userKey, strconv.Itoa(4)); err != nil {
		t.Errorf("makeRefund() error = %v", err)
	}

	type resultSet struct {
		voucher model.Voucher
		user    model.User
	}

	voucherExpected, _ := model.NewVoucher(250, userKey)
	userExpected := model.User{
		Name:     "Jane",
		LastName: "Doe",
		Balance:  150,
		Orders:   []int{4},
	}

	want := resultSet{
		voucherExpected,
		userExpected,
	}

	dbs.usr.Find(userKey)
	dbs.vch.Find(model.GenerateKeyForVoucher(userKey))

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
			Orders:   []int{4},
		},
	}

	dbs.usr = model.NewUserHandler()
	for _, u := range users {
		tmp := u
		dbs.usr.Usr = &tmp
		dbs.usr.AddToDB()
	}

	dbs.ord = model.NewOrderHandler()
	for _, o := range orders {
		tmp := o
		dbs.ord.Ord = &tmp
		dbs.ord.AddToDB()
	}

	dbs.vch = model.NewVoucherHandler()
}
