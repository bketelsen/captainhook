# captainhook

A generic webhook endpoint that runs scripts based on the URL called

This tool was built as part of a CI orchestration process, to be called when
Docker trusted builds finish.  It explicitly ignores the posted data from the webhook
because that would be `insecure`, which is `bad`. 


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

## Install

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

- consider making the POST data from the webhook available in some way
- more logs


Copyright 2014, Brian Ketelsen and Kelsey Hightower
LICENSE information found in LICENSE file.
