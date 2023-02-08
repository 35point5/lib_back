package main

import (
	"libback/ctrl"
	"libback/db"
)

func main() {

	database := db.NewDatabase()
	ctr := ctrl.MustNewCtrl(database)
	r := ctr.SetRouter()
	err := r.Run("127.0.0.1:8080")
	if err != nil {
		panic("Run Engine Failure!")
	}
	//obj := db.Book{Title: "aaa"}
	//tp := reflect.TypeOf(obj)
	//val := reflect.ValueOf(obj)
	//for i := 0; i < tp.NumField(); i++ {
	//	if true {
	//		fmt.Println(tp.Field(i).Name, val.Field(i).String())
	//	}
	//}
	//ctrl.MustNewCtrl(database)
}
