package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/linebot"
	"io/ioutil"
	"log"
	"m/models"
	"net/http"
	"net/url"
	"os"
)

func generateUrl(lat float64, lon float64) *url.URL{
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

func retrieveRestaurants(latitude float64, longitude float64) *models.GNavi{
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


func envLoad() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	envLoad()
	bot, err := linebot.New(
		os.Getenv("LINE_CHANNEL_SECRET"),
		os.Getenv("LINE_CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}

	//Setup HTTP Server
	http.HandleFunc("/callback", func(writer http.ResponseWriter, request *http.Request) {
		events, err := bot.ParseRequest(request)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				writer.WriteHeader(http.StatusBadRequest)
			} else {
				writer.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					_, err = bot.ReplyMessage(
						event.ReplyToken,
						linebot.NewTextMessage(
							message.Text),
						).Do()
					if err != nil {
						log.Fatal(err)
					}

				case *linebot.LocationMessage:
					lat, lon := message.Latitude, message.Longitude
					restaurants := retrieveRestaurants(lat, lon)
					if len(restaurants.Rest) > 0 {
						_, err = bot.ReplyMessage(
							event.ReplyToken,
							linebot.NewTextMessage(restaurants.Rest[0].Name),
						).Do()
						if err != nil {
							log.Fatal(err)
						}
					} else {
						_, err = bot.ReplyMessage(
							event.ReplyToken,
							linebot.NewLocationMessage(
								"近くにラーメン屋が見つからなかったよ😢",
								message.Address,
								message.Latitude,
								message.Longitude,
							),
						).Do()
						if err != nil {
							log.Fatal(err)
						}
					}
				}
			}
		}
	})

	if err := http.ListenAndServe(":" + os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}
}
