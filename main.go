package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"time"
	"io/ioutil"
	"github.com/BurntSushi/toml"
	"net/http"
	"gopkg.in/alecthomas/kingpin.v2"
)

var config *Config
var alertAPI AlertAPI

var (
	flagSync = kingpin.Flag("sync", "Sync services from your config.toml - (Attention: The complete database will be purged!)").Default("false").Short('s').Bool()
	flagHTTP = kingpin.Flag("http", "Prevent HTTP webinterface and json-api from starting!").Default("false").Short('h').Bool()
)

func main() {
	kingpin.Version(survVersion)
	kingpin.Parse()

	fmt.Println("Welcome to |SurV|", survVersion, "- starting up...")

	if(!readConfig(configFile)){
		panic("--> Sync wasn't successfull!")
	}
	
	session, err := mgo.Dial(config.Mongodb)
	defer session.Close()
	if(err != nil){
		fmt.Println(err)
	}

	if(*flagSync){
		if(!syncServices(session)){
			panic("--> Sync wasn't successfull!")
		}
	}

	startAlert()

	go createRunners(session)

	if(!(*flagHTTP)){
		go startHTTP(session)
	}

	scheduleLoop(session)
}

func scheduleLoop(s *mgo.Session){
	fmt.Println("--> Starting schedulement loop...")
	for ; ;  {
		time.Sleep(time.Second)
		schedule(s)
	}
}

func createRunners(s *mgo.Session){
	for i := 1; i <= config.NumWorkers; i++ {
		go runner(s, i)
		time.Sleep(time.Millisecond)
	}
}

func readConfig(conf string) bool {
	fmt.Println("--> Reading Config file ", conf, "...")
	rawData, err := ioutil.ReadFile(conf)
	if(err != nil){
		fmt.Println(err)
		return false
	}
	tomlData := string(rawData)

	fmt.Println("--> Parsing TOML...")
	if _, err := toml.Decode(tomlData, &config); err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func startHTTP(s *mgo.Session){
	fmt.Println("--> Starting HTTP API...")
	http.HandleFunc("/", apiMain)
	http.HandleFunc("/services", func(w http.ResponseWriter, r *http.Request) {
		apiServices(s, w, r)
	})
	http.ListenAndServe(":" + config.ApiPort, nil)
}
