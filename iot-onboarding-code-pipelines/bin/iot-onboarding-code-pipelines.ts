#!/usr/bin/env node
import 'source-map-support/register';
import cdk = require('@aws-cdk/core');
import { IotOnboardingCodePipelinesStack } from '../lib/iot-onboarding-code-pipelines-stack';


const app = new cdk.App();
const envName = app.node.tryGetContext("envName");

let stack = new IotOnboardingCodePipelinesStack(app, 'IotOnboardingCodePipelinesStack' + envName, {
    description: "Deploys the IoT Device connectivity pipeline to run the CDK deployment (qs-1rmapn8de)",
});
stack.templateOptions.metadata = { "QuickStartDocumentation": { EntrypointName: "Parameters for launching the deployment pipeline" } }