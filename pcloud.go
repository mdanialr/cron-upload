package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/mdanialr/go-cron-upload-to-cloud/internal/provider/pcloud"
	"github.com/mdanialr/go-cron-upload-to-cloud/internal/service"
)

func TryGetToken(cl *http.Client) {
	// GET DIGEST
	res, err := cl.Get(pcloud.GetDigestUrl())
	if err != nil {
		log.Fatalln("failed to when sending get digest request to pCloud API:", err)
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalln("failed reading response body after sending request to get digest:", err)
	}

	var jsonDigestResponse pcloud.StdResponse
	if err = json.Unmarshal(b, &jsonDigestResponse); err != nil {
		log.Fatalln("failed unmarshalling response body to json DigestResponse model:", err)
	}

	// GENERATE TOKEN via DIGEST AUTH
	if jsonDigestResponse.Result != 0 {
		js, _ := service.PrettyJson(b)
		fmt.Println(js)
		log.Fatalln("response from DigestResponse return non-0 value.")
	}
	us := pcloud.User{Username: "", Password: ""}
	res, err = cl.Get(us.GenerateTokenUrl(jsonDigestResponse.Digest))
	if err != nil {
		log.Fatalln("failed when sending request to generate token from pCloud API:", err)
	}
	defer res.Body.Close()

	b, err = io.ReadAll(res.Body)
	if err != nil {
		log.Fatalln("failed reading response body after sending request to generate token:", err)
	}

	var jsonTokenResponse pcloud.StdResponse
	if err = json.Unmarshal(b, &jsonTokenResponse); err != nil {
		log.Fatalln("failed unmarshalling response body to json TokenResponse model:", err)
	}

	if jsonTokenResponse.Result != 0 {
		js, _ := service.PrettyJson(b)
		fmt.Println(js)
		log.Fatalln("response from TokenResponse return non-0 value.")
	}
	fmt.Println("TOKEN:", jsonTokenResponse.Auth)
}

func TryLogout(cl *http.Client) {
	// LOGOUT / REMOVE / INVALIDATE TOKEN
	//res, err := cl.Get(pcloud.GetLogoutUrl(""))
	//if err != nil {
	//	log.Fatalln("failed to when sending logout request to pCloud API:", err)
	//}
	//defer res.Body.Close()
	//
	//b, err := io.ReadAll(res.Body)
	//if err != nil {
	//	log.Fatalln("failed reading response body after sending request to logout:", err)
	//}
	//
	//var jsonLogoutResponse pcloud.StdResponse
	//if err = json.Unmarshal(b, &jsonLogoutResponse); err != nil {
	//	log.Fatalln("failed unmarshalling response body to json LogoutResponse model:", err)
	//}
	//
	//if jsonLogoutResponse.Result != 0 {
	//	js, _ := service.PrettyJson(b)
	//	fmt.Println(js)
	//	log.Fatalln("response from LogoutResponse return non-0 value.")
	//}
	//fmt.Println("Deleted successfully:", jsonLogoutResponse.IsDeleted)
}

func TryPrintQuota(cl *http.Client, token string) {
	// PRINT STORAGE QUOTA
	res, err := cl.Get(pcloud.GetQuotaUrl(token))
	if err != nil {
		log.Fatalln("failed to when sending userinfo (quota) request to pCloud API:", err)
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalln("failed reading response body after sending request to userinfo (quota):", err)
	}

	var jsonQuotaResponse pcloud.StdResponse
	if err = json.Unmarshal(b, &jsonQuotaResponse); err != nil {
		log.Fatalln("failed unmarshalling response body to json QuotaResponse model:", err)
	}

	if jsonQuotaResponse.Result != 0 {
		js, _ := service.PrettyJson(b)
		fmt.Println(js)
		log.Fatalln("response from QuotaResponse return non-0 value.")
	}

	aQuota, err := service.BytesToAnyBit(jsonQuotaResponse.Quota, "Gb")
	if err != nil {
		log.Fatalln("failed to convert bytes to Gibibyte for available quota:", err)
	}
	uQuota, err := service.BytesToAnyBit(jsonQuotaResponse.UsedQuota, "Mb")
	if err != nil {
		log.Fatalln("failed to convert bytes to Gibibyte for used quota:", err)
	}
	fmt.Println("Available:", aQuota)
	fmt.Println("Used:", uQuota)
}

func TryCreateFolder(cl *http.Client, token string) {
	const SAMPLE = "vps/backup/db"
	// NOTES: should be created one by one. COULD NOT create folders recursively.
	// 1# /vps
	// 2# /vps/backup
	// 3# /vps/backup/db

	folders := strings.Split(SAMPLE, "/")
	var tmpFolders string
	for _, folder := range folders {
		tmpFolders += "/" + folder
		fmt.Println("Working on:", tmpFolders)
		// CREATE FOLDER one by one
		res, err := cl.Get(pcloud.GetCreateFolderUrl(token, tmpFolders))
		if err != nil {
			log.Fatalln("failed to when sending create folder request to pCloud API:", err)
		}
		defer res.Body.Close()

		b, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatalln("failed reading response body after sending request to create folder:", err)
		}

		var jsonCreateFolderResponse pcloud.StdResponse
		if err = json.Unmarshal(b, &jsonCreateFolderResponse); err != nil {
			log.Fatalln("failed unmarshalling response body to json CreateFolderResponse model:", err)
		}

		if jsonCreateFolderResponse.Result != 0 {
			js, _ := service.PrettyJson(b)
			fmt.Println(js)
			log.Fatalln("response from CreateFolderResponse return non-0 value.")
		}
		fmt.Println("Created:", jsonCreateFolderResponse.IsCreated)
	}
}
