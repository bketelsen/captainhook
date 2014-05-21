# captainhook


A generic webhook endpoint that runs scripts based on the URL called

## Quick Start

Create the `configdir`

```
mkdir ~/captainhook
```

Add a script 

```
```

```
captainhook httpd -configdir=`pwd`
```

from a browser or curl, call http://localhost:8080/endpoint1

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
