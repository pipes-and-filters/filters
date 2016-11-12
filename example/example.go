package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pipes-and-filters/filters"
)

func main() {
	chain, err := filters.ChainFile("chain.yml")
	if err != nil {
		log.Fatal(err)
	}
	e := Example{chain}
	err = http.ListenAndServe(fmt.Sprintf(":8080"), e)
	if err != nil {
		log.Fatal(err)
	}
}

type Example struct {
	Chain filters.Chain
}

func (e Example) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	exec, err := e.Chain.Exec()
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		log.Print(err)
	}
	exec.SetInput(r.Body)
	exec.SetOutput(w)
	err = exec.Run()
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		log.Print(err)
	}
}
