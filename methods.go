package main

import(
	"github.com/tatsushid/go-fastping"
	"net"
	"time"
	"net/http"
	"errors"
	"strings"
)

func method(task Task)(*mResult){
	if(task.Method == "methodPing"){
		return methodPing(task.Target)
	}else if(task.Method == "methodHTTP"){
		return methodHTTP(task.Target)
	}else if(task.Method == "methodTCP"){
		return methodTCP(task.Target)
	}else {
		var result= new(mResult)
		result.Method = task.Method
		result.Count = 0
		result.Success = false
		result.Error = errors.New("Unknown method - no test was run!")
		return result
	}
}

func methodPing(ip string)(*mResult){
	var dur time.Duration
	fin := false
	fail := false

	p := fastping.NewPinger()
	ra, err := net.ResolveIPAddr("ip4:icmp", ip)
	if err != nil {
		fail = true
	}
	p.AddIPAddr(ra)
	//p.Network("udp")
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		//fmt.Printf("IP Addr: %s receive, RTT: %v\n", addr.String(), rtt)
		fin = true
		dur = rtt
	}
	p.OnIdle = func() {
		if(!fin){
			fail = true
		}
	}
	err = p.Run()
	if err != nil {
		fail = true
	}

	var result = new(mResult)
	result.Method = "methodPing"
	result.Count = float64(dur.Seconds())

	if(!fin || fail){
		result.Error = errors.New("Host couldn't be reached!")
		result.Success = false
	}else{
		result.Success = true
	}

	return result
}


func methodHTTP(url string)(*mResult){
	var result = new(mResult)
	result.Method = "methodHTTP"

	if(strings.Contains(url, "http://") || strings.Contains(url, "https://") || strings.Contains(url, "ftp://")) {
		start := time.Now()
		_, err := http.Get(url)
		dur := time.Since(start)

		if err != nil {
			result.Success = false
			result.Error = err
		} else {
			result.Success = true
		}
		result.Count = float64(dur.Seconds())
	}else{
		result.Success = false
		result.Error = errors.New("You URL is missing a protocol like http:// in front!")
		result.Count = 0
	}
	return result
}

func methodTCP(target string)(*mResult){
	var result = new(mResult)
	result.Method = "methodTCP"

	start := time.Now()
	_, err := net.Dial("tcp", target)
	dur := time.Since(start)

	result.Count = float64(dur.Seconds())

	if(err != nil){
		result.Success = false
		result.Error = err
	}else{
		result.Success = true
	}

	return result
}