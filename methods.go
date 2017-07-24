package main

import(
	"github.com/tatsushid/go-fastping"
	"net"
	"time"
	"net/http"
	"errors"
	"strings"
	"os/exec"
	"regexp"
	"runtime"
)

func method(task Task)(*mResult){
	if(task.Method == "methodFastPing"){
		return methodFastPing(task.Target, false)
	}else if(task.Method == "methodFastPingUDP"){
		return methodFastPing(task.Target, true)
	}else if(task.Method == "methodSystemPing"){
		return methodSystemPing(task.Target)
	}else if(task.Method == "methodHTTP"){
		return methodHTTP(task.Target)
	}else if(task.Method == "methodTCP"){
		return methodTCP(task.Target)
	}else {
		var result = new(mResult)
		result.Method = task.Method
		result.Count = 0
		result.Success = false
		result.Error = errors.New("Unknown method - no test was run!")
		return result
	}
}

func methodSystemPing(host string)(*mResult){
	var result = new(mResult)
	result.Method = "methodSystemPing"
	if(runtime.GOOS == "darwin" || runtime.GOOS == "linux") {
		cmd, err := exec.Command("ping", "-c 4", host).Output()
		if !(err != nil) {
			reg := regexp.MustCompile(`min\/avg\/max\/\w+ = (?P<min>\d+\.\d+)\/(?P<avg>\d+\.\d+)\/(?P<max>\d+\.\d+)\/(?P<mdev>\d+\.\d+) ms`)
			res := reg.FindAllStringSubmatch(string(cmd), -1)

			if len(res) > 0 {
				average, _ := time.ParseDuration(res[0][1] + "ms")
				if (average > 0) {
					result.Count = float64(average.Seconds())
					result.Success = true
					return result
				} else {
					result.Count = 0
					result.Success = false
					result.Error = errors.New("Error parsing average time!")
					return result
				}
			} else {
				result.Count = 0
				result.Success = false
				result.Error = errors.New("Error parsing `ping` command!")
				return result
			}
		}else{
			result.Method = "methodSystemPing"
			result.Count = 0
			result.Success = false
			result.Error = errors.New("Supplied IP or hostname was incorrect!")
			return result
		}
	}else{
		result.Method = "methodSystemPing"
		result.Count = 0
		result.Success = false
		result.Error = errors.New("You can only use this method on Linux or macOS!")
		return result
	}
}

func methodFastPing(ip string, udp bool)(*mResult){
	var dur time.Duration
	fin := false
	fail := false

	p := fastping.NewPinger()
	ra, err := net.ResolveIPAddr("ip4:icmp", ip)
	if err != nil {
		fail = true
	}
	p.AddIPAddr(ra)
	if(udp) {
		p.Network("udp")
	}
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
		resp, err := http.Get(url)
		dur := time.Since(start)

		if(err != nil || resp.StatusCode != 200){
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