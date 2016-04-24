package main

import (
	"fmt"
	"os"

	log "github.com/cihub/seelog"
	"github.com/go-martini/martini"

	"herefriend/lib"
	"herefriend/server/routes"
)

const gPidFile = "./herefriend.pid"

func main() {
	log.Info("Start")

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
