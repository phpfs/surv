package main

import (
	"testing"
)

func TestMethodPing(t *testing.T){
	run1 := methodPing("8.8.8.8")
	if(run1.Success){
		t.Log(run1, "\n\nmethodPing Test: 8.8.8.8 - Successfull")
	}else{
		t.Error(run1, "\n\nmethodPing Test: 8.8.8.8 - Failed")
	}

	run2 := methodPing("192.168.168.192")
	if(!run2.Success){
		t.Log(run1, "\n\nmethodPing Test: 192.168.168.192 - Failed")
	}else{
		t.Error(run1, "\n\nmethodPing Test: 192.168.168.192 - Successfull")
	}
}