#!/bin/sh

env=$1
#update this variable to specify the name of your loval env
LOCAL_ENV_NAME=dev

echo "**********************************************"
echo "* IOT Onboarding: Create Template deployement for env '$env' "
echo "***********************************************"
if [ -z "$env" ]
then
    echo "Environment Must not be Empty"
    echo "Usage:"
    echo "sh deploy.sh <env>"
else
echo "0-Getting account ID"
rm -rf out
mkdir out
aws sts  get-caller-identity > out/id.json
cat out/id.json
accountId=$(jq -r .Account out/id.json)
analysis=4ddbea29-6ac7-43f5-8fcb-b4222e02582a
dataset=0cd2431d-ac85-46bf-8609-8b5d7293765c
region=$(aws configure get region)

echo "1-Creating requests for account $accountId"
cp 0-templateFromAnalysis.json.tpl out/0-templateFromAnalysis.json
cp 0-templatePermission.json out/0-templatePermission.json
if [ $env == $LOCAL_ENV_NAME ]
    then
    sed -i '' "s/__ACCOUNT_ID/$accountId/g" out/0-templateFromAnalysis.json
    sed -i '' "s/__ANALYSIS_ID/$analysis/g" out/0-templateFromAnalysis.json
    sed -i '' "s/__DATASET_ID/$dataset/g" out/0-templateFromAnalysis.json
    sed -i '' "s/__AWS_REGION/$region/g" out/0-templateFromAnalysis.json 
    else
    sed -i "s/__ACCOUNT_ID/$accountId/g" out/0-templateFromAnalysis.json
    sed -i "s/__ANALYSIS_ID/$analysis/g" out/0-templateFromAnalysis.json
    sed -i "s/__DATASET_ID/$dataset/g" out/0-templateFromAnalysis.json
    sed -i "s/__AWS_REGION/$region/g" out/0-templateFromAnalysis.json
fi

echo "2-Creating public template in account $accountId"
sh 0-create_template_from_analysis.sh $env $accountId

fi
