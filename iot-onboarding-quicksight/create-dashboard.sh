#!/bin/sh

env=$1
bucket=$2
adminUserName=$3
sourceTemplateArn=$4
qsUserRegion=$5

#update this variable to specify the name of your loval env
LOCAL_ENV_NAME=local_dev_env

echo "**********************************************"
echo "* IOT Onboarding: Create Quicksight Dashboard from template for env '$env' "
echo "environement:     $env"
echo "bucket:           $bucket"
echo "admin:            $adminUserName"
echo "source template:  $sourceTemplateArn"
echo "***********************************************"
if [ -z "$env" ] || [ -z "$bucket" ] || [ -z "$adminUserName" ] || [ -z "$sourceTemplateArn" ]
then
    echo "Environment, Bucket, Admin user and source template Must not be Empty"
    echo "Usage:"
    echo "sh quicksight.sh <env> <bucket> <adminUserArn> <sourceTemplateArn>"
else
#Getting local accountID and region
echo "0.1-Getting account ID,  region annd Infrastructure configuration from artifact bucket"
rm -rf out
mkdir out
aws sts  get-caller-identity > out/id.json
accountId=$(jq -r .Account out/id.json)
sourceTemplateArn=$(echo $sourceTemplateArn | sed 's/\//\\\//g')
adminUserArn=$(echo "arn:aws:quicksight:us-east-1:$accountId:user/default/$adminUserName" | sed 's/\//\\\//g')
aws s3 cp s3://$bucket/infra-config-$env.json infra-config.json
cat infra-config.json
region=$(jq -r .region infra-config.json)
glueDbName=$(jq -r .glueDbName infra-config.json)
athenaTableName=$(jq -r .athenaTableName infra-config.json)
adminUserArn=$(echo "arn:aws:quicksight:$qsUserRegion:$accountId:user/default/$adminUserName" | sed 's/\//\\\//g')

echo "Current region: $region"
echo "Account: $accountId"
echo "Admin user: $adminUserArn"
echo "Source Template: $sourceTemplateArn"
echo "Glue DB name: $glueDbName"
echo "Athena Table: $athenaTableName"

#Copping teamplate of CLI json intput files to be able to substitute values
echo "0.2 Copying CLI request templates for subtitutions"
cp 1-datasource.json.tpl out/1-datasource.json
cp 2-dataset.json.tpl out/2-dataset.json
cp 3-DashboardFromTemplate.json.tpl out/3-DashboardFromTemplate.json

#Creating Quicksight Datasource
echo "1-Create Datasource $accountId"
RANDOM=$$
datasourceId="RigadoDatasource$env$RANDOM"

if [ $env == $LOCAL_ENV_NAME ]
    then
    sed -i '' "s/__ACCOUNT_ID/$accountId/g" out/1-datasource.json
    sed -i '' "s/__AWS_REGION/$region/g" out/1-datasource.json
    sed -i '' "s/__DATA_SOURCE_ID/$datasourceId/g" out/1-datasource.json
    sed -i '' "s/__ADMIN_USER_ARN/$adminUserArn/g" out/1-datasource.json
    else
    sed -i  "s/__ACCOUNT_ID/$accountId/g" out/1-datasource.json
    sed -i  "s/__AWS_REGION/$region/g" out/1-datasource.json
    sed -i  "s/__DATA_SOURCE_ID/$datasourceId/g" out/1-datasource.json
    sed -i  "s/__ADMIN_USER_ARN/$adminUserArn/g" out/1-datasource.json
fi

echo "Final Datasource template: "
cat out/1-datasource.json
sh 1-create-datasource.sh $env $accountId $datasourceId
    rc=$?
    if [ $rc -ne 0 ]; then
      echo "An Error Occured during datasoutrce creation. Existing with status $rc" >&2
      exit $rc
    fi

#Creating Quicksight Dataset
echo "2-Create Dataset $accountId"
datasourceid=$(jq -r .DataSourceId out/datasource.json)
datasetId="RigadoDataset$env$RANDOM"

if [ $env == $LOCAL_ENV_NAME ]
    then
    sed -i '' "s/__DATA_SET_ID/$datasetId/g" out/2-dataset.json
    sed -i '' "s/__ACCOUNT_ID/$accountId/g" out/2-dataset.json
    sed -i '' "s/__DATASOURCE_ID/$datasourceid/g" out/2-dataset.json
    sed -i '' "s/__GLUE_DB/$glueDbName/g" out/2-dataset.json
    sed -i '' "s/__ATHENA_TABLE_NAME/$athenaTableName/g" out/2-dataset.json
    sed -i '' "s/__AWS_REGION/$region/g" out/2-dataset.json
    sed -i '' "s/__ADMIN_USER_ARN/$adminUserArn/g" out/2-dataset.json
    else
    sed -i  "s/__DATA_SET_ID/$datasetId/g" out/2-dataset.json
    sed -i  "s/__ACCOUNT_ID/$accountId/g" out/2-dataset.json
    sed -i  "s/__DATASOURCE_ID/$datasourceid/g" out/2-dataset.json
    sed -i  "s/__GLUE_DB/$glueDbName/g" out/2-dataset.json
    sed -i  "s/__ATHENA_TABLE_NAME/$athenaTableName/g" out/2-dataset.json
    sed -i  "s/__AWS_REGION/$region/g" out/2-dataset.json
    sed -i  "s/__ADMIN_USER_ARN/$adminUserArn/g" out/2-dataset.json
fi

echo "Final Dataset template: "
cat out/2-dataset.json
sh 2-create-dataset.sh $env $accountId $datasetId
    rc=$?
    if [ $rc -ne 0 ]; then
      echo "An Error Occured during dataset creation. Existing with status $rc" >&2
      exit $rc
    fi

#Creating Quicksight Dashboard from template
echo "2-Create Dashboard $accountId"
dhasboardIds="RigadoDashboard$env$RANDOM"


if [ $env == $LOCAL_ENV_NAME ]
    then
  sed -i '' "s/__DATA_SET_ID/$datasetId/g" out/3-DashboardFromTemplate.json
  sed -i '' "s/__ACCOUNT_ID/$accountId/g" out/3-DashboardFromTemplate.json
  sed -i '' "s/__SOURCE_TEMPLATE_ARN/$sourceTemplateArn/g" out/3-DashboardFromTemplate.json
  sed -i '' "s/__AWS_REGION/$region/g" out/3-DashboardFromTemplate.json
  sed -i '' "s/__ADMIN_USER_ARN/$adminUserArn/g" out/3-DashboardFromTemplate.json
  sed -i '' "s/__DASHBOARD_ID/$dhasboardIds/g" out/3-DashboardFromTemplate.json
  sed -i '' "s/__ENV/$env/g" out/3-DashboardFromTemplate.json
    else
  sed -i  "s/__DATA_SET_ID/$datasetId/g" out/3-DashboardFromTemplate.json
  sed -i  "s/__ACCOUNT_ID/$accountId/g" out/3-DashboardFromTemplate.json
  sed -i  "s/__SOURCE_TEMPLATE_ARN/$sourceTemplateArn/g" out/3-DashboardFromTemplate.json
  sed -i  "s/__AWS_REGION/$region/g" out/3-DashboardFromTemplate.json
  sed -i  "s/__ADMIN_USER_ARN/$adminUserArn/g" out/3-DashboardFromTemplate.json
  sed -i  "s/__DASHBOARD_ID/$dhasboardIds/g" out/3-DashboardFromTemplate.json
  sed -i  "s/__ENV/$env/g" out/3-DashboardFromTemplate.json
fi

echo "Final Dashboard template: "
cat out/3-DashboardFromTemplate.json
sh 3-create_dashboard_from_template.sh $env $accountId
    rc=$?
    if [ $rc -ne 0 ]; then
      echo "An Error Occured during dashbord creation. Existing with status $rc" >&2
      exit $rc
    fi

echo "Successfully created dashboard $dhasboardIds for user $adminUserArn"
fi
