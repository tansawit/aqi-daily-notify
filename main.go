package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type airAQIResponse struct {
	Status string `json:"status"`
	Data   struct {
		Aqi          int `json:"aqi"`
		Idx          int `json:"idx"`
		Attributions []struct {
			URL  string `json:"url"`
			Name string `json:"name"`
			Logo string `json:"logo,omitempty"`
		} `json:"attributions"`
		City struct {
			Geo  []float64 `json:"geo"`
			Name string    `json:"name"`
			URL  string    `json:"url"`
		} `json:"city"`
		Dominentpol string `json:"dominentpol"`
		Iaqi        struct {
			Co struct {
				V float64 `json:"v"`
			} `json:"co"`
			H struct {
				V int `json:"v"`
			} `json:"h"`
			No2 struct {
				V float64 `json:"v"`
			} `json:"no2"`
			O3 struct {
				V float64 `json:"v"`
			} `json:"o3"`
			P struct {
				V int `json:"v"`
			} `json:"p"`
			Pm10 struct {
				V int `json:"v"`
			} `json:"pm10"`
			Pm25 struct {
				V int `json:"v"`
			} `json:"pm25"`
			So2 struct {
				V float64 `json:"v"`
			} `json:"so2"`
			T struct {
				V float64 `json:"v"`
			} `json:"t"`
			W struct {
				V float64 `json:"v"`
			} `json:"w"`
		} `json:"iaqi"`
		Time struct {
			S  string `json:"s"`
			Tz string `json:"tz"`
			V  int    `json:"v"`
		} `json:"time"`
		Debug struct {
			Sync time.Time `json:"sync"`
		} `json:"debug"`
	} `json:"data"`
}

type location struct {
	Name string  `json:"name"`
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

func getAQI(lat float64, long float64) int {
	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.waqi.info/feed/geo:%f;%f/?token=c68f5c2d123625f3b476ae72b80601ab992fbdd0", lat, long), nil)
	if err != nil {
		panic(err)
	}
	parseFormErr := req.ParseForm()
	if parseFormErr != nil {
		fmt.Println(parseFormErr)
	}

	var AQIResponse airAQIResponse

	resp, err := client.Do(req)

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&AQIResponse)
	if err != nil {
		panic(err)
	}
	return AQIResponse.Data.Iaqi.Pm25.V
}

func sendRequest(places []location) {
	apiURL := "https://notify-api.line.me/api/notify"
	data := url.Values{}
	message := "\nDaily AQI Update:"
	for _, loc := range places {
		message += fmt.Sprintf("\n- %s: %d", loc.Name, getAQI(loc.Lat, loc.Long))
	}
	data.Set("message", message)
	u, _ := url.ParseRequestURI(apiURL)
	urlStr := u.String()

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
	r.Header.Add("Authorization", "Bearer xsx74pLcGF3OpPywOWflz9rJGyc8x94Y1iR8IQTE5SO")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, _ := client.Do(r)
	fmt.Println(resp.Status)
}

func main() {
	var places = []location{location{"Condo", 13.722070, 100.534140}, location{"Work", 13.722590, 100.520400}, location{"Home", 13.794310, 100.397170}}
	sendRequest(places)
}
