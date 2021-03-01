
# Copyright 2019 Amazon.com, Inc. or its affiliates. All Rights Reserved.
#  SPDX-License-Identifier: MIT-0
env=$1
bucket=$2
stackName=$3

echo "**********************************************"
echo "* Rigado Sensor Quickstart Depoyment: IOT Sitewise in environemnt '$env' and artifact bucket '$bucket'"
echo "***********************************************"
if [ -z "$env" ] || [ -z "$bucket" ]
then
    echo "Environment Must not be Empty"
    echo "Usage:"
    echo "sh deploy.sh <env> <bucket> [<stackName>]"
else
echo "1.0 Deleting all portals named RigadoMonitorPortal$env"
mkdir sitewiseResources
aws iotsitewise list-portals > sitewiseResources/portals.json
portalids=$(jq -c '.portalSummaries[] | select(.name=="RigadoMonitorPortal'$env'") | .id' sitewiseResources/portals.json)
if [ -z "$portalids" ]
then
    echo "Found no portal named RigadoMonitorPortal$env"
else
echo "Found portals $portalids"
#1.Iterating on prortal
for portalid in $portalids
do
#removing '"' in portal iD
portalid=$(echo $portalid | tr '"' ' ')
aws iotsitewise list-projects --portal-id $portalid > sitewiseResources/projects.json
projectIds=$(jq -c '.projectSummaries[] | .id' sitewiseResources/projects.json)
    #2.Iterating on projects
    for projectid in $projectIds
    do
    #removing '"' in project iD
    projectid=$(echo $projectid | tr '"' ' ')
    echo "Deleting project $projectid"
    aws iotsitewise delete-project --project-id $projectid
    done
echo "Deleting portal $portalid"
aws iotsitewise delete-portal --portal-id $portalid
done
fi

fi
