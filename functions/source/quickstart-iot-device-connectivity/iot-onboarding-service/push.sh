#!/bin/sh
env=$1
bucket=$2
echo "**********************************************"
echo "* Deploying Rigado onboarding service for environment '$env' and bucket $bucket"
echo "***********************************************"
if [ -z "$env" ] || [ -z "$bucket" ]
then
    echo "Environment and Bucket Must not be Empty"
    echo "Usage:"
    echo "sh deploy.sh <env> <bucket>"
else
echo "Publishing new function code"
aws s3 cp main.zip s3://$bucket/$env/iotOnBoarding/main.zip
aws lambda update-function-code --function-name iotOnBoarding$env --zip-file fileb://main.zip
fi
