package main

import (
	"fmt"
	"net/http"

	"github.com/kitabisa/gorm"
	_ "github.com/kitabisa/gorm/dialects/sqlite"
	"github.com/kitabisa/newrelic-context"
	"github.com/kitabisa/newrelic-context/nrgorm"
)

var db *gorm.DB

func initDB() *gorm.DB {
	db, err := gorm.Open("sqlite3", "./foo.db")
	if err != nil {
		panic(err)
	}
	nrgorm.AddGormCallbacks(db)
	return db
}

type Product struct {
	ID   int
	Name string
}

func catalogPage(db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		var products []Product
		db := nrcontext.SetTxnToGorm(req.Context(), db)
		db.Find(&products)
		for i, v := range products {
			rw.Write([]byte(fmt.Sprintf("%v. %v\n", i, v.Name)))
		}
	})
}

func other_main() {
	db = initDB()
	defer db.Close()

	handler := catalogPage(db)
	nrmiddleware, _ := nrcontext.NewMiddleware("test-app", "my-license-key")
	handler = nrmiddleware.Handler(handler)

	http.ListenAndServe(":8000", handler)
}
