package main

import (
	"blog/router"
)

func main() {

	r := router.Router()
	_ = r.Run(":8988")
}
