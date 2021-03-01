
# Copyright 2019 Amazon.com, Inc. or its affiliates. All Rights Reserved.
#  SPDX-License-Identifier: MIT-0
env=$1
bucket=$2
stackName=$3

ASSET_MODEL_PREFIX=QsTest
ASSET_MODEL_NAME_EPAPERDISPLAY=RigadoEPaperDisplay$ASSET_MODEL_PREFIX$env
ASSET_MODEL_NAME_HOBO_MX=RigadoHoboMX100$ASSET_MODEL_PREFIX$env
ASSET_MODEL_NAME_RS40OCCUPANCY=RigadoRS40_Occupancy_Sensor$ASSET_MODEL_PREFIX$env
ASSET_MODEL_NAME_RUUVITAG=RigadoRuuviTag$ASSET_MODEL_PREFIX$env
ASSET_MODEL_NAMES1_SENSOR=RigadoS1_Sensor$ASSET_MODEL_PREFIX$env
ASSET_MODEL_NAME_ROOT=RigadoSensors$ASSET_MODEL_PREFIX$env

echo "**********************************************"
echo "* Rigado Sensor Quickstart Depoyment: Cleasing Sitewise asset models in environemnt '$env' and artifact bucket '$bucket'"
echo "***********************************************"
if [ -z "$env" ] || [ -z "$bucket" ]
then
    echo "Environment Must not be Empty"
    echo "Usage:"
    echo "sh deploy.sh <env> <bucket> [<stackName>]"
else
echo "1.0 Deleting asset models"
mkdir sitewiseResources
aws iotsitewise list-asset-models > sitewiseResources/models.json
modelId=$(jq -c '.assetModelSummaries[] | select(.name=="'$ASSET_MODEL_NAME_ROOT'") | .id' sitewiseResources/models.json | sed 's/"/ /g')
echo "deleting $modelId from sitewise"
aws iotsitewise delete-asset-model --asset-model-id $modelId
modelId=$(jq -c '.assetModelSummaries[] | select(.name=="'$ASSET_MODEL_NAMES1_SENSOR'") | .id' sitewiseResources/models.json | sed 's/"/ /g')
echo "deleting $modelId from sitewise"
aws iotsitewise delete-asset-model --asset-model-id $modelId
modelId=$(jq -c '.assetModelSummaries[] | select(.name=="'$ASSET_MODEL_NAME_RUUVITAG'") | .id' sitewiseResources/models.json | sed 's/"/ /g')
echo "deleting $modelId from sitewise"
aws iotsitewise delete-asset-model --asset-model-id $modelId
modelId=$(jq -c '.assetModelSummaries[] | select(.name=="'$ASSET_MODEL_NAME_RS40OCCUPANCY'") | .id' sitewiseResources/models.json | sed 's/"/ /g')
echo "deleting $modelId from sitewise"
aws iotsitewise delete-asset-model --asset-model-id $modelId
modelId=$(jq -c '.assetModelSummaries[] | select(.name=="'$ASSET_MODEL_NAME_HOBO_MX'") | .id' sitewiseResources/models.json | sed 's/"/ /g')
echo "deleting $modelId from sitewise"
aws iotsitewise delete-asset-model --asset-model-id $modelId
modelId=$(jq -c '.assetModelSummaries[] | select(.name=="'$ASSET_MODEL_NAME_EPAPERDISPLAY'") | .id' sitewiseResources/models.json | sed 's/"/ /g')
echo "deleting $modelId from sitewise"
aws iotsitewise delete-asset-model --asset-model-id $modelId

fi
