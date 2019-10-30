package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/linebot"
	"log"
	"m/pkg/client"
	"net/http"
	"os"
)

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
					restaurants := client.RetrieveRestaurants(lat, lon)

					if len(restaurants.Rest) > 0 {
						var flexColumns []*linebot.BubbleContainer
						for i := range restaurants.Rest {
							access := fmt.Sprintf("不明")
							if restaurants.Rest[i].Access.Walk != "" {
								access = fmt.Sprintf("%s %s %s 分", restaurants.Rest[i].Access.Line, restaurants.Rest[i].Access.Station, restaurants.Rest[i].Access.Walk)
							}
							var flexContent = &linebot.BubbleContainer{
								Type: linebot.FlexContainerTypeBubble,
								Body: &linebot.BoxComponent{
									Type:   linebot.FlexComponentTypeBox,
									Layout: linebot.FlexBoxLayoutTypeVertical,
									Contents: []linebot.FlexComponent{
										&linebot.TextComponent{
											Type:   linebot.FlexComponentTypeText,
											Text:   restaurants.Rest[i].Name,
											Weight: linebot.FlexTextWeightTypeBold,
											Size:   linebot.FlexTextSizeTypeXl,
										},
										&linebot.BoxComponent{
											Type:   linebot.FlexComponentTypeBox,
											Layout: linebot.FlexBoxLayoutTypeVertical,
											Contents: []linebot.FlexComponent{
												&linebot.BoxComponent{
													Type:   linebot.FlexComponentTypeBox,
													Layout: linebot.FlexBoxLayoutTypeBaseline,
													Contents: []linebot.FlexComponent{
														&linebot.TextComponent{
															Type:  linebot.FlexComponentTypeText,
															Text:  "Access",
															Size:  linebot.FlexTextSizeTypeSm,
															Flex:  linebot.IntPtr(1),
															Color: "#aaaaaa",
														},
														&linebot.TextComponent{
															Type:  linebot.FlexComponentTypeText,
															Text:  access,
															Wrap:  true,
															Size:  linebot.FlexTextSizeTypeSm,
															Flex:  linebot.IntPtr(5),
															Color: "#666666",
														},
													},
													Spacing: linebot.FlexComponentSpacingTypeSm,
												},
												&linebot.BoxComponent{
													Type:   linebot.FlexComponentTypeBox,
													Layout: linebot.FlexBoxLayoutTypeBaseline,
													Contents: []linebot.FlexComponent{
														&linebot.TextComponent{
															Type:  linebot.FlexComponentTypeText,
															Text:  "Place",
															Size:  linebot.FlexTextSizeTypeSm,
															Flex:  linebot.IntPtr(1),
															Color: "#aaaaaa",
														},
														&linebot.TextComponent{
															Type:  linebot.FlexComponentTypeText,
															Text:  restaurants.Rest[i].Address,
															Wrap:  true,
															Size:  linebot.FlexTextSizeTypeSm,
															Flex:  linebot.IntPtr(5),
															Color: "#666666",
														},
													},
													Spacing: linebot.FlexComponentSpacingTypeSm,
												},
												&linebot.BoxComponent{
													Type:   linebot.FlexComponentTypeBox,
													Layout: linebot.FlexBoxLayoutTypeBaseline,
													Contents: []linebot.FlexComponent{
														&linebot.TextComponent{
															Type:  linebot.FlexComponentTypeText,
															Text:  "Tel",
															Size:  linebot.FlexTextSizeTypeSm,
															Flex:  linebot.IntPtr(1),
															Color: "#aaaaaa",
														},
														&linebot.TextComponent{
															Type:  linebot.FlexComponentTypeText,
															Text:  restaurants.Rest[i].Tel,
															Wrap:  true,
															Size:  linebot.FlexTextSizeTypeSm,
															Flex:  linebot.IntPtr(5),
															Color: "#666666",
														},
													},
													Spacing: linebot.FlexComponentSpacingTypeSm,
												},
											},
											Margin:  linebot.FlexComponentMarginTypeLg,
											Spacing: linebot.FlexComponentSpacingTypeSm,
										},
									},
								},
								Footer: &linebot.BoxComponent{
									Type:   linebot.FlexComponentTypeBox,
									Layout: linebot.FlexBoxLayoutTypeVertical,
									Contents: []linebot.FlexComponent{
										&linebot.ButtonComponent{
											Type:   linebot.FlexComponentTypeButton,
											Action: linebot.NewURIAction("WEBSITE", restaurants.Rest[i].URL),
											Height: linebot.FlexButtonHeightTypeSm,
											Style:  linebot.FlexButtonStyleTypeLink,
										},
										&linebot.SpacerComponent{
											Type: linebot.FlexComponentTypeSpacer,
											Size: linebot.FlexSpacerSizeTypeSm,
										},
									},
									Spacing: linebot.FlexComponentSpacingTypeSm,
									Flex:    linebot.IntPtr(0),
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
						if len(message.Address) > 0 {
							_, err = bot.ReplyMessage(
								event.ReplyToken,
								linebot.NewLocationMessage(
									"近くにラーメン屋が見つからなかったよ😢",
									message.Address,
									message.Latitude,
									message.Longitude,
								),
							).Do()
						} else {
							_, err = bot.ReplyMessage(
								event.ReplyToken,
								linebot.NewLocationMessage(
									"近くにラーメン屋が見つからなかったよ😢",
									"住所不明",
									message.Latitude,
									message.Longitude,
								),
							).Do()
						}
						if err != nil {
							log.Fatal(err)
						}
					}
				}
			}
		}
	})

	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}
}
