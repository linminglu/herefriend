package main

import (
	"fmt"
	"os"

	"github.com/go-martini/martini"

	"herefriend/lib"
	"herefriend/server/routes"
)

const gPidFile = "/var/run/herefriend.pid"

func main() {
	f, err := os.Create(gPidFile)
	if nil == err {
		f.WriteString(fmt.Sprintf("%d", os.Getpid()))
		f.Close()
	} else {
		fmt.Println(err)
	}

	defer lib.CloseSQL()

	m := martini.Classic()
	routes.InstallRoutes(m)
	m.RunOnAddr(":8080")
}
