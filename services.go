package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func syncServices(s *mgo.Session) bool {
	session := s.Copy()
	defer session.Close()

	fmt.Println("Purging DB...")
	err := session.DB("surv").C("services").DropCollection()

	fmt.Println("Uploading new Config to DB...")
	c := session.DB("surv").C("services")
	for  _, service := range config.Services {
		service.Status = true
		err = c.Insert(service)
		if(err != nil){
			fmt.Println(err)
			return false
		}
	}

	fmt.Println("Finished Upload successfully!")
	return true
}

func serviceStatus(s *mgo.Session, id string, status bool){
	session := s.Copy()
	defer session.Close()

	survs := session.DB("surv").C("services")

	var S Service

	err := survs.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&S)
	if(err != nil){
		fmt.Println(err)
	}

	if(S.Status != status){
		go alert(S.Name, status)

		S.Status = status
		err = survs.Update(bson.M{"_id": S.Id}, &S)
		if (err != nil) {
			fmt.Println(err)
		}
	}
}
