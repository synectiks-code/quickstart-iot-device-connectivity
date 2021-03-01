<p align="center">
  <h2 align="center">AWS quickstart-iot-device-connectivity</h2>

  <p align="center">
    An AWS landing zone for IOT device connectivity in partnership with Aws Partner <a href="https://www.rigado.com/market-solutions/smart-hospitality-retail-solutions-powered-by-aws-iot/?did=pa_card&trk=pa_card">Rigado</a>
  </p>
</p>



<!-- TABLE OF CONTENTS -->
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
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
    <li><a href="#acknowledgements">Acknowledgements</a></li>
  </ol>
</details>

## TODO BEFORE LAUNCH
* remove git credentials
* change git repo/owner to aws quiickstart and remove dev branch
* change qQS template URL
* add link to the rigado alegro kit
* add info about cognito users for cleanup
* add links to frameworks
* freeze version of CDK and NPM
* Add cost optimization (pushdown predicate, Recrawl Policy, device sampling, quicksight refresh schedule)
* Add FAQ about sitewise dashoard delete error because of exising project)
* Change default public template ARN

<!-- ABOUT THE PROJECT -->
## About The Project

This AWS quickstart aims at hellping AWS IOT customers to quickly get started with an IOT landing zone on AWS including:
* A REST microservice to onboard devices and gateway by serial number. The service creates the AWS IOT Core resources to secuerly connect to AWS MQTT Broker.
* An IOT Datalake ingesting the data from the long term storage and analytics
* An example AWS Quicksight Dahsboard to display data form the datalake (Compatible devices only)
* An example IOT Device real time Monitorinig dahsboard uusing AWS IOT SiteWise (Compatible Devices Only)

![Alt text](images/iot-quickstart-archtecture.png?raw=true "Title")


The QuickStart is being released in partnership with AWS IOT and Travel adn Hospitality Competency Partner Rigado and compatible with [Rigado](https://www.rigado.com/market-solutions/smart-hospitality-retail-solutions-powered-by-aws-iot/?did=pa_card&trk=pa_card) newly launched [Alegro Kit](https://www.rigado.com/market-solutions/smart-hospitality-retail-solutions-powered-by-aws-iot/?did=pa_card&trk=pa_card)

### About this repository

This repository includes the following folder:

#### e2e
This folder contains end-to-end tests for the onboarding microservice and MQTT connectivity tests that validate that onboardded devices can connect to the AWS IOT Core MQTT Broker. It use newman, a CLI tool allowing to run Postman tests and mosquitto as an MQTT Client

#### iot-onboarding-code-pipelines
This folder contains an AWS CDK project that builds a AWS Code Pipeline project which is used to deploy the architecture decribed above. We use this method to be able to provide consistent build experience for our CDK project independently from builders environement (NodeJS version...). The pipeline has the following steps

![Alt text](images/quickstart-cicd.png?raw=true "Title")

#### iot-onboarding-data-processing
This folder contains a Python ETL (Extract Load Transform) script that flatten the device Json messages to be queried by Amazon Athena and Amazon Quicksight. This ETL script is run by a Glue Job.

#### iot-onboarding-infra
This folder contains a CDK project that builds most of the infrastructure components described above except the Quicksight and Sitewise Dahsboards which are not yet supported by AWS CloudFormation at the time of construction of this quickstart.

#### iot-onboarding-quicksight
This folder contains a linux shell script that automates the creation of an AWS QuickSight Dahsboard based on a public template. Note that this requires for the target account to have activated Amazon QuickSight (add link here). Also, the example dashboard assumes the following structure for MQTT messages from the device (based on AWS Partner Rigado). 

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

The base topic can be configured as an input parametter from the CICD pipeline CloudFormation Stack and the IOT Datalake uses glue crawlers to dynamically identify the data structure of the incoming MQTT messages. This means that QuickStart users who use different device configurations can quickly adapy the dashboard to their specific need.

**Note:** Using AWS CLI comes with limitations compared to CloudFormation and some resources (such as dashboard, dataset and datasource) may need to be manually deleted to be updated or in case of failure during deployement. We hope that providing this autromated dahsboard allows you to move faster by relying on an example. We will move this to a more robust infrastructure as code solution when available.
**Note2:** The script will only fully succeed for QuickSight Enterprise users. QuickSights users with a standard subscription will need to manually create the dashboard on top of the created Dataset.

The example dashboard looks as follows:

![Alt text](images/quicksight.png?raw=true "Title")


#### iot-onboarding-service
This folder contains the Golang Code of the onboarding service Lambda function. The function sits behind an AWS API gateway REST API exposing the following Services
```
POST {{baseUrl}}api/onboard/{{deviceName}}
GET {{baseUrl}}api/onboard/{{deviceName}}
DELET {{baseUrl}}api/onboard/{{deviceName}}
```
These endpoints respetively create, retreive and delete a device or gateway including the following AWS IOT resources:
* A Device Certificate
* An IOT Thing and Associated policy to publish on the base topic provided as parameter to the quickstart CICD cloudformatioin template

The service Create and Retreive enpoints return all the data needed to setup the device for AWS connectivity and the message structure is as follow:
```json
{
    "serialNumber": "<device serial number>",
    "deviceName": "<device name = device serial number>",
    "thingId": "<ID of the AWS IOT Core Thing>",
    "credential": {
        "certificateArn": "<ARN of the certificate created for the IOT Thing>",
        "certificateId": "<ID of the certificate created for the IOT Thing>",
        "certificatePem": "<PEM Certificateg>",
        "privateKey": "<Private Key>",
        "publicKey": "<Public Key>"
    },
    "mqttEndpoint": "<MQTT enpoint of the IOT COre project>",
    "error": {
        "code": "<error code>",
        "msg": "<error message>",
        "type": "<error type>"
    }
}
```

The service is secured by Amazon Cognito and a random user is created during infrastructure deployment along with a refresh token. To access the service, the quickStart owner needs to access credentials stored in a S3 bucket following the stack successful build and generate temporary credentials in the form of a Cognito Access token. THese tem;porary credentials can then be used to access the device configuration data. More information on this flow is previded below.

**Note:** As part of the partnership with Rigado on this quickstart, the Rigado team created a web-based Wizard for their Alegro Kit user that takes care of generating the temporary credentials and setting up the devices remotely based on the credentials generated by this Microservice. More information at [Rigado.com](add rigado kit url)

#### iot-onboarding-sitewise
This folder contains a linux shell script that builds AWS IOT SiteWise resourtces needed to build a real time dashboard. These resources include:
* a Device model hierarchy, composed of a root device and 4 child devices (based on the Rigado Alegro Kit content)
* A sitewise project and portal

When working with Rigado devices A few manual steps are required to create the assets and add them to a dashoard. The folloowing result can be obtain in just a few minutes with the Rigado Allegro Kit.

**Note 1:** See the AWS IOT SItewise documentation in order to follow required step prior to deployement (Such as creatinng an AWS SSO user)
**Note 2:** Contrary to the datalake part, the IOT Core broker rule that ingests the data into Sitewise is not model-agnostic. This means that non-Rigado-kit-users need to update both the CDK script in the __iot-onboarding-infra__ folder and the sitewise shell script to acomodate for their device specificity. We hope that the code we provided here is a solid example alowing these user to quickly build their real-time pipeline and may add additional out-of-the-box support for other IOT partners in the future.

![Alt text](images/sitewise.png?raw=true "Title")




### Built With

This project use the folowing tools and frameworks:
* [Golang](https://getbootstrap.com)
* [Python](https://getbootstrap.com)
* [AWS CDK](https://jquery.com)
* [newman](https://laravel.com)
* [jq](https://laravel.com)
* [mosquitto](https://laravel.com)


<!-- GETTING STARTED -->
## Getting Started

To get started with this quickstart, follow the steps below (make sure to follow the prerequisit secctiono first)

### Prerequisites

#### Service Quotas

You need to ensure the following quotas requirements are met in your account. If not the case, create a request a request from the [AWS Service Quotas Dashboard](https://console.aws.amazon.com/servicequotas/home/)

Input         |    Quota Name |         Value | Comments
------------- | ------------- | ------------- | -------------
CodeBuild | Maximum number of concurrent running builds*	| 2 | The QuickStart uses a pipeline with paralel build steps. you need at this 2 concurent builds allowed. see [AWS CodeBuild Quotas](https://docs.aws.amazon.com/codebuild/latest/userguide/limits.html)

#### AWS SSO activation (Optional if you don'd want to use AWS IOT dashboord sitewise)
AWS SSO provides identity federation for SiteWise Monitor so that you can control access to your portals. With AWS SSO, your users sign in with their corporate email and password instead of an AWS account Follow the steps under Enabling AWS SSO in the [AWS IOT sitewise documentationn](https://docs.aws.amazon.com/iot-sitewise/latest/userguide/monitor-getting-started.html)

#### Create a Quicksight account (Optional if you don't want to use AWS Quicksigth dashboard or are already signed up)
If you haven't already, sign up for quicksight using the steps inthe [AWS documentation](https://docs.aws.amazon.com/quicksight/latest/user/signing-up.html). If you plan to deploy the default dashboard, you need an QuickSight Enterprise account

#### Validate your email adrdess with SES
This quickstart uses the email address provided in input form the AWS cloudFormation template as both sender and receiver of email notification. These notification will provide you with the key credentials to use the device onboarding MIicroservice. More specifically, for users of the Rigado [Alegro Kit](https://www.rigado.com/market-solutions/smart-hospitality-retail-solutions-powered-by-aws-iot/?did=pa_card&trk=pa_card), the email will provide the data necessary to use the Rigado Wizard to automatically onboard the Rigado Gateway. In order to be able to use this email address, SES requires you to verify the provided email as described in the [Amazon SES Documentation](https://docs.aws.amazon.com/ses/latest/DeveloperGuide/verify-email-addresses.html). 

#### Enable Logginng for AWS IOT Core
Enablling logging for AWS IOT Core will facilitate troubleshooting of device connectivity. This is especially useful is you are not using a rigado device. See instruction in the [AWS IOT Core documentation](https://docs.aws.amazon.com/iot/latest/developerguide/configure-logging.html)


### Installation

1. To get started with the deployment download the [AWS CloudFormation template](iot-onboarding-code-pipelines/iot-onboarding-int.yml). Note that you can also fork this repository and use it as a base for your own IOT project.

2. Go to the AWS cloudFormation console and launch the stack. The folowing parametters are required inputs
![Alt text](images/cloudformation-form.png?raw=true "Title")

Input  | Description
------------- | -------------
contactEmail  | Email of an administrator ussed for the AWS IOT sitewise portal creation. (see [AWS documentation](https://docs.aws.amazon.com/iot-sitewise/latest/userguide/administer-portals.html#portal-change-admins))
quickSightAdminUserName  | The use name of a Quicksight user with an ADMIN role. You can list all quicksight users by going to the [QuickSight administration Screen](https://us-east-1.quicksight.aws.amazon.com/sn/admin).
rootMqttTopic  | The root MQTT topic your devices publish to. If you are using the Rigado [Alegro Kit](https://www.rigado.com/market-solutions/smart-hospitality-retail-solutions-powered-by-aws-iot/?did=pa_card&trk=pa_card), the default value (data/#) does not need to be changed.
sourceTemplateArn  | use: **arn:aws:quicksight:eu-central-1:660526416360:template/iotOnboardingRigadoQuicksightPublicTemplatedev** This is a static location of a public quicksight dashboard template that we created for the purpose of this Quickstart. This allows you to get started quickly with a fully functional dashboard. Note that this example dahsbord is created specifically for the users of the Rigado [Alegro Kit](https://www.rigado.com/market-solutions/smart-hospitality-retail-solutions-powered-by-aws-iot/?did=pa_card&trk=pa_card). If you are not using rigado devices you will need to create you own dataset, analysis and dahsboard based on the devices you use. The ETL process and Glue Crawler that ingest the data from the IOT Broker to make them available in quicksight are data-model agnostic so you just need to link the created glue table as a datasource in QuickSight

Once the template run is successful, go to [AWS Code Pipeline](https://console.aws.amazon.com/codesuite/codepipeline). You should see the pipeline executing as below:
![Alt text](images/quickstart-cicd-3.png?raw=true "Title")

If you click on the pipeline name, you can see the steps of the pipeline running:
![Alt text](images/quickstart-cicd-2.png?raw=true "Title")


## Connecting Devices
### Rigado Devices
If you are using Rigado devices using the [Alegro Kit](https://www.rigado.com/market-solutions/smart-hospitality-retail-solutions-powered-by-aws-iot/?did=pa_card&trk=pa_card), go to the rigado wizard as explain in the email received after the devices are activated. Enter the data as received in the email you should have received from the QuickStart script:
```
AWS IOT Connectivity QuickStart Output Values


--------------------------------------------------------------------------------------------
| Cognito URL | https://iot-onboarding-quickstart-<env>.auth.us-east-1.amazoncognito.com/oauth2/token
--------------------------------------------------------------------------------------------
| API Gateway URL | https://<api_id>.execute-api.us-east-1.amazonaws.com/
--------------------------------------------------------------------------------------------
| Client ID | 228v...t9c3
--------------------------------------------------------------------------------------------
| Refresh Token | eyJjdHkiOiJKV1QiLCJl...slN29FrDNqHWo_0e5U85ow
```
Following the completion of the Rigado Wizard flow, the Rigado gateway will be setup to automatically and securely send traffic to AWS IOT Core using MQTT. To validat the the traffic is flowing through correctly, go to the AWS IOT Core console and subscribe top the topic provided in input of the CloudFromationn Template (by default data/# for Rigado). Provided at least one sensor has been turned on, you should see messages flowing through as shown below. 
![Alt text](images/iot-core-mqtt-test.png?raw=true "Title")

Once you validated that the traffic is flowing throught, see section **Visualizing the data** below to visualize the data.

### Non-Rigado Devices
For non rigado devices, you need to manually configure you device in ordre to alloow it to communicate with AWS IOT Core. The Lambda function deployed by the quickstart allows you to generate the device certificate and keypair to be used for secured communication between the device and IOT Core. It also creates the IOT Thing and appropriate policy. It returns:
-	The Mqtt Endpoint of the AWS IOT Core broker. 
-	The certificate generated on the fly 
-	The keypair (public and private key) generated on the fly
-	Other information about the device

Below and example of exchange with the onboarding service

We first get the session token by using the client ID and the refresh token

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
The response contains the necessary certificates and keys to be used to configure the device. Refer to the [AWS IOT Core documentation](https://docs.aws.amazon.com/iot/latest/developerguide/connect-to-iot.html) for details on configuring IOT device to communicate with AWS.

## Checking Connnectivity
In this section we will validate that the divices have been setup correctly and the traffic is flowing as expected.

## Visualizing the data  (Alegro Kit Users only)
As describe in the archotecture above, the data sent to the IOT Core MQTT broker is ingested into an IOT Datalake for aposteriori visualization and into AWS SiteWise for live monitoring. The paragraphe below describe how to sccess these visualization.
### Visualizing the data using Amazon QuickSight
In ordre to speed up your IOT project, this QuickStart deploys a predefinned dahsboard within AWS QuickSight. To benefit from this feature, you need:
- To have an AWS QuickSight Enterprise customer
- To be using the Rigado Allegro Kit as devices (while the Glue ETL ingestinng and processing the data is not device agnostic, the Dahsbords does make assumption onf the name of the field received annd will therefore only work out-of-the-box for Rigado Allegro Kit users)
By Default, the data is processed by a scheduled Glue Crawler and ETL ever 24H. this default value is choosen to minimize the cost of running the IOT datalake initially and can be updated easily by changing a CRON expression in the AWS CDK script or directly from the AWS console. 
When accessing the AWS QuickSight Dashbord for the first time, you need to provide access to the Amazon S3 Bucket that contains the refined data (by refined data, we mean the data processed by the ETL script). Giving QuickSIght Access to Amazon S3 buckert is described in the [AWS documentation](https://docs.aws.amazon.com/quicksight/latest/user/troubleshoot-connect-athena.html).
Following the deployment of the QuickStart, you should see ann analysis called :Rigado QuickStart Dahsbord in your quicksight account as shown below.
![Alt text](images/rigado-dahsboard.png?raw=true "Title")
This Dashboard is configured to query 48 hours of data in the past (this is to limit both cost and improve dashboard load time as the quantity of data increases in the future). THere are multiple ways you can change this setup while scaling with large amount of data by using [QuickSight SPICE](https://docs.aws.amazon.com/quicksight/latest/user/spice.html). Note that using SPICE will come with an additional cost.

**Note for non-alegro Kit user:** If you are not an allegro kit user, you will need to create you own Analysis and  datasource targeting the Athena Table for refined data mentioned earlier. This can be done in just a few clicks following the AWS QUickSIght documenttation. Note that the Glue job that refined the data is device agnistic as it justs flatten the JSON nested fields. It may, however not lead to practical result for deeply nested data.

### Visualizing the data using AWS SiteWise
Similarly to teh QuickSight dashboard the QuickStart created Sitewise Assets Model creates 1 root asset model and 4 children assets models. It also creates a Portal. In order to start visualizing the data in the prortal, you need to following teh setps below:
1. Go to AWS IoT QuickSight and select Build > Models
2. Choose the Asset model that correcponds to your Rigado Device (if the device you are using does not correspond to any existing asset model, refer to AWS IOT Sitewise documentation to route the traffic of your device to the aporpriate alias)
3. Create an asset for this asset model using the deviceId in the device name
4. Once created, go to "Edit" and enter a property alias for each of the model Measurement. For consistency with the IOT Core Broker rule, the alias value must be as follow:
```
<deviceId><MeasurementNameWithoutDoubleQuotes>
```
See example below for device **ffcfed4dd3ab**
![Alt text](images/sitewise-property-alias-setup.png?raw=true "Title")
Repeat this for all devices sending traffic behing the Rigado Gateway. (not that, using the Quicksight Dashboard, you can have a list of all devices sending traffic though the Gateway)

5. Once the asset is created youu can access the portal created by the QuickStart or create a portal from scratch following the AWS IOT Sitewise documentation. It will then just take a few minutes to add your assets to dedicated dashboards.


From this point, you can use the created portal to design dashboard for your devices as descrived in the AWS IOT Sitewise documentation.

**Note for non-alegro Kit user:** If you are not an allegro kit user, you will need to create your own AWS IOT Core Broker rule (following the same model than gthe one created in thei quickStart) to ingest the properly formated data into SiteWise. You will also need to manually create the Assets models and Assets following teh AWS SIteWise documentation.


## Cleaninng Up
In this quickstart, we use a combination of CLI and CDK for AWS Resources deployement. This is because some services like QuickSight and Sitewise are not supported by CloudFromation jut yet. Consequently, several manual steps will be required to clean up the deployed resources in the use account. These steps are described below:
1. Empty the S3 buckets
Identify the bucket created by the stack (they are prefixed by "iotonboardinginfrastack") and ensure you clean the content of these buckets before deleting the stack.
2. Delete infra stack
Go to CloudFormation and Delete the infrastructuure stack named **IOTOnboardingInfraStackint**
3. Delete Code Pipeline stack
Go to CloudFormation and Delete code pipeline stack you created
4. Clean up QuickSight Dashboard
You can manually delete the resources created in quicksight followiing the AWS Quicksight Documentation. Note that, if youo created a QuickSIght Account just for the purpose of this QuickStart you can unsubscribe to the service by following the steps described [here](https://docs.aws.amazon.com/quicksight/latest/user/closing-account.html)
5. Clean up Sitewise Dashboard
You need to delete the following resources (deletion procedure is described in teh AWS Documentation),:
- SiteWise Assets.
- Sitewise Assets Models (the quickStart creates 1 root asset model andn 4 children assets models).
- Sitewise Projects and Dashboards.

## FAQ

### My Quicksight dashborad display "error"
Errors displayed in the widgets of the QuickSight Dashboard typically results form one of the following causes:
1. Access to the S3 Bucket created by the CDK script
**solution:** Follow the AWS QuickSightt instruction to provide access to the bucket. The Bucket name
2. No Data is available yet in the datalake (data either not yet crawled or not yet received)

### My quicksight deployment script fails with error
### My code pipeline action fails with error about concurent build allowed
Some AWS accounnts are by default configure to only allow 1 concurent build using Code Pipeline. A workaround to this issue is to retry the pipeline stage. this will rety onnly the failed build. A long term solution consists in requesting a Code Build Service limite increase by contactinng AWS Support.

### My Sitewise script fails with ResourceAlreadyExistsException
```
An error occurred (ResourceAlreadyExistsException) when calling the CreateAssetModel operation: Another resource is already using the name RigadoHoboMX100QsTestint 
```
**solution:** Delete the Sitewise resources

### I am using the Alegro kit and I don't see any column in my Quicksight dataset
Data can take up to one hour to propagate to IOT Datalake due to the schedule chosen for the Crawler for cost optimization purpose.
**solution:** Ensure data is flowing using IOT Core Monitor and manually trgger the crawlers and jobs.

### Deletinng the infrastructure stack fails
**solution:** Empty and Delete the buckets


<!-- LICENSE -->
## License

Distributed under the MIT License. See `LICENSE` for more information.


<!-- CONTACT -->
## Contact

Quickstart Team - [@your_twitter](https://twitter.com/your_username) - email@example.com

Project Link: [https://github.com/your_username/repo_name](https://github.com/your_username/repo_name)



<!-- ACKNOWLEDGEMENTS -->
## Acknowledgements


