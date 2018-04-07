# EstiConsole
Version: v2.0.0

## What the heck is this thing?
EstiConsole is a program that was developed to remotely control console processes.

## How do I run it?

~~~~
$ openssl genrsa -out server.key 2048
$ openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
~~~~

## How do I build it?
EstiConsole requires a few dependencies to build.
TEMP
~~~~
go get github.com/c9s/goprocinfo/linux
go get google.golang.org/grpc
go get github.com/jroimartin/gocui
go get github.com/nu7hatch/gouuid
go get github.com/howeyc/gopass
~~~~

## Why did you write it in Golang?
Because golang is cool.
