# EstiConsole README WIP
Version: v2.0.1

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

## TODO
* Unzip file for extraction (and time delay to delete zip cache)
* Add configuration option to disable logging for process
* Add config reload

systemd unit file
~~~~
[Unit]
Description=EstiConsole

[Service]
WorkingDirectory=/home/estinet/EstiConsole
User=estinet

Restart=always
ExecStart=/home/estinet/EstiConsole/esticonsole
ExecStop=/home/estinet/EstiConsole/esticli -masterkey /home/estinet/EstiConsole/masterkey.key instancestop

[Install]
WantedBy=multi-user.target
~~~~
