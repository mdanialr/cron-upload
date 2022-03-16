package main

import (
	"encoding/json"
	"fmt"
	"github.com/mdanialr/go-cron-upload-to-cloud/service"
	"io"
	"log"
	"net/http"

	"github.com/mdanialr/go-cron-upload-to-cloud/pcloud"
)

func main() {
	cl := &http.Client{}

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

	var jsonDigestResponse pcloud.DigestResponse
	if err = json.Unmarshal(b, &jsonDigestResponse); err != nil {
		log.Fatalln("failed unmarshalling response body to json DigestResponse model:", err)
	}

	// GENERATE TOKEN via DIGEST AUTH
	if jsonDigestResponse.Result != 0 {
		js, _ := service.PrettyJson(b)
		fmt.Println(js)
		log.Fatalln("response from DigestResponse return non-0 value.")
	}
	us := pcloud.User{Username: "mdanialrma@gmail.com", Password: "vT83Cn2wcrdfHFIXndYkZ2rEdZYM03"}
	res, err = cl.Get(us.GenerateTokenUrl(jsonDigestResponse.Digest))
	if err != nil {
		log.Fatalln("failed when sending request to generate token from pCloud API:", err)
	}
	defer res.Body.Close()

	b, err = io.ReadAll(res.Body)
	if err != nil {
		log.Fatalln("failed reading response body after sending request to generate token:", err)
	}

	var jsonTokenResponse pcloud.TokenResponse
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
