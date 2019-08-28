# Ultra simple DNS server support file storage

[![Build Status](https://travis-ci.org/feiyuw/udns.svg?branch=master)](https://travis-ci.org/feiyuw/udns)

## Features

- [x] support A record
- [x] support to read DNS records from local file
- [x] support to forward to parent DNS server if not found
- [x] support to fetch myip by "my.ip"
- [x] support to reload DNS records when local file changed
- [ ] support CNAME record

## Installation (RHEL/CentOS)

1. Download rpm from [udns-0.1.0-1.x86_64.rpm](https://github.com/feiyuw/udns/releases/download/0.1.0/udns-0.1.0-1.x86_64.rpm)
1. yum install udns-0.1.0-1.x86_64.rpm
1. set DNS records at /opt/udns/local.dns
1. update /etc/resolv.conf
1. enjoy
