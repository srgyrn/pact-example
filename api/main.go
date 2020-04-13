package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/srgyrn/pact-example/api/model"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type DBGateway interface {
	AddToDB() error
	Find(key string) error
	Delete(key string) bool
}

type BalanceHolder interface {
	UpdateBalance(balance float32) (float32, error)
}

var dbs struct {
	vch *model.VoucherHandler
	usr *model.UserHandler
	ord *model.OrderHandler
}

func main() {
	err := initDBs()

	if err != nil {
		fmt.Println(err)
	}

	router := httprouter.New()
	router.POST("/order/:orderID/refund/", refundHandler)

	log.Fatal(http.ListenAndServe(":8090", router))
}

func refundHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params)  {
	oid := ps.ByName("orderID")
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	postBody := struct {
		UserKey string `json:"user_key"`
	}{}
	json.Unmarshal(body, &postBody)

	if len(strings.TrimSpace(postBody.UserKey)) == 0 {
		//http.Error(w, string(body), http.StatusBadRequest)
		http.Error(w, fmt.Sprintf("%v \n %v", string(body), postBody.UserKey), http.StatusBadRequest)
		return
	}

	err := makeRefund(postBody.UserKey, oid)
	if !errors.Is(err, nil) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.NewEncoder(w).Encode(&postBody); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func makeRefund(userKey, orderID string) error {
	err := dbs.usr.Find(userKey)
	if !errors.Is(err, nil) {
		return fmt.Errorf("user not found: %s", userKey)
	}

	err = dbs.ord.Find(orderID)
	if !errors.Is(err, nil) {
		return fmt.Errorf("order not found")
	}

	order := dbs.ord.Ord

	if order.IsDeleted {
		return fmt.Errorf("order already refunded")
	}

	order.IsDeleted = true

	refundToVoucher := false
	if order.ShippingCountryZone == model.ZoneMena && order.PaymentWay == model.CashOnDelivery {
		refundToVoucher = true
	}

	user := dbs.usr.Usr

	if !refundToVoucher {
		user.UpdateBalance(order.Total)
		return nil
	}

	err = dbs.vch.Find(model.GenerateKeyForVoucher(userKey))

	if err != nil && dbs.vch.Account == nil {
		va, err := model.NewVoucher(order.Total, model.GenerateKeyForUser(user))

		if !errors.Is(err, nil) {
			return err
		}

		dbs.vch.Account = &va
		err = dbs.vch.AddToDB()

		return nil
	}

	dbs.vch.UpdateBalance(order.Total)
	return nil
}


func initDBs() error {
	dbs.vch = model.NewVoucherHandler()
	dbs.usr = model.NewUserHandler()
	dbs.ord = model.NewOrderHandler()

	userJSON, err := openDataFile("users")
	if err != nil {
		return err
	}

	orderJSON, err := openDataFile("orders")
	if err != nil {
		return err
	}

	dbs.usr.BulkInsert(userJSON)
	dbs.ord.BulkInsert(orderJSON)

	return nil
}

func openDataFile(fileName string) ([]byte, error) {
	dataFile, err := os.Open("./api/data/" + fileName + ".json")
	defer dataFile.Close()

	if err != nil {
		return nil, fmt.Errorf("cannot open %s.json", fileName)
	}

	byteValue, err := ioutil.ReadAll(dataFile)
	if err != nil {
		fmt.Printf("failed to byte\n%s", err)
		os.Exit(1)
	}

	return byteValue, nil
}