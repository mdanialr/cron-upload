package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mdanialr/cron-upload/internal/config"
	"github.com/mdanialr/cron-upload/internal/provider/gdrive"
)

var (
	isDrive bool          // whether to use Google Drive as provider or not
	isInit  bool          // initialize token in conjunction with Google Drive provider
	conf    *config.Model // global variable that would be used in this main pkg
)

func setupFlags() {
	flag.BoolVar(&isInit, "init", false, "retrieve token.json by using auth.json for Google Drive provider")
	flag.BoolVar(&isDrive, "drive", false, "use Google Drive as provider to upload files")
	flag.Parse()
}

func main() {
	timer := time.Now()
	setupFlags()

	f, err := os.ReadFile("app-config.yml")
	if err != nil {
		log.Fatalln("failed to read config file:", err)
	}

	conf, err = config.NewConfig(bytes.NewReader(f))
	if err != nil {
		log.Fatalln("failed to create new config instance:", err)
	}
	if err := conf.Sanitization(); err != nil {
		log.Fatalln("failed to sanitize config file, please make sure config file has valid values:", err)
	}
	if err := conf.Upload.Sanitization(); err != nil {
		log.Fatalln("failed to sanitize upload in config file, please make sure upload section has valid values:", err)
	}

	// if running using drive as provider
	if isDrive && isInit {
		// if init params also included then init token first, before running the Google Drive's Job
		cl := &http.Client{}

		if err := gdrive.InitToken(conf, cl); err != nil {
			log.Fatalln("failed to initialize token for Google Drive:", err)
		}
	}

	// if running using drive as provider
	if isDrive && !isInit {
		// if init not included then run the job
		gdrive.GoogleDrive(conf)
	}

	fmt.Println("\nElapsed time:", time.Since(timer))
}
