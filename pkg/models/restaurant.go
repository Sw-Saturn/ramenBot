package models

import "time"

//GNavi is that can be acquired from the GNavi API
type GNavi struct {
	Rest []struct {
		ID         string    `json:"id"`
		UpdateDate time.Time `json:"update_date"`
		Name       string    `json:"name"`
		NameKana   string    `json:"name_kana"`
		Latitude   string    `json:"latitude"`
		Longitude  string    `json:"longitude"`
		Category   string    `json:"category"`
		URL        string    `json:"url"`
		URLMobile  string    `json:"url_mobile"`
		CouponURL  struct {
			Pc     string `json:"pc"`
			Mobile string `json:"mobile"`
		} `json:"coupon_url"`
		ImageURL struct {
			ShopImage1 string `json:"shop_image1"`
			ShopImage2 string `json:"shop_image2"`
			Qrcode     string `json:"qrcode"`
		} `json:"image_url"`
		Address  string `json:"address"`
		Tel      string `json:"tel"`
		TelSub   string `json:"tel_sub"`
		Fax      string `json:"fax"`
		Opentime string `json:"opentime"`
		Holiday  string `json:"holiday"`
		Access   struct {
			Line        string `json:"line"`
			Station     string `json:"station"`
			StationExit string `json:"station_exit"`
			Walk        string `json:"walk"`
			Note        string `json:"note"`
		} `json:"access"`
		ParkingLots string `json:"parking_lots"`
		Pr          struct {
			PrShort string `json:"pr_short"`
			PrLong  string `json:"pr_long"`
		} `json:"pr"`
		Code struct {
			Areacode      string   `json:"areacode"`
			Areaname      string   `json:"areaname"`
			Prefcode      string   `json:"prefcode"`
			Prefname      string   `json:"prefname"`
			AreacodeS     string   `json:"areacode_s"`
			AreanameS     string   `json:"areaname_s"`
			CategoryCodeL []string `json:"category_code_l"`
			CategoryNameL []string `json:"category_name_l"`
			CategoryCodeS []string `json:"category_code_s"`
			CategoryNameS []string `json:"category_name_s"`
		} `json:"code"`
	} `json:"rest"`
}
