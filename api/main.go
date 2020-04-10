package main

import (
	"fmt"
	"github.com/srgyrn/pact-example/api/model"
	"io/ioutil"
	"os"
)

type DBGateway interface {
	AddToDB() error
	Find(key string) (interface{}, error)
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
}

func refundOrder() {

}

func initDBs() error {
	dbs.vch = model.NewVoucherHandler()
	dbs.usr = model.NewUserHandler()
	dbs.ord = model.NewOrderDB()

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

	fmt.Println(dbs.usr, dbs.ord)

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