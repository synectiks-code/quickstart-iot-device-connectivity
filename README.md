<p align="center">
  <h2 align="center">AWS quickstart-iot-device-connectivity</h2>

  <p align="center">
    An AWS landing zone for IoT device connectivity in partnership with AWS IoT Partner <a href="https://www.rigado.com/market-solutions/smart-hospitality-retail-solutions-powered-by-aws-iot/?did=pa_card&trk=pa_card">Rigado</a>
  </p>
</p>

<!-- TABLE OF CONTENT -->
<details open="open">
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
  </ol>
</details>

## TODO BEFORE LAUNCH
* remove git credentials
* change git repo/owner to aws quickstart and remove dev branch
* change qQS template URL
* Change default public template ARN
* add flic button.
* add case issue in FAQ when IOT core rule fail.
* add description of new parameters (env and quicksight user region) : aws quicksight list-users  --aws-account-id 045495081976 --namespace default
* comments on the duplicate fields issue
<!-- ABOUT THE PROJECT -->
## About The Project

This AWS quickstart aims at helping AWS IoT customers to quickly get started with an IoT landing zone on AWS including:
* A REST microservice to onboard devices and gateways by serial number. The service creates the AWS IoT Core resources to securely connect to AWS MQTT Broker.
* An IoT Datalake ingesting the data from the long term storage and analytics
* An example AWS Quicksight Dahsboard to display data form the datalake (Compatible devices only)
* An example IoT Device real-time Monitoring dahsboard using AWS IoT SiteWise (Compatible Devices Only)

![Alt text](images/iot-quickstart-archtecture.png?raw=true "Title")

Typical use cases can be Smart Kitchen or Smart Retail Store.

The QuickStart is being released in partnership with AWS IoT and Travel and Hospitality Competency Partner [Rigado](https://www.rigado.com/market-solutions/smart-hospitality-retail-solutions-powered-by-aws-iot/?did=pa_card&trk=pa_card) and compatible with [Rigado Alegro Kit](https://www.rigado.com/market-solutions/smart-hospitality-retail-solutions-powered-by-aws-iot/?did=pa_card&trk=pa_card)

### About this repository

This repository includes the following folders.

#### e2e
This folder contains end-to-end tests for the onboarding microservice and MQTT connectivity tests that validate that onboarded devices can connect to the AWS IoT Core MQTT Broker. It uses newman, a CLI tool allowing to run Postman tests and mosquitto as an MQTT Client.

#### iot-onboarding-code-pipelines
This folder contains an AWS CDK project that builds an AWS Code Pipeline project which is used to deploy the architecture decribed above. We use this method to provide consistent build experience for our CDK project independently from builders environement (NodeJS version...). The pipeline has the following steps.

![Alt text](images/quickstart-cicd.png?raw=true "Title")

#### iot-onboarding-data-processing
This folder contains a Python ETL (Extract Load Transform) script that flatten the device Json messages to be queried by Amazon Athena and Amazon QuickSight. This ETL script is run by an AWS Glue Job.

#### iot-onboarding-infra
This folder contains a CDK project that builds most of the infrastructure components described above except the AWS Quicksight and AWS IoT SiteWise Dahsboards which are not yet supported by AWS CloudFormation at the time of construction of this QuickStart.

#### iot-onboarding-quicksight
This folder contains a linux shell script that automates the creation of an Amazon QuickSight Dahsboard based on a public template. Note that this requires for the target account to have an [activated Amazon QuickSight Enterprise account](https://docs.aws.amazon.com/quicksight/latest/user/signing-up.html). Also, the example dashboard assumes the following structure for MQTT messages from the device (based on AWS Partner Rigado). 

```
Topic: data/#
Body:
{
    measurements: {
        <mesurement name>: <data>,
        ...
    },
    device: {
        deviceId: <serial number>
        gatewayid: <rigado gateway ID>,
        capabilityModelId: <urn:vendor:model:*>",
    },
}
```

An example message from a rigado Device:

```json
{
  "device": {
    "gatewayId": "C0300A1930-00366",
    "deviceId": "ac233fa2129e",
    "capabilityModelId": "urn:rigado:S1_Sensor:2",
    "dtmi": "dtmi:rigado:S1Sensor;1"
  },
  "measurements": {
    "batteryLevel": 100,
    "temperature": 19.02734375,
    "humidity": 47.7890625,
    "rssi": -63
  }
}
```

The base MQTT topic can be configured as an input parametter from the CICD pipeline CloudFormation Stack and the IoT Datalake uses glue crawlers to dynamically identify the data structure of the incoming MQTT messages. This means that QuickStart users who use different device configurations can quickly adapt the dashboard to their specific need.

**Note:** Using AWS CLI comes with limitations compared to AWS CloudFormation and some resources (such as dashboard, dataset and datasource) may need to be manually deleted to be updated or in case of failure during deployement. We hope that providing this autromated dahsboard allows you to move faster by relying on an example and may move this to a more robust infrastructure as code solution when available.
**Note2:** The script will only fully succeed for Amazon QuickSight Enterprise users. QuickSights users with a standard subscription will need to manually create the dashboard on top of the created Dataset.

The example QuickSight Dashboard looks as follows:

![Alt text](images/quicksight.png?raw=true "Title")


#### iot-onboarding-service
This folder contains the Golang Code of the onboarding service Lambda function. The function sits behind an AWS API gateway REST API exposing the following Services
```
POST {{baseUrl}}api/onboard/{{deviceName}}
GET {{baseUrl}}api/onboard/{{deviceName}}
DELET {{baseUrl}}api/onboard/{{deviceName}}
```
These endpoints respetively create, retreive and delete a device or gateway including the following AWS IoT resources:
* A Device Certificate
* An IoT Thing and Associated policy to publish on the base topic provided as parameter to the quickstart CICD CloudFormation template.

The service Creates and Retreives enpoints return all the data needed to setup the device for AWS connectivity and the message structure is as follow:
```json
{
    "serialNumber": "<device serial number>",
    "deviceName": "<device name = device serial number>",
    "thingId": "<ID of the AWS IoT Core Thing>",
    "credential": {
        "certificateArn": "<ARN of the certificate created for the IoT Thing>",
        "certificateId": "<ID of the certificate created for the IoT Thing>",
        "certificatePem": "<PEM Certificateg>",
        "privateKey": "<Private Key>",
        "publicKey": "<Public Key>"
    },
    "mqttEndpoint": "<MQTT enpoint of the IoT COre project>",
    "error": {
        "code": "<error code>",
        "msg": "<error message>",
        "type": "<error type>"
    }
}
```

The service is secured by Amazon Cognito and a random user is created during infrastructure deployment along with a refresh token. To access the service, the quickStart owner needs to access credentials stored in a S3 bucket following the stack successful build and generate temporary credentials in the form of a Cognito Access token. These temporary credentials can then be used to access the device configuration data. More information on this flow is provided below.

**Note:** As part of the partnership with Rigado on this QuickStart, the Rigado team created a web-based Wizard for their Alegro Kit users that takes care of generating the temporary credentials and setting up the devices remotely. More information at [Rigado.com](add rigado kit url)

#### iot-onboarding-sitewise
This folder contains a linux shell script that builds AWS IoT SiteWise resources needed to build a real time dashboard. These resources include:
* A device model hierarchy, composed of a root device and 4 child devices (based on the Rigado Alegro Kit content)
* A AWS IoT Sitewise project and portal.

**Note 1:** A few manual steps are required to create the assets and add them to a dashoard. The following result can be obtain in just a few minutes with the Rigado Allegro Kit.
**Note 2:** See the AWS IoT SiteWise documentation in order to follow required steps prior to deployement (Such as creatinng an AWS SSO user)
**Note 3:** Contrary to the datalake part, the IoT Core broker rule that ingests the data into AWS IoT SiteWise is not model-agnostic. This means that non-Rigado-kit users need to update both the CDK script in the __iot-onboarding-infra__ folder and the sitewise shell script to acommodate for their device specificity. We hope that the code we provided here is a good example allowing these users to quickly build their real-time pipeline and may add additional out-of-the-box support for other IoT partners in the future.

![Alt text](images/sitewise.png?raw=true "Title")

### Built With

This project use the folowing tools and frameworks:
* [Golang](https://golang.org/)
* [Python](https://www.python.org/)
* [AWS CDK](https://aws.amazon.com/cdk/)
* [newman](https://learning.postman.com/docs/running-collections/using-newman-cli/command-line-integration-with-newman/)
* [jq](https://stedolan.github.io/jq/)
* [mosquitto](https://mosquitto.org/)


<!-- GETTING STARTED -->
## Getting Started

To get started with this AWS Quickstart, follow the steps below (make sure to follow the prerequisits section first)

### Prerequisites

#### Service Quotas

You need to ensure the following Quota requirements are met in your account. If not the case, create a request a request from the [AWS Service Quotas Dashboard](https://console.aws.amazon.com/servicequotas/home/)

Input         |    Quota Name |         Value | Comments
------------- | ------------- | ------------- | -------------
CodeBuild | Maximum number of concurrent running builds*	| 2 | The QuickStart uses a pipeline with parallel build steps. You need at this 2 concurent builds allowed. see [AWS CodeBuild Quotas](https://docs.aws.amazon.com/codebuild/latest/userguide/limits.html)

#### AWS SSO activation (Optional if you don'y want to use AWS IoT SiteWise)
AWS SSO provides identity federation for AWS IoT SiteWise Monitor so that you can control access to your portals. With AWS SSO, users sign-in with their corporate email and password instead of an AWS account. Follow the steps under **Enabling AWS SSO** in the [AWS IoT sitewise documentationn](https://docs.aws.amazon.com/iot-sitewise/latest/userguide/monitor-getting-started.html)

#### Create a Quicksight account (Optional if you don't want to use AWS QuickSigth or are already signed up)
If you haven't already, sign-up for Amazon QuickSight using the steps in the [AWS documentation](https://docs.aws.amazon.com/quicksight/latest/user/signing-up.html). If you plan to deploy the default dashboard, you need an Amazon QuickSight Enterprise account.

#### Validate your email adrdess with SES
This QuickStart uses the email address provided in input to the AWS CloudFormation template as both sender and receiver of email notifications. These notifications will provide you with the key credentials to use the device onboarding Micro-services. More specifically, for users of the [Rigado Alegro Kit](https://www.rigado.com/market-solutions/smart-hospitality-retail-solutions-powered-by-aws-iot/?did=pa_card&trk=pa_card), the email will provide the data necessary to use the Rigado Wizard to automatically onboard the Rigado Gateway. In order to be able to use this email address, SES requires you to verify the provided email as described in the [Amazon SES Documentation](https://docs.aws.amazon.com/ses/latest/DeveloperGuide/verify-email-addresses.html). 

#### Enable Logginng for AWS IoT Core
Enablling logging for AWS IoT Core will facilitate troubleshooting of device connectivity. This is especially useful is you are not using a Rigado device. See instruction in the [AWS IoT Core documentation](https://docs.aws.amazon.com/iot/latest/developerguide/configure-logging.html)


### Installation

1. To get started with the deployment download the [AWS CloudFormation template](iot-onboarding-code-pipelines/iot-onboarding-int.yml). Note that you can also fork this repository and use it as a base for your own IoT project.

2. Go to the AWS cloudFormation console and launch the stack. The following parametters are required inputs
![Alt text](images/cloudformation-form.png?raw=true "Title")

Input  | Description
------------- | -------------
contactEmail  | Email of an administrator ussed for the AWS IoT SiteWise portal creation. (see [AWS documentation](https://docs.aws.amazon.com/iot-sitewise/latest/userguide/administer-portals.html#portal-change-admins))
environment  | You can leave this parametter default value to 'int'. This parametter can be used to create multiple stacks in the same AWS region and the same AWS account (for example, one for development and one for production)
quickSightAdminUserName  | The username of an Amazon QuickSight user with an **ADMIN** role. You can list all Amazon QuickSight users by going to the [Amazon QuickSight administration Screen](https://us-east-1.quicksight.aws.amazon.com/sn/admin). If Omitted, the CICD pipeline will not include the QuickSight dahsboard
quickSightAdminUserRegion  | The region in which the above quickSIght user was created. It is usually the regionn in which the use subscribed to Amazon QuickSight and can be obtain by the command ```aws quicksight list-users  --aws-account-id <account-id> --namespace <namespace or default>```
rootMqttTopic  | The root MQTT topic your devices publishes to. If you are using the Rigado [Alegro Kit](https://www.rigado.com/market-solutions/smart-hospitality-retail-solutions-powered-by-aws-iot/?did=pa_card&trk=pa_card), the default value (data/#) does not need to be changed.
sourceTemplateArn  | use: **arn:aws:quicksight:eu-central-1:660526416360:template/iotOnboardingRigadoQuicksightPublicTemplatedev** This is a static location of a public QuickSight dashboard template that we created for the purpose of this QuickStart. This allows you to get started quickly with a fully functional dashboard. Note that this example dahsboard is created specifically for the users of the [Rigado Alegro Kit](https://www.rigado.com/market-solutions/smart-hospitality-retail-solutions-powered-by-aws-iot/?did=pa_card&trk=pa_card). If you are not using rigado devices you will need to create you own dataset, analysis and dahsboard based on the devices you use. The ETL process and Glue Crawler that ingest the data from the IoT Broker to make them available in Amazon QuickSight are data-model agnostic so you just need to link the created AWS Glue table as a datasource in Amazon QuickSight

Once the template run is successful, go to [AWS Code Pipeline](https://console.aws.amazon.com/codesuite/codepipeline). You should see the pipeline executing as below:
![Alt text](images/quickstart-cicd-3.png?raw=true "Title")

If you click on the pipeline name, you can see the steps of the pipeline running:
![Alt text](images/quickstart-cicd-2.png?raw=true "Title")


## Connecting Devices
### Rigado Devices
If you are using Rigado devices using the [Alegro Kit](https://www.rigado.com/market-solutions/smart-hospitality-retail-solutions-powered-by-aws-iot/?did=pa_card&trk=pa_card), go to the rigado wizard as explained in the email received after the devices were activated. Enter the data as received in the email you should have received from the QuickStart script:
```
AWS IoT Connectivity QuickStart Output Values


--------------------------------------------------------------------------------------------
| Cognito URL | https://iot-onboarding-quickstart-<env>.auth.us-east-1.amazoncognito.com/oauth2/token
--------------------------------------------------------------------------------------------
| API Gateway URL | https://<api_id>.execute-api.us-east-1.amazonaws.com/
--------------------------------------------------------------------------------------------
| Client ID | 228v...t9c3
--------------------------------------------------------------------------------------------
| Refresh Token | eyJjdHkiOiJKV1QiLCJl...slN29FrDNqHWo_0e5U85ow
```
Following the completion of the Rigado Wizard flow, the Rigado gateway will be setup to automatically and securely send traffic to AWS IoT Core using MQTT. To validate that traffic is flowing throught correctly, go to the AWS IoT Core console and subscribe to the MQTT topic provided in input of the CloudFormation Template (by default **data/#** for Rigado). Provided at least one sensor has been turned on, you should see messages flowing through as shown below. 
![Alt text](images/iot-core-mqtt-test.png?raw=true "Title")

Once you validated that the traffic is going throught, see section __Visualizing the data  (Alegro Kit Users only)__ below to visualize the data.

**Note:** For convenience reason, the build script create a new Amazon Cognito user each time it runs. If you are using the pipeline to manage your own CICD process, you should cleanup your Amazon cognito user pool users regularly to avoid unecessary charges. After a few runs, going to your Amazon Cognito console should look as below:
![Alt text](images/cognito-users.png?raw=true "Title")
You only need one active user to be able to access the API gateway to onboard devices

### Non-Rigado Devices
For non-rigado devices (or devices not supported by the allegro kit), you need to manually configure you device in ordre to allow it to communicate with AWS IoT Core. The Lambda function deployed by the quickstart allows you to generate the device certificate and keypair to be used for secured communication between the device and IoT Core. It also creates the IoT Thing and appropriate policy. It returns:
-	The Mqtt Endpoint of the AWS IoT Core broker. 
-	The certificate generated during the service call
-	The keypair (public and private key) generated during the service call
-	Other information about the device

Below and example of exchange with the onboarding service. We first get the session token by using the client ID and the refresh token

```curl --location --request POST 'https://<cognito_domain>.auth.<region>.amazoncognito.com/oauth2/token?grant_type=refresh_token&client_id=<client_id>&refresh_token=<refresh_token> \
--header 'Content-Type: application/x-www-form-urlencoded'
{
    "id_token": "<id_token>",
    "access_token": "<access_token>",
    "expires_in": 3600,
    "token_type": "Bearer"
}
```
Then we can call the onboarding service to create the device. 
```
curl --location --request POST '<api_gateway_url>/api/onboard/<device_serial_number> \
--header 'Authorization: Bearer <access_token>'
{
    "serialNumber": "<device_serial_number> ",
    "deviceName": "<device_serial_number> ",
    "thingId": "<iot_thing_id>",
    "credential": {
        "certificateId": "<certificate_id>",
        "certificatePem": "<certificate data>",
        "privateKey": "<private key data>",
        "publicKey": "<public key data>"
    },
    "mqttEndpoint": "https://data.iot.<region>.amazonaws.com",
    "error": {
        "code": "",
        "msg": "",
        "type": ""
    }
}
```
The response contains the necessary certificates and keys to be used to configure the device. Refer to the [AWS IoT Core documentation](https://docs.aws.amazon.com/iot/latest/developerguide/connect-to-iot.html) for details on configuring IoT device to communicate with AWS.

## Checking Connnectivity
In this section we will validate that the devices has been setup correctly and the traffic is flowing as expected.

## Visualizing the data  (Alegro Kit Users only)
As described in the architecture above, the data sent to the IoT Core MQTT broker is ingested into an IoT Datalake for a posteriori visualization and into AWS IoT SiteWise for live monitoring. The paragraph below describes how to setup these visualizations.

### Visualizing the data using Amazon QuickSight
In order to speed up your IoT project, this QuickStart deploys a predefined dahsboard within Amazon QuickSight. To benefit from this feature, you need:
- To have an Amazon QuickSight Enterprise customer
- To be using the Rigado Allegro Kit as devices (while the Glue ETL ingesting and processing the data is not device agnostic, the Dahsboards does make assumptions on the name of the fields received and will therefore only work out-of-the-box for Rigado Allegro Kit users)
By Default, the data is processed by a scheduled Glue Crawler and ETL ever 24H. this default value is choosen to minimize the cost of running the IoT datalake initially and can be updated easily by changing a CRON expression in the AWS CDK script or directly from the AWS Glue console. The QuickStart also creates a glue trigger that can start the data refresh on demand directly from the AWS Glue Console.

When accessing the Amazon QuickSight Dashboard for the first time, you need to provide access to the Amazon S3 Bucket that contains the refined data (by refined data, we mean the data processed by the ETL script). Instructions to give Amazon QuickSight Access to Amazon S3 bucket is described in the [AWS documentation](https://docs.aws.amazon.com/quicksight/latest/user/troubleshoot-connect-athena.html).
Following the deployment of the QuickStart, you should see an Analysis called **Rigado QuickStart Dahsbord** in your Amazon QuickSight account in the selected regions as shown below.
![Alt text](images/rigado-dahsboard.png?raw=true "Title")

This Dashboard is configured to query 48 hours of data in the past (this is to limit both cost and improve dashboard load time as the quantity of data increases in the future). There are multiple ways you can change this setup while scaling with large amounts of data by using [QuickSight SPICE](https://docs.aws.amazon.com/quicksight/latest/user/spice.html). Note that using SPICE will come with an additional cost.

The Glue ETL that processes the data into a flat structure is also optimized to only query 48 hours of data in the past using [push down predicate](https://docs.aws.amazon.com/glue/latest/dg/aws-glue-programming-etl-partitions.html). this can easily be changed with a minor update to the Python script directly accessible from the AWS Glue Console.

**Note for non-alegro Kit user:** If you are not a Rigado allegro kit user, you will need to create you own Analysis and  datasource targeting the Athena Table for refined data mentioned earlier. This can be done in just a few clicks following the Amazon QuickSight documentation. The Glue job that refines the data is device agnostic as it justs flatten the JSON nested fields. It may, however not lead to practical result for deeply nested data.

### Visualizing the data using AWS SiteWise
The QuickStart creates an AWS IoT Sitewise Assets Model Hierarchy composed 1 root asset model and 4 children assets models. It also creates a Portal. In order to start visualizing the data in the portal, you need to follow the steps below:
1. Go to AWS IoT SiteWise and select Build > Models
2. Choose the Asset model that corresponds to your Rigado Device (if the device you are using does not correspond to any existing asset model, refer to AWS IoT Sitewise documentation to create a dedicated asset model and route the traffic of your device through the apropriate alias using AWS IoT Core)
3. Create an asset under this asset model using the deviceId in the device name
4. Once created, go to "Edit" and enter a property alias for each of the model Measurements. For consistency with the IoT Core Broker rule, the alias value must be as follow:
```
<deviceId><MeasurementNameWithoutDoubleQuotes>
```
See example below for device **ffcfed4dd3ab**
![Alt text](images/sitewise-property-alias-setup.png?raw=true "Title")

Repeat this for all devices sending traffic behing the Rigado Gateway. (Using the Amazon QuickSight Dashboard, you can have a list of all devices sending traffic though the Gateway and use this list too setup love monitoring with AWS IoT SiteWise)

5. Once the asset is created you can access the portal created by the QuickStart or create a portal from scratch following the AWS IoT SiteWise documentation. It will then just take a few minutes to add your assets to dedicated dashboards.

From this point, you can use the created portal to design dashboards for your devices as described in the AWS IoT SiteWise documentation.

**Note for non-alegro Kit user:** If you are not an allegro kit user, you will need to create your own AWS IoT Core Broker rule (following the same model than the one created in the QuickStart) to ingest the properly formated data into AWS IoT SiteWise. You will also need to manually create the Assets Models and Assets following the AWS IoT SiteWise documentation.


## Cleaninng Up
In this quickstart, we use a combination of CLI and CDK for AWS Resources deployement. This is because some services like Amazon QuickSight and AWS IoT Sitewise are not supported by CloudFormation jut yet. Consequently, several manual steps will be required to clean up the deployed resources in the user account. These steps are described below:
1. Empty the Amazon S3 buckets
Identify the buckets created by the stack (they are prefixed by "iotonboardinginfrastack") and ensure you clean the content of these buckets before deleting the stack.
2. Delete the infrastructure CloudFormation stack
Go to CloudFormation and Delete the infrastructuure stack starting with **IoTOnboardingInfraStack**
3. Delete Code Pipeline CloudFormation stack
Go to CloudFormation and Delete code pipeline stack you created
4. Clean-up Amazon QuickSight Dashboard
You can manually delete the resources created in Amazon QuickSight following the Amazon Quicksight Documentation. If you created an Amazon QuickSight Account just for the purpose of this QuickStart you can unsubscribe to the service by following the steps described [here](https://docs.aws.amazon.com/quicksight/latest/user/closing-account.html)
5. Clean-up AWS IoT Sitewise Dashboard
You need to delete the following resources (The deletion procedure is provided in the AWS Documentation),:
- SiteWise Assets.
- Sitewise Assets Models (the quickStart creates 1 root asset model and 4 child asset models).
- Sitewise Projects and Dashboards.


## FAQ

### My Amazon QuickSight dashborad displays errors in the Widgets
Errors displayed in the widgets of the QuickSight Dashboard typically results from one of the following causes:
1. Access to the S3 Bucket created by the CDK script to store the data
**solution:** Follow the Amazon QuickSightt instructions to provide access to the bucket. 
2. No Data is not available yet in the datalake (data either not yet crawled or not yet received)
**solution:** If data has already been received, you can manually trigger the on-demand trigger from the Glue console

### My quicksight deployment script fails with an error
This can happen when the deployment of a previous version has been only partially completed. As a result, some resources already exists and others not yet.
**solution:** Go to QuickSight and delete the existing dashboard and dataset before you retry the pipeline stage.

### My code pipeline action fails with errors about concurent build allowed
Some AWS accounts are by default configured to only allow 1 concurent build using Code Pipeline. A workaround to this issue is to retry the pipeline stage. This will only retry the failed build. A long term solution consists in requesting a Code Build Service limit increase by contacting AWS Support.

### My Sitewise script fails with ResourceAlreadyExistsException
```
An error occurred (ResourceAlreadyExistsException) when calling the CreateAssetModel operation: Another resource is already using the name RigadoHoboMX100QsTestint 
```
**solution:** Delete the Sitewise resources and re-run the pipeline stage.

### I am using the Alegro kit and I don't see any column in my AWS Quicksight dataset
Data takes up to one day to propagate to IoT Datalake due to the schedule chosen for the Crawler for cost optimization reason.
**solution:** Ensure data is flowing using IoT Core Monitor and manually trigger the crawlers and jobs using the ON_DEMAND trigger from the Glue Console

### Deleting the infrastructure stack fails
A common issue when deleteing the stack is when Amazon S3 buckets contain data. In this case, they fail to delete to avoid unintended data loss.
**solution:** S3 Buckets must be manually emptied and deleted prior to retying the stack deletion.

### Deleting AWS IoT SiteWise Portal fails
A common issue while deleting a Sitewise portal is the existence of a project associated to the portal. 
**solution:** You can delete the project from the created portal if you can login as administrator. If you are unable to login as an admin, you can delete the underlying project using the AWS IoT Sitewise CLI.

### I added a new AWS IoT Core rule for AWS IoT Sitewise for a device displayed inn my QuickSight dashboard but I can see traffic in Sitewise
First, you should validate that the device has been setup correctly by subscribing to the MQTT topic or downloading raw data ingested in S3. Once you validate that the traffic is indeed flowing correctly, make sure that the field name used in the rule matches the raw data. It is important to note that the fields name ingested in the datalake are lowercased and should therefore not be used as-is in IoT Core rules.

<!-- LICENSE -->
## License

Distributed under the MIT License. See `LICENSE` for more information.
<!-- CONTACT -->
## Contact

Quickstart Team - [@your_twitter](https://twitter.com/your_username) - email@example.com
Project Link: [https://github.com/your_username/repo_name](https://github.com/your_username/repo_name)

