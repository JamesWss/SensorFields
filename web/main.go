package main

import (
	"net/http"

	"github.com/ericxiao417/sensor-farm/web/controller"
)

func main() {
	controller.Initialize()

	http.ListenAndServe(":3000", nil)
}
