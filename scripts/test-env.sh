#!/bin/bash

# yum install git gcc
# git clone https://github.com/sanderdescamps/terraform-provider-flash.git

# yum -y install https://dl.fedoraproject.org/pub/epel/epel-release-latest-7.noarch.rpm
# yum -y install jq

wget https://go.dev/dl/go1.19.linux-amd64.tar.gz

rm -rf /usr/local/go && tar -C /usr/local -xzf go1.19.linux-amd64.tar.gz

export PATH=$PATH:/usr/local/go/bin



export PURE_USERNAME=pureuser
export PURE_PASSWORD=pureuser
export PURE_TARGET=flasharray1.testdrive.local

TF_ACC=1 go test -run <test_name> -v ./...