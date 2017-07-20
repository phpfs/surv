# SurV
Have you ever wondered why monitoring has to be so hard?
I mean - don't get me wrong - Prometheus and Nagios are great, but sometimes just a little overpowered.    
That's why I created SurV, an application written fully in GoLang, to easily monitor your services!

## Installation
Make sure, you have MongoDB installed and running!   

After you added your services to `config.toml` (see Configuration), it's as easy as:
```bash
git clone https://github.com/phpfs/surv.git
cd surv
go build
```
You now have an executable named `./surv` that you can run!

## Configuration
If you want to change the path to you config.toml, edit config.go and then `go build` again!   

Besides that, all configuration is made in config.toml:
```toml
numWorkers = 4
mongodb = "localhost"
apiPort = "9010"
webPort = "9020"

[alert]
typ = "alertTelegram"
target = "Chat_ID"
auth = "Token-For-HTTP-API"

[[services]]
name = "Test Google"
target = "8.8.8.8"
method = "methodPing"
[services.cron]
every = 60

[[services]]
name = "Test Localhost"
target = "http://localhost:8000"
method = "methodHTTP"
[services.cron]
every = 10

[[services]]
name = "Test Bing"
target = "https://bing.de"
method = "methodHTTP"
[services.cron]
every = 90
```

---> Let's look at this step by step: 
```toml
numWorkers = 4
mongodb = "localhost"
apiPort = "9010"
webPort = "9020"
```
In this part, you can specify the number of Workers which will work on your monitoring tasks. This number should be about the same as the number of services you want to monitor!   
Your MongoDB string should be a complete URL containing authentification, ports and hostnames.
At last, you have to specify the ports on which the API and the Web-Page (Coming soon!) are listening (see API for more details!).
```toml
[alert]
typ = "alertTelegram"
target = "Chat_ID"
auth = "Token-For-HTTP-API"
```
This Example is prepared for use with a telegram bot as the notificator.   
If you just want to receive alert on the Command Line, set `typ` to "alertCMD" and leave `target` and `auth` blank!
(You can easily discover your ChatID by contacting @cid_bot on Telegram!)
```toml
[[services]]
name = "Test Google"
target = "8.8.8.8"
method = "methodPing"
[services.cron]
every = 60
```
For each service you want to monitor, you simply repeat this paragraph.
Specify a `name` for your service, a fitting `target`and choose a `method` (see Methods) to use.   
At last, specify an interval named `every` after which your service should be rechecked!


## Methods
Currently, SurV supports 3 methods to check a service's availability:
1. -methodPing-
```toml
[[services]]
name = "Test Google"
target = "8.8.8.8"
method = "methodPing"
[services.cron]
every = 60
```
It is important that you use a plain IP as the `target`!     
2. -methodHTTP-
```toml
[[services]]
name = "Test Localhost"
target = "http://localhost:8000"
method = "methodHTTP"
[services.cron]
every = 10
```
Make sure to supplie a valid URL containing a protocol like http:// or ftp://!     
3. -methodTCP-
```toml
[[services]]
name = "IMAP"
target = "secureimap.t-online.de:993"
method = "methodTCP"
[services.cron]
every = 50
```
SurV will try to establish a TCP connection with service listening on the port you specified!
This time, you have to supplie a domain or IP with the fitting port!

## API
Endpoints:
- `/` - Check if SurV-API is running
- `/services` - GET all services that registered and monitored including their current status and their last check timestamp!
