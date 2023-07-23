package main

import (
	"github.com/healer1219/gin-web-framework/bootstrap"
)

func main() {
	defer bootstrap.RealeaseDB()
	bootstrap.BootUp()
}
