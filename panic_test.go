package main

import (
	"testing"
	"net/http"
)

func TestPanicInServiceHandlerFunction(t *testing.T) {
	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request){
		w.Write([]byte("nice job!"))
		panic("fuck off!")
	})

	http.ListenAndServe(":8080", nil)
	t.Log("I am done!")
}

