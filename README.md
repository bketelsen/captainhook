# captainhook

[![Build Status](https://drone.gopheracademy.com/api/badges/bketelsen/captainhook/status.svg)](https://drone.gopheracademy.com/bketelsen/captainhook)

A generic webhook endpoint that runs scripts based on the URL called

This tool was built as part of a CI orchestration process, to be called when
Docker trusted builds finish.  It explicitly ignores the posted data from the webhook
because that would be `insecure`, which is `bad`. 

## Shoulders of Giants

Captainhook would not be possible if not for all of the great projects it depends on. Please see [SHOULDERS.md](SHOULDERS.md) to see a list of them.

## Quick Start

### Install captainhook

`go get github.com/bketelsen/captainhook`

### Create the `configdir`

```
mkdir ~/captainhook
```

### Add a script 

```
{
    "scripts": [
        {
            "command": "ls",
            "args": [
                "-l",
                "-a"
            ]
        },
        {
            "command": "echo",
            "args": [
		    "hello"
		    ]
        }
    ]
}
```
Name this script `endpoint1.json`

### Start the service

```
captainhook -configdir ~/captainhook
```

### Test using curl

```
curl -X POST http://localhost:8080/endpoint1
```

### Configure calling webhooks
Each script you create in the `configdir` will be executed when
the corresponding endpoint is called. 

If you have a script called `deployBigApp.json` you would trigger
it by posting to http://your.captainhook.url/deployBigApp.

The scripts in the json file are executed sequentially, and the output is logged
and returned to the caller in the response, which always has an HTTP status code
of 200 (OK) even if your scripts didn't work.  This is intentional, to avoid causing
errors in external services like Docker or Github, which might not like you returning
statuses other than 200 (OK).

### Accessing the Request POST Body 
You'll sometimes need to access the POST data of the request for information such as a callback URL. 
You can pass the raw POST data to a script by adding {{POST}} to the script arguments.

```json
{
    "scripts": [
        {
            "command": "echo",
            "args": [
            "{{POST}}"
            ]
        }
    ]
}
``` 

### Limiting access for webhooks
You can limit who can call your webhooks by specifying "allowedNetworks" in the json config.

```json
{
    "scripts": [
        {
            "command": "echo"
        }
    ],
    "allowedNetworks": [
        "10.0.0.0/8",
        "127.0.0.1/32"
    ]
}
```
This would allow your hook to be called from the 10.0.0.0/8 network, or from localhost.

### Supporting proxy headers for client IP

Only enable proxy support if you are on a trusted network behind a reverse proxy. End-users with direct network access can subvert the allowedNetworks restriction if proxy support is on.
```
captainhook -configdir ~/captainhook -enable-proxy -proxy-header X-Forwarded-For
```

## Docker
```
docker pull bketelsen/captainhook
mkdir /some/local/config
$EDITOR /some/local/config/myhook.json
docker run -d -v /some/local/config:/config bketelsen/captainhook
```

## Install

captainhook requires Go 1.10 to build locally with 'vgo' support.

`go get github.com/bketelsen/captainhook`

## Build

Download

```
mkdir -p $GOPATH/src/github.com/bketelsen
cd $GOPATH/src/github.com/bketelsen
git clone git@github.com:bketelsen/captainhook.git
```

```
go build .
```
## To Do

- more logs


Copyright 2014, Brian Ketelsen and Kelsey Hightower
LICENSE information found in LICENSE file.
