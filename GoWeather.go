package main

import(
	"fmt"
	"log"
	"time"
	"os"
	
	"net/http"
	"net/url"
	"io/ioutil"
	"encoding/json"
)

	/*		Declarations		*/

//Structs
type ConfigData struct {
	OWM string `json:"openweathermap"`
	UnitsOWM string `json:"unitsOWM"`
	
	WB string `json:"weatherbit"`
	UnitsWB string `json:"unitsWB"`

	Lang string `json:"lang"`
}

type Weather struct {
	Description string
	Tempature string
	Humidity string
	Precipitation string
}

type Response struct {
	OMW *Weather
	WB *Weather
}

//Config
var config ConfigData

	/*		Functions		*/

//Init
func init() {
	getConfig()
	//Check if	
	if (config.OWM == "" && config.WB != "") {
		fmt.Println("Error, no API keys given")
		os.Exit(0)
	}
}

func main() {
	response := GetCurrent("Berlin,DE")
	fmt.Println(response)
}

//Util
func getConfig() {
	data, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
    }
	
    err = json.Unmarshal(data, &config)
    if err != nil {
        log.Fatal("Error during Unmarshal(): ", err)
    }
}

//Actions
func GetCurrent(search string) *Response {
	newResponse := new(Response)

	lat, lon := geocodeMapsAPI(search)
	if config.OWM != "" {
		newResponse.OMW = openWeatherMapAPI(lat, lon)
	}
	if config.WB != "" {
		newResponse.WB = weatherBitAPI(lat, lon)
	}
	
	return newResponse
}

//APIs
func openWeatherMapAPI(lat string, lon string) *Weather {
	//Generic GET request to Open Weather Map API
	client := &http.Client{
		Timeout: time.Second * 10,
    }
    
    url := "https://api.openweathermap.org/data/2.5/weather?lat=" + lat + "&lon=" + lon + "&lang=" + config.Lang + "&units=" + config.UnitsOWM + "&APPID=" + config.OWM
    
    req, err := http.NewRequest("GET", url, nil)
    
    if err != nil {
		log.Fatal(err)
	}
	
	resp, err := client.Do(req)
	
	if err != nil {
        log.Fatal(err)
    }
    
    responseData, err := ioutil.ReadAll(resp.Body)
    
    if err != nil {
        log.Fatal(err)
    }
    
	defer resp.Body.Close()
	
	var data map[string]json.RawMessage 
	if err := json.Unmarshal(responseData, &data); err != nil {
		log.Fatal(err)
	}
	
	newWeather := new(Weather)
	
	//Unstructured JSON for description
	descData := data["weather"]
	var descObj []map[string]interface{}
	if err := json.Unmarshal(descData, &descObj); err != nil {
		log.Fatal(err)
	}
	newWeather.Description = fmt.Sprintf("%v", descObj[0]["description"])
	
	//Unstructured JSON for tempature & humidity
	tempData := data["main"]
	var tempObj map[string]interface{}
	if err := json.Unmarshal(tempData, &tempObj); err != nil {
		log.Fatal(err)
	}
	newWeather.Tempature = fmt.Sprintf("%v", tempObj["temp"])
	newWeather.Humidity = fmt.Sprintf("%v", tempObj["humidity"])
	newWeather.Precipitation = "N/A"
	
	return newWeather
}

func weatherBitAPI(lat string, lon string) *Weather {
	client := &http.Client{
		Timeout: time.Second * 10,
    }
    
    url := "http://api.weatherbit.io/v2.0/current?lang=" + config.Lang + "&units=" + config.UnitsWB + "&lon=" + lon + "&lat=" + lat + "&include=minutely&key=" + config.WB
    
    req, err := http.NewRequest("GET", url, nil)
    
    if err != nil {
		log.Fatal(err)
	}
	
	resp, err := client.Do(req)
	
	if err != nil {
        log.Fatal(err)
    }
    
    responseData, err := ioutil.ReadAll(resp.Body)
    
    if err != nil {
        log.Fatal(err)
    }
    
	defer resp.Body.Close()
	
	var data map[string]json.RawMessage 
	if err := json.Unmarshal(responseData, &data); err != nil {
		log.Fatal(err)
	}
	
	newWeather := new(Weather)
	
	//Unstructured Data
	weatherData := data["data"]
	var weatherObj []map[string]interface{}
	if err := json.Unmarshal(weatherData, &weatherObj); err != nil {
		log.Fatal(err)
	}
	
	newWeather.Tempature = fmt.Sprintf("%v", weatherObj[0]["temp"])
	newWeather.Humidity = fmt.Sprintf("%v", weatherObj[0]["rh"])
	newWeather.Precipitation = fmt.Sprintf("%v", weatherObj[0]["precip"])

	
	//Have to do this again cause need raw Json and not an interface
	weatherData2 := data["data"]
	var weatherObj2 []map[string]json.RawMessage 
	if err := json.Unmarshal(weatherData2, &weatherObj2); err != nil {
		log.Fatal(err)
	}
	
	//Unstructured Data for Description
	descData := weatherObj2[0]["weather"]
	var descObj map[string]interface{}
	if err := json.Unmarshal(descData, &descObj); err != nil {
		log.Fatal(err)
	}
	newWeather.Description = fmt.Sprintf("%v", descObj["description"])
	
	return newWeather
}

func geocodeMapsAPI(search string) (string, string) {
	//Generic GET request to Open Weather Map API
	client := &http.Client{
		Timeout: time.Second * 10,
    }
    
    query := url.QueryEscape(search)
    url := "https://geocode.maps.co/search?q=" + query
        
    req, err := http.NewRequest("GET", url, nil)
    
    if err != nil {
		log.Fatal(err)
	}
	
	resp, err := client.Do(req)
	
	if err != nil {
        log.Fatal(err)
    }
    
    responseData, err := ioutil.ReadAll(resp.Body)
    
    if err != nil {
        log.Fatal(err)
    }
    
	defer resp.Body.Close()
	
	var data []map[string]interface{}
	if err := json.Unmarshal(responseData, &data); err != nil {
		log.Fatal(err)
	}
	
	lat := fmt.Sprintf("%v", data[0]["lat"])
	lon := fmt.Sprintf("%v", data[0]["lon"])
	
	return lat, lon
}