
# Copyright 2019 Amazon.com, Inc. or its affiliates. All Rights Reserved.
#  SPDX-License-Identifier: MIT-0
env=$1
bucket=$2
portalContactEmail=$3
LOCAL_ENV_NAME=dev

echo "**********************************************"
echo "* IOT Onboarding Depoyment: IOT Sitewise in environemnt '$env' and artifact bucket '$bucket'"
echo "***********************************************"
if [ -z "$env" ] || [ -z "$bucket" ] || [ -z "$portalContactEmail" ]
then
    echo "Environment, Bucket and contact email Must not be Empty"
    echo "Usage:"
    echo "sh deploy.sh <env> <bucket> <portalContactEmail>"
else


echo "0.Checking is portal Already exists"
aws iotsitewise list-portals > portals.json
portalids=$(jq -c '.portalSummaries[] | select(.name=="RigadoMonitorPortaldev") | .id' portals.json)
if [ -z "$portalids" ]
then

ASSET_MODEL_PREFIX=QsTest
ASSET_MODEL_NAME_EPAPERDISPLAY=RigadoEPaperDisplay$ASSET_MODEL_PREFIX$env
ASSET_MODEL_NAME_HOBO_MX=RigadoHoboMX100$ASSET_MODEL_PREFIX$env
ASSET_MODEL_NAME_RS40OCCUPANCY=RigadoRS40_Occupancy_Sensor$ASSET_MODEL_PREFIX$env
ASSET_MODEL_NAME_RUUVITAG=RigadoRuuviTag$ASSET_MODEL_PREFIX$env
ASSET_MODEL_NAMES1_SENSOR=RigadoS1_Sensor$ASSET_MODEL_PREFIX$env
ASSET_MODEL_NAME_ROOT=RigadoSensors$ASSET_MODEL_PREFIX$env

PORTAL_NAME=RigadoMonitorPortal$env
PORTAL_CONTACT_EMAIL=$portalContactEmail
SITEWISE_PROJECT_NAME=RigadoMonitorProject$env

echo "0. getting infra configuration for environement $env"
aws s3 cp s3://$bucket/infra-config-$env.json infra-config.json
cat infra-config.json
iotSitewiseServiceRoleArn=$(jq -r .iotSitewiseRole infra-config.json)

echo "1. Preping local files"
#rm -rf sitewiseResources
mkdir sitewiseResources

##################################################
# I. Asset Models
###########################

echo "0.1 Create Child Asset Models"
#aws iotsitewise create-asset-model --asset-model-name $ASSET_MODEL_NAME_EPAPERDISPLAY --asset-model-properties file://sitewise-model-properties-epaper-display.json > sitewiseResources/sitewiseAssetModelEPaperDisplay.json
aws iotsitewise create-asset-model --asset-model-name $ASSET_MODEL_NAME_HOBO_MX --asset-model-properties file://sitewise-model-properties-hobomx.json > sitewiseResources/sitewiseAssetModelHoboMx.json
aws iotsitewise create-asset-model --asset-model-name $ASSET_MODEL_NAME_RS40OCCUPANCY --asset-model-properties file://sitewise-model-properties-s40-occupancy.json > sitewiseResources/sitewiseAssetModelRs40Occupancy.json
aws iotsitewise create-asset-model --asset-model-name $ASSET_MODEL_NAME_RUUVITAG --asset-model-properties file://sitewise-model-properties-ruuvitag.json > sitewiseResources/sitewiseAssetModelRuuvitag.json
aws iotsitewise create-asset-model --asset-model-name $ASSET_MODEL_NAMES1_SENSOR --asset-model-properties file://sitewise-model-properties-s1-sensor.json > sitewiseResources/sitewiseAssetModelS1Sensor.json

hoboMxModelId=$(jq -r .assetModelId sitewiseResources/sitewiseAssetModelHoboMx.json)
rs40ModelId=$(jq -r .assetModelId sitewiseResources/sitewiseAssetModelRs40Occupancy.json)
ruuvitagModelId=$(jq -r .assetModelId sitewiseResources/sitewiseAssetModelRuuvitag.json)
s1ModelId=$(jq -r .assetModelId sitewiseResources/sitewiseAssetModelS1Sensor.json)

echo "0.1 Create Root Asset Models for models: $hoboMxModelId, $rs40ModelId, $ruuvitagModelId, $s1ModelId"
cp sitewise-root-asset-model.json.tpl sitewiseResources/sitewise-root-asset-model.json
echo "Sitewise resources folder content:"
ls sitewiseResources
if [ $env == $LOCAL_ENV_NAME ]
    then
    sed -i '' "s/HOBOMX100_ASSET_MODEL_ID/$hoboMxModelId/g" sitewiseResources/sitewise-root-asset-model.json
    sed -i '' "s/RS40_ASSET_MODEL_ID/$rs40ModelId/g" sitewiseResources/sitewise-root-asset-model.json
    sed -i '' "s/RUUVI_ASSET_MODEL_ID/$ruuvitagModelId/g" sitewiseResources/sitewise-root-asset-model.json
    sed -i '' "s/S1_ASSET_MODEL_ID/$s1ModelId/g" sitewiseResources/sitewise-root-asset-model.json       
    else
    sed -i "s/HOBOMX100_ASSET_MODEL_ID/$hoboMxModelId/g" sitewiseResources/sitewise-root-asset-model.json
    sed -i "s/RS40_ASSET_MODEL_ID/$rs40ModelId/g" sitewiseResources/sitewise-root-asset-model.json
    sed -i "s/RUUVI_ASSET_MODEL_ID/$ruuvitagModelId/g" sitewiseResources/sitewise-root-asset-model.json
    sed -i "s/S1_ASSET_MODEL_ID/$s1ModelId/g" sitewiseResources/sitewise-root-asset-model.json
fi

#The asset model creation can take a little time and trigger a race condition. We implement 
# exponential retry 
echo "Waiting 5 sec for asset models creation completion"
sleep 5
aws iotsitewise create-asset-model --asset-model-name $ASSET_MODEL_NAME_ROOT --cli-input-json file://sitewiseResources/sitewise-root-asset-model.json
rc=$?
if [ $rc -ne 0 ]; then
      echo "Root asset model Creation Failed: Retry ing after 10 seconds"
      sleep 10
      aws iotsitewise create-asset-model --asset-model-name $ASSET_MODEL_NAME_ROOT --cli-input-json file://sitewiseResources/sitewise-root-asset-model.json
      rc=$?
      if [ $rc -ne 0 ]; then
        echo "Root asset model Creation Failed: Retry ing after 30 seconds"
        sleep 30
        aws iotsitewise create-asset-model --asset-model-name $ASSET_MODEL_NAME_ROOT --cli-input-json file://sitewiseResources/sitewise-root-asset-model.json
      fi 
fi

##################################################
# III. Portal
###########################
echo "2.1 Create Portal"
aws iotsitewise create-portal --portal-name $PORTAL_NAME --portal-contact-email $PORTAL_CONTACT_EMAIL --role-arn $iotSitewiseServiceRoleArn > sitewiseResources/sitewisePortalRes.json
portalUrl=$(jq -r .portalStartUrl sitewiseResources/sitewisePortalRes.json)
portalId=$(jq -r .portalId sitewiseResources/sitewisePortalRes.json)
#echo "2.2 Create Project "
#aws iotsitewise create-project --portal-id $portalId --project-name $SITEWISE_PROJECT_NAME > sitewiseResources/sitewiseProject.json
rc=$?
if [ $rc -ne 0 ]; then
      echo "Sitewise Portal Build Failed! Please contact the AWS Quickstart Support" >&2
      exit $rc
else
echo "****************************************************************"
echo "* SiteWise Portal Build Successed                              *"
echo "****************************************************************"
echo "* -------------------------------------------------------------"
echo "* IAM Role Arn: $iotSitewiseServiceRoleArn      "
echo "* -------------------------------------------------------------"
echo "* Portal Url: $portalUrl      "
echo "* -------------------------------------------------------------"
echo "* Portal ID: $portalId      "
echo "* -------------------------------------------------------------"
echo "*"
echo "* Important:                                                   "
echo "* -------------------------------------------------------------"
echo "* You need to go to your AWS console and create an             *"
echo "* administrator user for your AWS IOT Sitewise Project         *"
echo "****************************************************************"
fi

else    
    echo "Portal RigadoMonitorPortal$env already exit in account. Delete all resources to re-run this script"
fi
fi

