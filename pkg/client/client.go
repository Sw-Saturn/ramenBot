package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"m/pkg/models"
	"net/http"
	"net/url"
	"os"
)

func generateUrl(lat float64, lon float64) *url.URL {
	u := &url.URL{}
	u.Scheme = "https"
	u.Host = "api.gnavi.co.jp"
	u.Path = "RestSearchAPI/v3/"
	q := u.Query()
	q.Set("keyid", os.Getenv("GNAVI_ACCESS_TOKEN"))
	q.Set("latitude", fmt.Sprintf("%f", lat))
	q.Set("longitude", fmt.Sprintf("%f", lon))
	q.Set("category_s", "RSFST08008")
	q.Set("category_s", "RSFST08008")
	q.Set("range", "4")
	u.RawQuery = q.Encode()
	return u
}

func RetrieveRestaurants(latitude float64, longitude float64) *models.GNavi {
	reqUrl := generateUrl(latitude, longitude)
	req, err := http.NewRequest("GET", reqUrl.String(), nil)
	if err != nil {
		log.Fatal(err)
	}
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(byteArray))

	data := new(models.GNavi)

	if err := json.Unmarshal(byteArray, data); err != nil {
		fmt.Println("JSON Unmarshal error:", err)
		return nil
	}
	return data
}
