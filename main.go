package main

import (
	"bytes"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mdanialr/cron-upload/internal/config"
	"github.com/mdanialr/cron-upload/internal/logger"
	"github.com/mdanialr/cron-upload/internal/provider/gdrive"
)

var (
	isDrive   bool          // whether to use Google Drive as provider or not
	isInit    bool          // initialize token in conjunction with Google Drive provider
	isRefresh bool          // exchange authorization code for new refresh token
	conf      *config.Model // global variable that would be used in this main pkg
)

func setupFlags() {
	flag.BoolVar(&isInit, "init", false, "retrieve token.json by using auth.json for Google Drive provider")
	flag.BoolVar(&isDrive, "drive", false, "use Google Drive as provider to upload files")
	flag.BoolVar(&isRefresh, "refresh", false, "exchange authorization code for new refresh token")
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

	// init internal logging
	if err := logger.InitLogger(conf); err != nil {
		log.Fatalln("failed to initialize internal logging:", err)
	}

	// if running using drive as provider
	if isDrive && isInit {
		logger.InfL.Println("START initialize token job")
		// if init params also included then init token first, before running the Google Drive's Job
		cl := &http.Client{}

		if err := gdrive.InitToken(conf, cl); err != nil {
			log.Fatalln("failed to initialize token for Google Drive:", err)
		}
	}

	// if refresh params also included then exchange authorization code for refresh token
	if isDrive && isRefresh {
		logger.InfL.Println("START renew refresh token job")
		if err := gdrive.Refresh(conf); err != nil {
			log.Fatalln("failed to exchange authorization code for refresh token:", err)
		}
	}

	// if running using drive as provider
	if isDrive && !isInit && !isRefresh {
		logger.InfL.Println("START job")
		// if init not included then run the job
		if err := gdrive.GoogleDrive(conf); err != nil {
			logger.ErrL.Println(err)
		}
	}

	logger.InfL.Println("END job in:", time.Since(timer))
}
