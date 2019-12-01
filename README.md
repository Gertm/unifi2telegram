# unifi2telegram
Send camera snapshot to telegram on motion.

This little program monitors the recording.log file made by unifi-video.  
All you have to do is configure where the cameras are in your network and set up a Telegram bot. ( https://core.telegram.org/bots )  

Make sure you enabled the snapshotting feature in your unifi cameras, you can set this up through the nvr dashboard.

WARNING: This is still largely untested, although I am using it myself and it works pretty well for me.

## How to build

* Make sure golang is installed
* Clone the repo
* go get -u ./...
* go build

Copy the resulting executable to /usr/bin/

There is an example service file for systemd included.

To cross-compile for Windows, use the with-env.ps1 Powershell file to do this:
```
with-env.ps1 GOOS=linux GOARCH=amd64 go build
```
