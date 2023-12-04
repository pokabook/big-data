package utils

import (
	"log"
	"net/http"
)

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func CheckCode(res *http.Response) {
	if res.StatusCode >= 400 {
		log.Fatalln(res.StatusCode)
	}
}
