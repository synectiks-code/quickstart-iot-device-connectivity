import cdk = require('@aws-cdk/core');
import codebuild = require('@aws-cdk/aws-codebuild');
import codepipeline = require('@aws-cdk/aws-codepipeline');
import codepipeline_actions = require('@aws-cdk/aws-codepipeline-actions');
import { CfnParameter, StackProps, RemovalPolicy } from "@aws-cdk/core";
import { Bucket } from "@aws-cdk/aws-s3";
import { Role, ServicePrincipal, ManagedPolicy } from "@aws-cdk/aws-iam";

//TODO: this will need to be removed after publication of teh quickstart
var GITHUB_TOKEN_SECRET_ID = "rollagrgithubtoken"

export class IotOnboardingCodePipelinesStack extends cdk.Stack {
  constructor(scope: cdk.Construct, id: string, props?: StackProps) {
    super(scope, id, props);

    const region = (props && props.env) ? props.env.region : ""
    const account = (props && props.env) ? props.env.account : ""

    const envName = this.node.tryGetContext("envName");
    //const gitHubRepo = "aws-quickstart/quickstart-iot-device-connectivity"
    const gitHubRepo = "grollat/quickstart-iot-device-connectivity"

    //CloudFormatiion Input Parmetters to be provided by end user:
    const contactEmail = new CfnParameter(this, "contactEmail", {
      type: "String",
      allowedPattern: "^([a-zA-Z0-9_\\-\\.]+)@([a-zA-Z0-9_\\-\\.]+)\\.([a-zA-Z]{2,5})$",
      description: "A contact email address for the solution administrator"
    });
    const quickSightAdminUserName = new CfnParameter(this, "quickSightAdminUserName", {
      type: "String",
      allowedPattern: ".+",
      description: "The Name of an existing Amin user created for Amazon Quicksihght (see quickstart guide)"
    });
    const sourceTemplateArn = new CfnParameter(this, "sourceTemplateArn", {
      type: "String",
      allowedPattern: ".+",
      description: "The Arn of a the source public template (see quickstart guide)"
    });
    const rootMqttTopic = new CfnParameter(this, "rootMqttTopic", {
      type: "String",
      allowedPattern: ".+",
      default: "data/#",
      description: "the root MQTT topic where onboarded devices publish (see quickstart guide)"
    });

    const artifactBucket = new Bucket(this, "iotOnboardingArtifacts", {
      removalPolicy: RemovalPolicy.DESTROY,
      versioned: true
    })

    //TODO: provide a more granular access to the code build pipeline
    const buildProjectRole = new Role(this, 'buildRole', {
      assumedBy: new ServicePrincipal('codebuild.amazonaws.com'),
      managedPolicies: [ManagedPolicy.fromAwsManagedPolicyName("AdministratorAccess")]
    })

    const infraBuild = new codebuild.PipelineProject(this, 'infraBuilProject', {
      projectName: "code-build-iot-onboarding-infra",
      role: buildProjectRole,
      buildSpec: codebuild.BuildSpec.fromObject({
        version: '0.2',
        phases: {
          install: {
            "runtime-versions": {
              nodejs: 10
            },
            commands: [
              'echo "CodeBuild is running in $AWS_REGION" && aws configure set region $AWS_REGION',
              'npm install -g aws-cdk',
              'cdk --version',
              'cd iot-onboarding-infra',
              'npm install'
            ]
          },
          build: {
            commands: [
              'echo "Build and Deploy Infrastructure"',
              'pwd && sh deploy.sh ' + envName + " " + artifactBucket.bucketName + " " + rootMqttTopic.valueAsString
            ],
          },
        },
        artifacts: {
          "discard-path": "yes",
          files: [
            'iot-onboarding-infra/infra-config-' + envName + '.json',
          ],
        },
      }),
      environment: {
        buildImage: codebuild.LinuxBuildImage.AMAZON_LINUX_2_3,
      },
    });

    const lambdaBuild = new codebuild.PipelineProject(this, 'lambdaBuilProject', {
      projectName: "code-build-iot-onboarding-lambda-" + envName,
      role: buildProjectRole,
      buildSpec: codebuild.BuildSpec.fromObject({
        version: '0.2',
        phases: {
          install: {
            "runtime-versions": {
              golang: 1.13
            }
          },
          build: {
            commands: [
              'echo "Build and Deploy lambda Function"',
              'cd iot-onboarding-service',
              'pwd && sh lbuild.sh ' + envName + " " + artifactBucket.bucketName
            ],
          },
        }
      }),
      environment: {
        buildImage: codebuild.LinuxBuildImage.AMAZON_LINUX_2_3,
      },
    });

    const glueEtlBuild = new codebuild.PipelineProject(this, 'glueETLBuilProject', {
      projectName: "code-build-iot-onboarding-etl-" + envName,
      role: buildProjectRole,
      buildSpec: codebuild.BuildSpec.fromObject({
        version: '0.2',
        phases: {
          build: {
            commands: [
              'echo "Uploading ETK script to s3"',
              'cd iot-onboarding-data-processing',
              'pwd && sh ./deploy.sh ' + envName + " " + artifactBucket.bucketName
            ],
          },
        }
      }),
      environment: {
        buildImage: codebuild.LinuxBuildImage.AMAZON_LINUX_2_3,
      },
    });

    const siteWiseBuild = new codebuild.PipelineProject(this, 'siteWiseBuildProject', {
      projectName: "code-build-iot-onboarding-sitewise-" + envName,
      role: buildProjectRole,
      buildSpec: codebuild.BuildSpec.fromObject({
        version: '0.2',
        phases: {
          build: {
            commands: [
              'echo "Building sitewise Assets model and project"',
              'cd iot-onboarding-sitewise',
              'pwd && sh ./sitewise.sh ' + envName + " " + artifactBucket.bucketName + " " + contactEmail.valueAsString
            ],
          },
        }
      }),
      environment: {
        buildImage: codebuild.LinuxBuildImage.AMAZON_LINUX_2_3,
      },
    });

    const quicksightBuild = new codebuild.PipelineProject(this, 'quicksightBuildProject', {
      projectName: "code-build-iot-onboarding-quicksight-" + envName,
      role: buildProjectRole,
      buildSpec: codebuild.BuildSpec.fromObject({
        version: '0.2',
        phases: {
          build: {
            commands: [
              'echo "Building Quicksight Dashboard"',
              'cd iot-onboarding-quicksight',
              'pwd && sh ./create-dashboard.sh ' + envName + " " + artifactBucket.bucketName + " " + quickSightAdminUserName.valueAsString + " " + sourceTemplateArn.valueAsString
            ],
          },
        }
      }),
      environment: {
        buildImage: codebuild.LinuxBuildImage.AMAZON_LINUX_2_3,
      },
    });

    const onboardingTest = new codebuild.PipelineProject(this, 'testProject', {
      projectName: "code-build-iot-onboarding-test-" + envName,
      role: buildProjectRole,
      buildSpec: codebuild.BuildSpec.fromObject({
        version: '0.2',
        phases: {
          install: {
            "runtime-versions": {
              nodejs: 10
            },
            commands: [
              "yum -y install epel-release",
              "yum -y install mosquitto",
              "npm install -g newman"
            ]
          },
          build: {
            commands: [
              'echo "Testing Deployed on boarding service"',
              'cd e2e',
              'pwd && sh ./test.sh ' + envName + " " + artifactBucket.bucketName
            ],
          },
        }
      }),
      environment: {
        buildImage: codebuild.LinuxBuildImage.AMAZON_LINUX_2_3,
      },
    });



    const sourceOutput = new codepipeline.Artifact();
    const cdkBuildOutputLambda = new codepipeline.Artifact('CdkBuildOutputLambda');
    const cdkBuildOutputETL = new codepipeline.Artifact('CdkBuildOutputETL');
    const cdkBuildOutputInfra = new codepipeline.Artifact('CdkBuildOutputInfra');
    const cdkBuildOutputTest = new codepipeline.Artifact('CdkBuildOutputTest');
    const siteWiseOutput = new codepipeline.Artifact('siteWiseOutput');
    const quickSightOutput = new codepipeline.Artifact('quickSightOutput');
    new codepipeline.Pipeline(this, 'IotOnboardingPipeline', {
      pipelineName: "code-pipeline-iot-onboarding-" + envName,
      stages: [
        {
          stageName: 'Source',
          actions: [
            new codepipeline_actions.GitHubSourceAction({
              actionName: 'GitHub_Source',
              owner: 'grollat',
              repo: gitHubRepo,
              //TODO: this will need to be removed after publication of teh quickstart
              oauthToken: cdk.SecretValue.secretsManager(GITHUB_TOKEN_SECRET_ID),
              output: sourceOutput,
            }),
          ],
        },
        {
          stageName: 'Build',
          actions: [
            new codepipeline_actions.CodeBuildAction({
              actionName: 'uploadELTScript',
              project: glueEtlBuild,
              input: sourceOutput,
              runOrder: 1,
              outputs: [cdkBuildOutputETL],
            }),
            new codepipeline_actions.CodeBuildAction({
              actionName: 'buildLambdaCode',
              project: lambdaBuild,
              input: sourceOutput,
              runOrder: 1,
              outputs: [cdkBuildOutputLambda],
            }),
            new codepipeline_actions.CodeBuildAction({
              actionName: 'deployInfra',
              project: infraBuild,
              input: sourceOutput,
              runOrder: 2,
              outputs: [cdkBuildOutputInfra],
            }),
          ],
        },
        {
          stageName: 'Test',
          actions: [
            new codepipeline_actions.CodeBuildAction({
              actionName: 'testOnboardingService',
              project: onboardingTest,
              input: sourceOutput,
              outputs: [cdkBuildOutputTest],
            }),
          ],
        },
        {
          stageName: 'Deploy',
          actions: [
            new codepipeline_actions.S3DeployAction({
              actionName: "deployInfraConfigToS3",
              bucket: artifactBucket,
              runOrder: 1,
              input: cdkBuildOutputInfra
            }),
            new codepipeline_actions.CodeBuildAction({
              actionName: 'setupSitewise',
              project: siteWiseBuild,
              input: sourceOutput,
              runOrder: 2,
              outputs: [siteWiseOutput],
            }),
            new codepipeline_actions.CodeBuildAction({
              actionName: 'setupQuicksight',
              project: quicksightBuild,
              input: sourceOutput,
              runOrder: 2,
              outputs: [quickSightOutput],
            })
          ],
        },
      ],
    });

  }

}








