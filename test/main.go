package main

import (
	"fmt"
	"net/http"

	"github.com/Contra-Culture/gorr"
)

func main() {
	d, r := gorr.New(
		func(cfg *gorr.DispatcherCfgr) {

		})
	fmt.Print(r.String())
	fmt.Println()
	http.ListenAndServe(":8080", d)
}
