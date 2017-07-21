package main

import (
	"gopkg.in/mgo.v2"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"time"
)

func runner(s *mgo.Session, id int){
	session := s.Clone()
	defer session.Close()

	fmt.Println("Worker ", id, " running...")

	tasks := session.DB("surv").C("tasks")

	var picked = mgo.Change{Update: bson.M{"status": "picked"}}

	for ; ; {
		var task Task
		err := tasks.Find(bson.M{"result": nil, "status": "pending"}).One(&task)
		if (err != nil && err != mgo.ErrNotFound) {
			fmt.Println(err)
		} else {
			if (err == mgo.ErrNotFound) {
				time.Sleep(time.Millisecond * 50)
			} else {
				_, err = tasks.Find(bson.M{"result": nil, "status": "pending", "_id": task.Id}).Apply(picked, &task)
				if(err == mgo.ErrNotFound){
					time.Sleep(50 * time.Millisecond)
				}else{
					task.Status = "running"
					err = tasks.Update(bson.M{"_id": task.Id}, &task)
					if (err != nil) {
						fmt.Println(err)
					}
					task.Result = method(task)

					if(!task.Result.Success){
						task.Result = nil
						time.Sleep(5 * time.Second)
						task.Result = method(task)
					}

					task.Status = "finished"
					task.Time = time.Now()
					task.Worker = id

					serviceStatus(session, task.Service, task.Result.Success)

					err = tasks.Update(bson.M{"_id": task.Id}, &task)
					if (err != nil) {
						fmt.Println(err)
					}
				}
			}
		}
	}
}