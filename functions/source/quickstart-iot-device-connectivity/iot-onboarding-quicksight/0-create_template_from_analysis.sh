env=$1
account=$2

aws quicksight delete-template --aws-account-id $account --template-id iotOnboardingRigadoQuicksightPublicTemplate$env
aws quicksight create-template --aws-account-id $account --template-id iotOnboardingRigadoQuicksightPublicTemplate$env --cli-input-json file://out/0-templateFromAnalysis.json
aws quicksight describe-template  --aws-account-id $account --template-id iotOnboardingRigadoQuicksightPublicTemplate$env 
aws quicksight update-template-permissions --aws-account-id $account --template-id iotOnboardingRigadoQuicksightPublicTemplate$env --grant-permissions file://out/0-templatePermission.json --profile default
