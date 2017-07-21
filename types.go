package main

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type (
	mResult struct {
		Method string `json:"method"`
		Count float64 `json:"count"`
		Success bool `json:"success"`
		Error error `json:"error"`
	}

	Task struct {
		Id bson.ObjectId `json:"id" bson:"_id,omitempty"`
		Name string `json:"name"`
		Service string `json:"service"`
		Method string `json:"method"`
		Target string `json:"target"`
		Status string `json:"status"`
		Time time.Time `json:"time"`
		Worker int `json:"worker"`
		Result *mResult `json:"result"`
	}

	Service struct {
		Id bson.ObjectId `json:"id" bson:"_id,omitempty"`
		Name string `json:"name"`
		Cron Cron `json:"cron"`
		Target string `json:"target"`
		Method string `json:"method"`
		Last time.Time `json:"last"`
		Status bool `json:"status"`
	}

	Config struct {
		NumWorkers int
		Mongodb string
		Token string
		ApiPort string
		WebPort string
		Services []Service `json:"services"`
		Alert Alert
	}

	Alert struct {
		Typ string
		Target string
		Auth string
	}

	Cron struct {
		Every int `json:"every"`
	}
)