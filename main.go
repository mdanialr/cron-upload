package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/mdanialr/cron-upload/internal/provider"
	"github.com/mdanialr/cron-upload/internal/provider/gdrive"
	awsS3 "github.com/mdanialr/cron-upload/internal/provider/s3"
	w "github.com/mdanialr/cron-upload/internal/worker"
	"github.com/mdanialr/cron-upload/internal/worker/v1"
	"github.com/mdanialr/cron-upload/pkg/config"
	h "github.com/mdanialr/cron-upload/pkg/helper"
	"github.com/mdanialr/cron-upload/pkg/logger"
	"github.com/mdanialr/cron-upload/pkg/scan"
	"github.com/spf13/viper"
)

var (
	configPath, logType string
	isTest              bool
)

func init() {
	flag.BoolVar(&isTest, "test", false, "test whether there is any error in the config file")
	flag.StringVar(&configPath, "path", ".", "locate the app config file. Default is set to current directory")
	flag.StringVar(&logType, "log", "stdout", "use '-log file' to write the logs to a file. Default is set to stdout")
	flag.Parse()
}

func main() {
	timer := time.Now()

	// init app config
	v, err := config.Init(configPath)
	if err != nil {
		log.Fatalln("failed to init config:", err)
	}
	// do some validation first to make sure all required fields are filled
	if err = config.Validate(v); err != nil {
		log.Fatalln("config file validation is failed:", err)
	}
	// sanitize config and setup necessary default value
	config.Sanitize(v)
	// make sure the provided provider name is currently supported
	if err = provider.ValidateSupportedClouds(v.GetString("provider.name")); err != nil {
		log.Fatalln(err)
	}
	// init logger and choose the log target output
	var lo logger.Writer
	switch logType {
	case "file":
		appLogPath := strings.TrimSuffix(v.GetString("log"), "/")
		appLogPath = fmt.Sprintf("%s/%s", appLogPath, "app")
		appLog, err := os.OpenFile(appLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0770)
		if err != nil {
			log.Fatalln("failed to init logger file for the app:", err)
		}
		defer appLog.Close()
		lo = logger.NewFile(appLog)
	default:
		lo = logger.NewFile(os.Stdout)
	}
	// init the log to do some prerequisite preparation
	lo.Init()
	// run only if not a test
	if !isTest {
		lo.WriteInf("START job")
		lo.WriteInf("")
	}
	// init wait group & cloud provider
	var wg sync.WaitGroup
	var chosenCloudProvider provider.Cloud
	// chose cloud provider based on value in the config file
	switch v.GetString("provider.name") {
	case string(provider.GoogleDrive):
		// init Google Drive cloud provider
		gDriveSvc, err := gdrive.Init(v.GetString("provider.cred"))
		if err != nil {
			lo.WriteErr("failed to init Google Drive:", err)
			return
		}
		chosenCloudProvider = gdrive.NewGoogleDriveProvider(gDriveSvc)
	case string(provider.S3Bucket):
		// additional validation for s3 bucket
		if err = config.ValidateS3Bucket(v); err != nil {
			log.Fatalln("config file validation for provider s3 bucket is failed:", err)
		}
		// init AWS S3 Bucket cloud provider
		s3Client, err := awsS3.Init(v.GetString("provider.cred"), v.GetString("provider.region"))
		if err != nil {
			lo.WriteErr("failed to init AWS S3 Bucket:", err)
			return
		}
		s3Ctx := context.Background()
		chosenCloudProvider = awsS3.NewS3BucketProvider(s3Ctx, v.GetString("provider.bucket"), s3Client)
	default:
		// because cloud provider is mandatory, throw error if there is no one provided
		lo.WriteErr("no cloud provider is provided")
		return
	}
	// run test if needed
	if isTest {
		testConfUploads := config.GetUploads(v)
		for _, testUpload := range testConfUploads {
			lo.WriteInf(h.LogStart("TEST Scan Local Dir", testUpload.Path))
			if _, err = scan.FilesAsc(testUpload.Path); err != nil {
				lo.WriteErr("failed to scan path:", err)
				return
			}
			lo.WriteInf(h.LogDone("TEST Scan Local Dir", testUpload.Path))
		}
		// test to create folder
		lo.WriteInf(h.LogStart("TEST Cloud Provider", ""))
		createdFolder, err := chosenCloudProvider.CreateFolder("test-cron-upload")
		if err != nil {
			lo.WriteErr("Failed to create a folder in the cloud provider:", err)
			return
		}
		if createdFolder == "" {
			lo.WriteErr("Failed to create a folder in the cloud provider:", err)
			return
		}
		lo.WriteInf(h.LogDone("TEST Cloud Provider", ""))
		// upload a dummy file then delete them
		lo.WriteInf(h.LogStart("TEST Upload", ""))
		dummyFile := io.NopCloser(strings.NewReader("hello world"))
		testPayload := provider.Payload{
			Name:   "test-cron-upload.txt",
			File:   dummyFile,
			Parent: []string{createdFolder},
		}
		createdTestPayload, err := chosenCloudProvider.UploadFile(&testPayload)
		if err != nil {
			lo.WriteErr("Failed to upload a test file:", err)
			return
		}
		lo.WriteInf(h.LogDone("TEST Upload", ""))
		// then delete them
		lo.WriteInf(h.LogStart("TEST Delete", ""))
		if err = chosenCloudProvider.Delete(createdTestPayload.Id); err != nil {
			lo.WriteErr("Failed to delete a test file:", err)
		}
		// also delete the dummy folder
		if err = chosenCloudProvider.Delete(createdFolder); err != nil {
			lo.WriteErr("Failed to delete a test folder:", err)
		}
		lo.WriteInf(h.LogDone("TEST Delete", ""))
		return
	}
	// init worker with the chosen cloud provider as dependency
	newWorker := worker.NewWorker(&wg, lo, chosenCloudProvider)
	// listing all available folder routes in the cloud provider
	cloudRoutes := newWorker.BuildRoutesProvider()
	// listing all routes from config file, then matched it with the routes from cloud provider
	for _, upload := range config.GetUploads(v) {
		trimmedRoute := strings.Trim(upload.Name, "/\\") // remove any the slash and back-slash
		route := fmt.Sprintf("%s/%s", v.GetString("root"), trimmedRoute)
		routeId := cloudRoutes[route]
		// create new routes in the cloud provider
		if routeId == "" {
			// use the root as starting point for the parent id
			routeId = cloudRoutes[v.GetString("root")]
			uploadRoutes := strings.Split(trimmedRoute, "/")
			// keep checking and creating the routes until reaching the last part of the route
			for _, ro := range uploadRoutes {
				currentParentId, err := newWorker.CheckCreateRoute(routeId, ro)
				if err != nil {
					routeId = "" // mark as empty so it does not trigger doTheJob
					lo.WriteErr("failed to check and create route", ro, "from", upload.Name, ":", err)
					break
				}
				routeId = currentParentId
				cloudRoutes[route] = currentParentId
			}
		}
		// do the job only if the route id is already known
		if routeId != "" {
			doTheJob(v, lo, &wg, newWorker, upload, routeId)
		}
	}
	wg.Wait() // block till all jobs are done

	lo.WriteInf("")
	lo.WriteInf("END job in:", time.Since(timer))
}

func doTheJob(
	v *viper.Viper,
	log logger.Writer,
	wg *sync.WaitGroup,
	worker w.I,
	upload config.UploadModel,
	routeId string,
) {
	// read local files from config
	files, err := scan.FilesAsc(upload.Path)
	if err != nil {
		log.WriteErr("failed to scan path:", err)
		return
	}
	// count retain from config but use retain from root as the default value
	retainMin := config.GetRetainExpiry(v, upload.Retain)
	// init necessary worker channels for this route
	listExpired := w.Channels{
		In:  make(chan *provider.Payload),
		Out: make(chan *provider.Payload),
	}
	listUnmatched := w.Channels{
		In:  listExpired.In,
		Out: make(chan *provider.Payload),
	}
	deleteFile := w.Channels{
		In: listExpired.Out,
	}
	uploadFile := w.Channels{
		In: listUnmatched.Out,
	}
	// spawn worker for delete & upload file since they can run independently
	for i := uint(1); i < v.GetUint("worker")+1; i++ {
		wg.Add(2)
		go worker.DeleteFile(deleteFile)
		go worker.UploadFile(uploadFile)
	}
	// one worker for each route is sufficient
	go worker.ListExpiredFiles(listExpired, routeId, retainMin)
	// because listening to In channel therefor we need to add wait group
	wg.Add(1)
	go worker.ListUnmatchedFiles(listUnmatched, routeId, files...)
}
