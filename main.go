package main

import (
	"gitlab.tiandy.com/lizewei08892/ginwebframework/bootstrap"
)

func main() {
	defer bootstrap.RealeaseDB()
	bootstrap.BootUp()
}
