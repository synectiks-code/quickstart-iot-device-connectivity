# Welcome to your CDK TypeScript project!

This is a blank project for TypeScript development with CDK.

The `cdk.json` file tells the CDK Toolkit how to execute your app.

## Useful commands

 * `npm run build`   compile typescript to js
 * `npm run watch`   watch for changes and compile
 * `npm run test`    perform the jest unit tests
 * `cdk deploy`      deploy this stack to your default AWS account/region
 * `cdk diff`        compare deployed stack with current state
 * `cdk synth`       emits the synthesized CloudFormation template

export email=papu.bhattacharya@synectiks.com
export env=dev
export bucket=iotonboardingcodepipelin-iotonboardingartifacts02-1x3qxi4duzifq
export mqttTopic=appkube-iot-mqtt
./deploy.sh $env $bucket $mqttTopic $email
