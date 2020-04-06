package main

type DBOperations interface {
	AddToDB() error
	Find(key string) (interface{}, error)
	Delete(key string) bool
}


func main() {
}
