#!/bin/bash

test -e /etc/systemd/system/udns.service && systemctl stop udns && systemctl disable udns && rm -f /etc/systemd/system/udns.service || echo "no udns.service"
test -e /etc/init/udns.conf && (stop udns; rm -f /etc/init/udns.conf) || echo "no udns.conf"
