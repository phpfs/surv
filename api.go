package main

import (
	"net/http"
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func apiMain(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	out := make(map[string] string)
	out["message"] = "SurV " + survVersion + " is running!"
	w.WriteHeader(200)
	fin, _ := json.Marshal(out)
	fmt.Fprintf(w, string(fin))
}

func apiServices(s *mgo.Session, w http.ResponseWriter, r *http.Request){
	session := s.Clone()
	defer session.Close()
	survs := session.DB("surv").C("services")

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	out := make(map[string] string)

	var services []Service
	err := survs.Find(bson.M{}).All(&services)
	if(err != nil){
		out["err"] = "MongoDB query errored!"
		w.WriteHeader(500)
		fin, _ := json.Marshal(out)
		fmt.Fprintf(w, string(fin))
		return
	}else{
		w.WriteHeader(200)
		fin, _ := json.Marshal(services)
		fmt.Fprintf(w, string(fin))
	}
}
