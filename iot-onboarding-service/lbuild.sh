#!/bin/sh

env=$1
bucket=$2
#update this variable to specify the name of your loval env
LOCAL_ENV_NAME=local_dev_env

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

if [ $env == $LOCAL_ENV_NAME ]
    then
    echo "Local environement build"
    else
    sh ./env.sh $env
fi
echo "2-Unit Testing"
rm main main.zip
#in local we do not set the GOOS env variable to have a MACOS build
if [ $env == $LOCAL_ENV_NAME ]
    then
    echo "2.1-Local Unit Testing"
    go build -o main src/main/main.go
    sh ./test.sh
    rc=$?
    if [ $rc -ne 0 ]; then
      echo "Existing Build with status $rc" >&2
      exit $rc
    fi
fi
echo "3-Building application Deployement"
export GOOS=linux
go build -o main src/main/main.go
if [ $env != $LOCAL_ENV_NAME ]; then
    echo "3.1-Unit Testing"
    sh ./test.sh
    rc=$?
    if [ $rc -ne 0 ]; then
      echo "Existing Build with status $rc" >&2
      exit $rc
    fi
fi
zip main.zip main
echo "4-Deploying to Lambda"
sh push.sh $env $bucket
if [ $env == $LOCAL_ENV_NAME ]; then
  cd e2e && sh test.sh $env && cd ..
fi
fi
