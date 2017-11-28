// Package gpserver handles the logic of the gPanel server
package gpserver

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/Ennovar/gPanel/pkg/api/bundle"
	"github.com/Ennovar/gPanel/pkg/gpaccount"
)

type Controller struct {
	Directory    string
	DocumentRoot string
	Bundles      map[string]*gpaccount.Controller
	ServerLogger *log.Logger
	APILogger    *log.Logger
}

func New() *Controller {
	bundles := make(map[string]*gpaccount.Controller)

	dirs, err := ioutil.ReadDir("bundles/")
	if err != nil {
		fmt.Println("error finding bundles:", err.Error())
	}

	f, err := os.OpenFile("server/logs/server_errors.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("error whilst trying to start server logging instance:", err.Error())
	}

	apiLogger := log.New(f, "API :: ", 3)

	for _, dir := range dirs {
		if dir.Name() == "default_bundle" || !dir.IsDir() {
			continue
		}

		if strings.HasPrefix(dir.Name(), "bundle_") {
			dirPath := "bundles/" + dir.Name() + "/"
			err, accPort, pubPort := bundle.GetPorts(dirPath)

			curBundle := gpaccount.New(dirPath, accPort, pubPort)

			err = curBundle.Start()
			err2 := curBundle.Public.Start()
			if err != nil || err2 != nil {
				fmt.Println("error starting bundle:", dir.Name())
			}

			bundles[strings.Replace(dir.Name(), "bundle_", "", 1)] = curBundle
		}
	}

	return &Controller{
		Directory:    "server/",
		DocumentRoot: "document_root/",
		Bundles:      bundles,
		ServerLogger: log.New(f, "SERVER :: ", 3),
		APILogger:    apiLogger,
	}
}
