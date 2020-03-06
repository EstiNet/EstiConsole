# EstiConsole 
Version: v2.0.2

EstiConsole is a program that was developed to remotely control console processes. It was originally developed to manage Minecraft servers at [EstiNet](estinet.net), and is now used in production!

## How do I run it?
Grab a release from the [releases page](https://github.com/EstiNet/EstiConsole/releases), put the server binary for your platform in a folder. You can put the client binary anywhere that you want to run it.

Simply running the server binary will generate the necessary config files in the directory from which the server binary is run.

## Using the client
You can get the help page for the client by doing:

~~~~
$ ./client-binary -h
~~~~

## Generating a key pair for SSL
~~~~
$ openssl genrsa -out server.key 2048
$ openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
~~~~

## How do I build it?
EstiConsole requires a few dependencies to build (still uses old GOPATH way, will switch to go modules soooon).

~~~~
go get github.com/c9s/goprocinfo/linux
go get google.golang.org/grpc
go get github.com/jroimartin/gocui
go get github.com/nu7hatch/gouuid
go get github.com/howeyc/gopass
~~~~

Simply run the buildServer.sh and buildCli.sh scripts from the directory head of the repository to build the binaries.

## Why did you write it in Golang?
Because golang is cool.

## TODO
* Unzip file for extraction (and time delay to delete zip cache)
* Add configuration option to disable logging for process
* Add config reload

## Sample Systemd Unit File
~~~~
[Unit]
Description=EstiConsole

[Service]
WorkingDirectory=/opt/esticonsole
User=user

Restart=always
ExecStart=/opt/esticonsole/esticonsole
ExecStop=/opt/esticonsole/esticli -masterkey /opt/esticonsole/masterkey.key instancestop

[Install]
WantedBy=multi-user.target
~~~~
