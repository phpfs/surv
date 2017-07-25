package main

import (
	"fmt"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func syncServices(session *mgo.Session) bool {
	fmt.Println("Purging DB...")
	err := session.DB("surv").DropDatabase()

	fmt.Println("Uploading new Config to DB...")
	c := session.DB("surv").C("services")
	for _, service := range config.Services {
		service.Status = true
		service.Change = time.Now()
		err = c.Insert(service)
		if err != nil {
			fmt.Println(err)
			return false
		}
	}

	fmt.Println("Finished Upload successfully!")
	return true
}

func serviceStatus(session *mgo.Session, id string, status bool) {
	survs := session.DB("surv").C("services")

	var S Service

	err := survs.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&S)
	if err != nil {
		fmt.Println(err)
	}

	if S.Status != status {
		if status || time.Duration(S.THold)*time.Second < time.Since(S.Change) {
			go alert(S.Name, status)
		}

		S.Status = status
		S.Change = time.Now()
		err = survs.Update(bson.M{"_id": S.Id}, &S)
		if err != nil {
			fmt.Println(err)
		}
	}
}
