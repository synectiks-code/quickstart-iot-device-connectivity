#!/bin/sh

env=$1
bucket=$2

echo "**********************************************"
echo "* IOT Onboarding Service deployement for env '$env' "
echo "***********************************************"
if [ -z "$env" ] || [ -z "$bucket" ]
then
    echo "Environment Must not be Empty"
    echo "Usage:"
    echo "sh deploy.sh <env> <bucket>"
else
echo "1-Setting up environement"
export GOPATH=$(pwd)
echo "GOPATH set to $GOPATH"
echo "Instaling dependencies"
echo "2-Unit Testing"
rm main main.zip
echo "3-Building application Deployement"
export GOOS=linux
go build -o main src/main/main.go
echo "3.1-Unit Testing"
sh ./test.sh
rc=$?
if [ $rc -ne 0 ]; then
  echo "Existing Build with status $rc" >&2
  exit $rc
fi

zip main.zip main
echo "4-Deploying to Lambda"
sh push.sh $env $bucket
cd e2e && sh test.sh $env && cd ..
fi
