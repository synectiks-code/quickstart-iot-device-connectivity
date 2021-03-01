#!/usr/bin/env node
import 'source-map-support/register';
import cdk = require('@aws-cdk/core');
import { IotOnboardingCodePipelinesStack } from '../lib/iot-onboarding-code-pipelines-stack';


const app = new cdk.App();
const envName = app.node.tryGetContext("envName");

new IotOnboardingCodePipelinesStack(app, 'IotOnboardingCodePipelinesStack' + envName);
