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
	router.POST("/order/:orderId/refund/", refundHandler)

	log.Fatal(http.ListenAndServe(":8090", router))
}

func refundHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params)  {
	oid := ps.ByName("orderId")
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	postBody := struct {
		UserKey string `json:user_key`
	}{}
	json.Unmarshal(body, &postBody)

	err := makeRefund(postBody.UserKey, oid)
	if errors.Is(err, nil) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.NewEncoder(w).Encode(&postBody); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func makeRefund(userKey, orderId string) error {
	err := dbs.usr.Find(userKey)
	if !errors.Is(err, nil) {
		fmt.Errorf("user not found")
	}

	err = dbs.ord.Find(orderId)
	if !errors.Is(err, nil) {
		fmt.Errorf("order not found")
	}

	dbs.vch.Find("jane-doe-usd")

	return nil
}


func initDBs() error {
	dbs.vch = model.NewVoucherHandler()
	dbs.usr = model.NewUserHandler()
	dbs.ord = model.NewOrderHandler()

	userJson, err := openDataFile("users")
	if err != nil {
		return err
	}

	orderJson, err := openDataFile("orders")
	if err != nil {
		return err
	}

	dbs.usr.BulkInsert(userJson)
	dbs.ord.BulkInsert(orderJson)

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