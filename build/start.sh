#!/bin/bash

# start
cd /opt/udns
test -e /etc/systemd/system && (
    cp -f udns.service /etc/systemd/system/udns.service && systemctl enable udns && systemctl restart udns  # systemd
    ) || (
    cp -f udns.conf /etc/init/udns.conf && (stop udns; start udns)  # upstart
    )
