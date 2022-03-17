package main

import "net/http"

func TryGetToken(cl *http.Client) {
	//// GET DIGEST
	//res, err := cl.Get(pcloud.GetDigestUrl())
	//if err != nil {
	//	log.Fatalln("failed to when sending get digest request to pCloud API:", err)
	//}
	//defer res.Body.Close()
	//
	//b, err := io.ReadAll(res.Body)
	//if err != nil {
	//	log.Fatalln("failed reading response body after sending request to get digest:", err)
	//}
	//
	//var jsonDigestResponse pcloud.StdResponse
	//if err = json.Unmarshal(b, &jsonDigestResponse); err != nil {
	//	log.Fatalln("failed unmarshalling response body to json DigestResponse model:", err)
	//}
	//
	//// GENERATE TOKEN via DIGEST AUTH
	//if jsonDigestResponse.Result != 0 {
	//	js, _ := service.PrettyJson(b)
	//	fmt.Println(js)
	//	log.Fatalln("response from DigestResponse return non-0 value.")
	//}
	//us := pcloud.User{Username: "", Password: ""}
	//res, err = cl.Get(us.GenerateTokenUrl(jsonDigestResponse.Digest))
	//if err != nil {
	//	log.Fatalln("failed when sending request to generate token from pCloud API:", err)
	//}
	//defer res.Body.Close()
	//
	//b, err = io.ReadAll(res.Body)
	//if err != nil {
	//	log.Fatalln("failed reading response body after sending request to generate token:", err)
	//}
	//
	//var jsonTokenResponse pcloud.StdResponse
	//if err = json.Unmarshal(b, &jsonTokenResponse); err != nil {
	//	log.Fatalln("failed unmarshalling response body to json TokenResponse model:", err)
	//}
	//
	//if jsonTokenResponse.Result != 0 {
	//	js, _ := service.PrettyJson(b)
	//	fmt.Println(js)
	//	log.Fatalln("response from TokenResponse return non-0 value.")
	//}
	//fmt.Println("TOKEN:", jsonTokenResponse.Auth)
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
