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
#### remove git credentials
#### change qQS template URL
#### Add instructions for SSO and QS account creatiion
#### Add cleaning instructions for dahsbord
#### Add instruction to get input for Rigado wizard
#### add where to find quicksight admin user name

<!-- ABOUT THE PROJECT -->
## About The Project

This AWS quickstart aims at hellping AWS IOT customers to quickly get started with an IOT landing zone on AWS including:
* A REST microservice to onboard devices and gateway by serial number. The service creates the AWS IOT Core resources to secuerly connect to AWS MQTT Broker.
* An IOT Datalake ingesting the data from the long term storage and analytics
* An example AWS Quicksight Dahsboard to display data form the datalake (Compatible devices only)
* An example IOT Device real time Monitorinig dahsboard uusing AWS IOT SiteWise (Compatible Devices Only)

![Alt text](images/iot-quickstart-archtecture.png?raw=true "Title")


The QuickStart is being released in partnership with AWS IOT and Travel adn Hospitality Competency Partner Rigado and compatible with [Rigado](https://www.rigado.com/market-solutions/smart-hospitality-retail-solutions-powered-by-aws-iot/?did=pa_card&trk=pa_card) newly launched [Alegro Kit](https://www.rigado.com/market-solutions/smart-hospitality-retail-solutions-powered-by-aws-iot/?did=pa_card&trk=pa_card)

### Built With

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


This section should list any major frameworks that you built your project using. Leave any add-ons/plugins for the acknowledgements section. Here are a few examples.
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

#### AWS SSO activation (Optional if you don'd want to use AWS IOT dashboord sitewise)
AWS SSO provides identity federation for SiteWise Monitor so that you can control access to your portals. With AWS SSO, your users sign in with their corporate email and password instead of an AWS account Follow the steps under Enabling AWS SSO in the [AWS IOT sitewise documentationn](https://docs.aws.amazon.com/iot-sitewise/latest/userguide/monitor-getting-started.html)

#### Create a Quicksight account (Optional if you don't want to use AWS Quicksigth dashboard or are already signed up)
If you haven't already, sign up for quicksight using the steps inthe [AWS documentation](https://docs.aws.amazon.com/quicksight/latest/user/signing-up.html)

### Installation

To get started with the deployment download the [AWS CloudFormation template](https://github.com/aws-quickstart/quickstart-iot-device-connectivity/raw/main/iot-onboarding-int.yml)
Not that you can also fork this repository and use it as a base for your own IOT project.

<!-- USAGE EXAMPLES -->
## Usage

Use this space to show useful examples of how a project can be used. Additional screenshots, code examples and demos work well in this space. You may also link to more resources.

_For more examples, please refer to the [Documentation](https://example.com)_



<!-- ROADMAP -->
## Roadmap

See the [open issues](https://github.com/othneildrew/Best-README-Template/issues) for a list of proposed features (and known issues).



<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to be learn, inspire, and create. Any contributions you make are **greatly appreciated**.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request



<!-- LICENSE -->
## License

Distributed under the MIT License. See `LICENSE` for more information.



<!-- CONTACT -->
## Contact

Your Name - [@your_twitter](https://twitter.com/your_username) - email@example.com

Project Link: [https://github.com/your_username/repo_name](https://github.com/your_username/repo_name)



<!-- ACKNOWLEDGEMENTS -->
## Acknowledgements


