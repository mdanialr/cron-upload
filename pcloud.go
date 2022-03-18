package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
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
		fmt.Println("Path:", jsonCreateFolderResponse.Metadata.Path)
		fmt.Println("Name:", jsonCreateFolderResponse.Metadata.Name)
		fmt.Println("Id:", jsonCreateFolderResponse.Metadata.Id)
	}
}

func TryClearTrash(cl *http.Client, token string) {
	// CLEAR TRASH
	res, err := cl.Get(pcloud.GetClearTrashUrl(token))
	if err != nil {
		log.Fatalln("failed to when sending clear trash request to pCloud API:", err)
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalln("failed reading response body after sending request to clear trash:", err)
	}

	var jsonStdResponse pcloud.StdResponse
	if err = json.Unmarshal(b, &jsonStdResponse); err != nil {
		log.Fatalln("failed unmarshalling response body to json StdResponse model:", err)
	}

	if jsonStdResponse.Result != 0 {
		js, _ := service.PrettyJson(b)
		fmt.Println(js)
		log.Fatalln("response from StdResponse return non-0 value.")
	}
	fmt.Println("Successfully cleared")
}

func TryDeleteFile(cl *http.Client, token string) {
	// DELETE FILE from the given File id and one by one
	fileId := []string{"11649279231", "11649279554", "11649279603"}
	for _, id := range fileId {
		fmt.Println("Working on File id:", id)

		res, err := cl.Get(pcloud.GetDeleteFileUrl(token, id))
		if err != nil {
			log.Fatalln("failed to when sending delete file request to pCloud API:", err)
		}
		defer res.Body.Close()

		b, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatalln("failed reading response body after sending request to delete file:", err)
		}

		var jsonDeleteFileResponse pcloud.DeleteFileResponse
		if err = json.Unmarshal(b, &jsonDeleteFileResponse); err != nil {
			log.Fatalln("failed unmarshalling response body to json DeleteFileResponse model:", err)
		}

		if jsonDeleteFileResponse.Result != 0 {
			js, _ := service.PrettyJson(b)
			fmt.Println(js)
			log.Fatalln("response from DeleteFileResponse return non-0 value for file id:", id)
		}
		fmt.Println("Is Deleted:", jsonDeleteFileResponse.Meta.IsDeleted)
		fmt.Println("Filename:", jsonDeleteFileResponse.Meta.Name)
	}
}

func TryUploadFile(cl *http.Client, token string, fPath string) {
	const folderId = "2586338097"
	// PREPARE THE FILE FIRST
	fl, err := os.Open(fPath)
	if err != nil {
		log.Fatalln("failed to open file from filepath:", err)
	}
	defer fl.Close()

	// PREPARE MULTIPART FORM-DATA
	var buf = &bytes.Buffer{}
	wr := multipart.NewWriter(buf)
	part, err := wr.CreateFormFile("file", filepath.Base(fl.Name()))
	if err != nil {
		log.Fatalln("failed to create multi part form data from the given file:", err)
	}
	io.Copy(part, fl)
	wr.Close()

	// SEND POST REQUEST TO pCloud API
	req, _ := http.NewRequest(http.MethodPost, pcloud.GetUploadFileUrl(token, folderId), buf)
	req.Header.Add("content-type", wr.FormDataContentType())
	res, err := cl.Do(req)
	if err != nil {
		log.Fatalln("failed to when sending userinfo (quota) request to pCloud API:", err)
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalln("failed reading response body after sending request to userinfo (quota):", err)
	}

	var jsonUploadFileResponse pcloud.StdResponse
	if err = json.Unmarshal(b, &jsonUploadFileResponse); err != nil {
		log.Fatalln("failed unmarshalling response body to json UploadFileResponse model:", err)
	}

	if jsonUploadFileResponse.Result != 0 {
		js, _ := service.PrettyJson(b)
		fmt.Println(js)
		log.Fatalln("response from UploadFileResponse return non-0 value.")
	}
	fmt.Printf("File %s successfully uploaded", fl.Name())
}
