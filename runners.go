package main

import (
	"gopkg.in/mgo.v2"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"time"
)

func runner(s *mgo.Session, id int){
	session := s.Copy()
	defer session.Close()

	fmt.Println("Worker ", id, " running...")

	tasks := session.DB("surv").C("tasks")

	var picked = mgo.Change{Update: bson.M{"status": "picked"}}

	for ; ; {
		var task Task
		err := tasks.Find(bson.M{"result": nil, "status": "pending"}).One(&task)
		if (err != nil && err != mgo.ErrNotFound) {
			fmt.Println("Runner", id, ": ", err)
			time.Sleep(20 * time.Second)
		} else {
			if (err == mgo.ErrNotFound) {
				time.Sleep(2 * time.Second)
			} else {
				_, err = tasks.Find(bson.M{"result": nil, "status": "pending", "_id": task.Id}).Apply(picked, &task)
				if(err == mgo.ErrNotFound){
					time.Sleep(time.Second)
				}else{
					task.Status = "running"
					err = tasks.Update(bson.M{"_id": task.Id}, &task)
					if (err != nil) {
						fmt.Println(err)
					}

					task.Result = method(task)

					i := 1
					for (i < 5 && !task.Result.Success){
						task.Result = nil
						time.Sleep(time.Duration(i) * 10 * time.Second)
						task.Result = method(task)
						i++
					}

					task.Status = "finished"
					task.Time = time.Now()
					task.Worker = id
					task.Tried = i

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