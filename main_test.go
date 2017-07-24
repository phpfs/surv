package main

import (
	"testing"
)

func TestMethodSystemPing(t *testing.T){
	run1 := methodSystemPing("8.8.8.8")
	if(run1.Success){
		t.Log(run1, "\n\nmethodSystemPing Test: 8.8.8.8 - Successfull")
	}else{
		t.Error(run1, "\n\nmethodSystemPing Test: 8.8.8.8 - Failed")
	}

	run2 := methodSystemPing("192.168.168.192")
	if(!run2.Success){
		t.Log(run1, "\n\nmethodSystemPing Test: 192.168.168.192 - Failed")
	}else{
		t.Error(run1, "\n\nmethodSystemPing Test: 192.168.168.192 - Successfull")
	}
}

func TestMethodHTTP(t *testing.T){
	run1 := methodHTTP("https://google.com")
	if(run1.Success){
		t.Log(run1, "\n\nmethodHTTP Test: https://google.com - Successfull")
	}else{
		t.Error(run1, "\n\nmethodHTTP Test: https://google.com - Failed")
	}

	run2 := methodSystemPing("google.com")
	if(!run2.Success){
		t.Log(run1, "\n\nmethodHTTP Test: google.com - Failed")
	}else{
		t.Error(run1, "\n\nmethodHTTP Test: google.com - Successfull")
	}
}

func TestMethodTCP(t *testing.T){
	run1 := methodTCP("secureimap.t-online.de:993")
	if(run1.Success){
		t.Log(run1, "\n\nmethodTCP Test: secureimap.t-online.de:993 - Successfull")
	}else{
		t.Error(run1, "\n\nmethodTCP Test: secureimap.t-online.de:993 - Failed")
	}

	run2 := methodSystemPing("secureimap.t-online.de:65023")
	if(!run2.Success){
		t.Log(run1, "\n\nmethodTCP Test: secureimap.t-online.de:65023 - Failed")
	}else{
		t.Error(run1, "\n\nmethodTCP Test: secureimap.t-online.de:65023 - Successfull")
	}
}