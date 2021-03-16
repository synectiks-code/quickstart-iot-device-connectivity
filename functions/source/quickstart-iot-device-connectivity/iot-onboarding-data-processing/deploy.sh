
# Copyright 2019 Amazon.com, Inc. or its affiliates. All Rights Reserved.
#  SPDX-License-Identifier: MIT-0

env=$1
bucket=$2
stackName=$3
#update this variable to specify the name of your loval env
LOCAL_ENV_NAME=dev

echo "**********************************************"
echo "* IOT Onboarding Data Platform Glue ETL '$env' and artifact bucket '$bucket'" 
echo "***********************************************"
if [ -z "$env" ] || [ -z "$bucket" ]
then
    echo "Environment or Bucket Must not be Empty"
    echo "Usage:"
    echo "sh deploy.sh <env>"
else
    OUTRegion=$(aws configure get region)
    echo "0-Pushing iotOnboardingSensorFlatteningJob.py script to S3"
    aws s3 cp iotOnboardingSensorFlatteningJob.py s3://$bucket/$env/etl/iotOnboardingSensorFlatteningJob.py
fi
