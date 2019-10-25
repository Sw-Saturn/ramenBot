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
						linebot.NewTextMessage(message.Text),
					).Do()
					if err != nil {
						log.Fatal(err)
					}

				case *linebot.LocationMessage:
					lat, lon := message.Latitude, message.Longitude
					restaurants := retrieveRestaurants(lat, lon)

					if len(restaurants.Rest) > 0 {
						var flexColumns []*linebot.BubbleContainer
						for i := range restaurants.Rest{
							var flexContent = &linebot.BubbleContainer{
								Type:   linebot.FlexContainerTypeBubble,
								Body:	&linebot.BoxComponent{
									Type: 		linebot.FlexComponentTypeBox,
									Layout:		linebot.FlexBoxLayoutTypeVertical,
									Contents:	[]linebot.FlexComponent{
										&linebot.TextComponent{
											Type:     linebot.FlexComponentTypeText,
											Text:     restaurants.Rest[i].Name,
											Weight:   linebot.FlexTextWeightTypeBold,
											Size:     linebot.FlexTextSizeTypeXl,
										},
										&linebot.BoxComponent{
											Type:     linebot.FlexComponentTypeBox,
											Layout:   linebot.FlexBoxLayoutTypeVertical,
											Contents: []linebot.FlexComponent{
												&linebot.BoxComponent{
													Type:     linebot.FlexComponentTypeBox,
													Layout:   linebot.FlexBoxLayoutTypeBaseline,
													Contents: []linebot.FlexComponent{
														&linebot.TextComponent{
															Type:	linebot.FlexComponentTypeText,
															Text:	"Access",
															Size:	linebot.FlexTextSizeTypeSm,
															Flex:	linebot.IntPtr(1),
															Color:	"#aaaaaa",
														},
														&linebot.TextComponent{
															Type:	linebot.FlexComponentTypeText,
															Text:	fmt.Sprintf("%s %s %s 分", restaurants.Rest[i].Access.Line, restaurants.Rest[i].Access.Station, restaurants.Rest[i].Access.Walk),
															Wrap:	true,
															Size:	linebot.FlexTextSizeTypeSm,
															Flex:	linebot.IntPtr(5),
															Color:	"#666666",
														},
													},
													Spacing:  linebot.FlexComponentSpacingTypeSm,
												},
												&linebot.BoxComponent{
													Type:     linebot.FlexComponentTypeBox,
													Layout:   linebot.FlexBoxLayoutTypeBaseline,
													Contents: []linebot.FlexComponent{
														&linebot.TextComponent{
															Type:	linebot.FlexComponentTypeText,
															Text:	"Place",
															Size:	linebot.FlexTextSizeTypeSm,
															Flex:	linebot.IntPtr(1),
															Color:	"#aaaaaa",
														},
														&linebot.TextComponent{
															Type:	linebot.FlexComponentTypeText,
															Text:	restaurants.Rest[i].Address,
															Wrap:	true,
															Size:	linebot.FlexTextSizeTypeSm,
															Flex:	linebot.IntPtr(5),
															Color:	"#666666",
														},
													},
													Spacing:  linebot.FlexComponentSpacingTypeSm,
												},
												&linebot.BoxComponent{
													Type:     linebot.FlexComponentTypeBox,
													Layout:   linebot.FlexBoxLayoutTypeBaseline,
													Contents: []linebot.FlexComponent{
														&linebot.TextComponent{
															Type:	linebot.FlexComponentTypeText,
															Text:	"Tel",
															Size:	linebot.FlexTextSizeTypeSm,
															Flex:	linebot.IntPtr(1),
															Color:	"#aaaaaa",
														},
														&linebot.TextComponent{
															Type:	linebot.FlexComponentTypeText,
															Text:	restaurants.Rest[i].Tel,
															Wrap:	true,
															Size:	linebot.FlexTextSizeTypeSm,
															Flex:	linebot.IntPtr(5),
															Color:	"#666666",
														},
													},
													Spacing:  linebot.FlexComponentSpacingTypeSm,
												},
											},
											Margin:   linebot.FlexComponentMarginTypeLg,
											Spacing:  linebot.FlexComponentSpacingTypeSm,
										},
									},
								},
								Footer:	&linebot.BoxComponent{
									Type:     linebot.FlexComponentTypeBox,
									Layout:   linebot.FlexBoxLayoutTypeVertical,
									Contents: []linebot.FlexComponent{
										&linebot.ButtonComponent{
											Type:    linebot.FlexComponentTypeButton,
											Action:  linebot.NewURIAction("WEBSITE", restaurants.Rest[i].URL),
											Height:  linebot.FlexButtonHeightTypeSm,
											Style:   linebot.FlexButtonStyleTypeLink,
										},
										&linebot.SpacerComponent{
											Type: linebot.FlexComponentTypeSpacer,
											Size: linebot.FlexSpacerSizeTypeSm,
										},
									},
									Spacing:  linebot.FlexComponentSpacingTypeSm,
									Flex:linebot.IntPtr(0),
								},
							}

							flexColumns = append(flexColumns, flexContent)
							if i == 5 {
								break
							}

						}
						carousel := &linebot.CarouselContainer{
							Type:     linebot.FlexContainerTypeCarousel,
							Contents: flexColumns,
						}
						_, err = bot.ReplyMessage(
							event.ReplyToken,
							linebot.NewTextMessage("はいいらっしゃいラーメンﾄﾞｰｰｰｰｰﾝ！！！！！！！！！食え食え食え食え食え食え食え食え食った？食った？食った？食った？食った？くった？食った？食った？食った？食った？食った？出てけ出てけ出てけ出てけ出てけ出てけ出てけ出てけ出てけ出てけ出てけ出てけ出てけ出てけ"),
							linebot.NewFlexMessage("alt", carousel),
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
