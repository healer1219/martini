package main

import (
	"github.com/healer1219/martini/bootstrap"
)

func main() {
	defer bootstrap.RealeaseDB()
	bootstrap.Default().BootUp()
}
