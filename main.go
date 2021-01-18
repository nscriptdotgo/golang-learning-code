package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var token string = os.Getenv("TOKEN")
var port string = os.Getenv("PORT")

// Constants
const (
	BaseURL = "https://api.openweathermap.org/data/2.5"
	Unit    = "imperial"
)

// Coord struct
type Coord struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

// Weather struct
type Weather struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

// Main struct, contains the actual weather data
type Main struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like_a_burrito"`
	TempMin   float64 `json:"temp_min"`
	TempMax   float64 `json:"temp_max"`
	Pressure  int     `json:"pressure"`
	Humidity  int     `json:"humidity"`
}

// Wind struct, information about wind
type Wind struct {
	Speed float64 `json:"speed"`
	Deg   int     `json:"deg"`
}

// Clouds struct
type Clouds struct {
	All int `json:"all"`
}

// Sys struct
type Sys struct {
	Type    int     `json:"type"`
	ID      int     `json:"id"`
	Message float64 `json:"message"`
	Country string  `json:"country"`
	Sunrise int     `json:"sunrise"`
	Sunset  int     `json:"sunset"`
}

// CurrentWeatherData struct for getting current weather.
type CurrentWeatherData struct {
	Coord      Coord     `json:"coords"`
	Weather    []Weather `json:"weather"`
	Base       string    `json:"base"`
	Main       Main      `json:"main"`
	Visibility int       `json:"visibility"`
	Wind       Wind      `json:"wind"`
	Clouds     Clouds    `json:"clouds"`
	Dt         int       `json:"dt"`
	Sys        Sys       `json:"sys"`
	Timezone   int       `json:"timezone"`
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Cod        int       `json:"cod"`
}

func getCurrentWeatherByZipCode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Declare and intialize vars needed for api call, map mux route parameter
	token := os.Getenv("TOKEN")
	vars := mux.Vars(r)
	zipCode := vars["zipCode"]
	// Make API call with above vars, returns Response object
	resp, err := http.Get(fmt.Sprintf("%s/weather?zip=%s,us&units=%s&appid=%s", BaseURL, zipCode, Unit, token))
	log.Print(resp.Request)
	// catch error if error
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}
	// Take response object and covert that into a byte slice - we need this to store the values into custom struct for currentweatherdata
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(string(responseData))
	// declare weatherDataObject var as type CurrentWeatherData
	var weatherDataObject CurrentWeatherData
	// Unmarshal the byte slice and store into new weatherDataObject var
	json.Unmarshal(responseData, &weatherDataObject)

	json.NewEncoder(w).Encode(weatherDataObject)
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/currentweather/{zipCode}", getCurrentWeatherByZipCode)
	if os.Getenv("ENVIRONMENT") == "prod" {
		log.Fatal(http.ListenAndServe(":"+port, myRouter))
	} else {
		log.Fatal(http.ListenAndServe(":8080", myRouter))
	}
}

func main() {
	log.Print("==== Starting OpenWeather Golang Custom API ====")
	log.Print("Loading dotenv variables from file...")
	if os.Getenv("ENVIRONMENT") != "prod" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error: ", err)
		}
	}
	log.Print("Environment variables loaded!")
	log.Print("==== Started OpenWeather Golang Custom API ====")
	handleRequests()
}
