package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

func schedule(s *mgo.Session){
	session := s.Clone()
	defer session.Close()

	survs := session.DB("surv").C("services")
	tasks := session.DB("surv").C("tasks")

	for ; ;  {
		var result []Service

		err := survs.Find(bson.M{}).All(&result)
		if(err != nil){
			fmt.Println(err)
		}

		for _, S := range result {
			if(time.Since(S.Last) > time.Duration(S.Cron.Every) * time.Second){
				task := new(Task)
				task.Service = S.Id.Hex()
				task.Method = S.Method
				task.Target = S.Target
				task.Name = S.Name
				task.Status = "pending"

				err := tasks.Insert(task)
				if(err != nil){
					fmt.Print(err)
				}

				S.Last = time.Now()
				err = survs.Update(bson.M{"_id": S.Id}, &S)
				if (err != nil) {
					fmt.Println(err)
				}
			}
		}
	}
}