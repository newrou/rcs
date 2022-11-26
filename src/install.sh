#!/bin/bash

ln -s /root/Develop/rcs/src/rcs /var/lib/rcs

go build rcs-make-crontab.go
#./rcs-make-crontab.go
go build rcs-server.go
