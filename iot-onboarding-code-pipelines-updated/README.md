# Welcome to your CDK TypeScript project

This is a blank project for CDK development with TypeScript.

The `cdk.json` file tells the CDK Toolkit how to execute your app.

## Useful commands

* `npm run build`   compile typescript to js
* `npm run watch`   watch for changes and compile
* `npm run test`    perform the jest unit tests
* `cdk deploy`      deploy this stack to your default AWS account/region
* `cdk diff`        compare deployed stack with current state
* `cdk synth`       emits the synthesized CloudFormation template

cdk --app bin/iot-onboarding-code-pipelines-updated.js synth
cdk --app bin/iot-onboarding-code-pipelines-updated.js  bootstrap
cdk --app bin/iot-onboarding-code-pipelines-updated.js deploy
cdk --app bin/iot-onboarding-code-pipelines-updated.js diff
cdk --app bin/iot-onboarding-code-pipelines-updated.js deploy --parameters contactEmail=papu.bhattacharya@synectiks.com \
--parameters quickSightAdminUserName=admin --parameters quickSightAdminUserRegion=us-east-1 \
--parameters sourceTemplateArn=arn:aws:quicksight:us-east-1:657907747545:templateiotOnboardingRigadoQuicksightPublicTemplatedev \
--parameters rootMqttTopic=appkube-iot-mqtt --parameters environment=dev --parameters gitHubUserName=synectiks-code \
--parameters githubtoken=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

