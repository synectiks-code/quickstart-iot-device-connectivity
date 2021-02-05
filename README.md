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

THis repository includes the following folder:
* #### e2e
This folder contain end to end tests for the onboarding microservice and MQTT connectivity tests that validate that onboardded devices can connect to the AWS IOT Core MQTT Broker. It use newman, a CLI tool allowing to run test build using Postman and mosquitto as an MQTT CLient
* #### iot-onboarding-code-pipelines
This folder contains an AWS CDK project that builds a AWS Code Pipeline project which is used to deploy the architecture decribed above. We use this method to be able to provide consistent build experience for our CDk project independently forom builders environement (NOdeJS version...). The pipeline has the folloriing setps:
![Alt text](images/quickstart-cicd.png?raw=true "Title")

* #### iot-onboarding-data-processing
This folder contains a Python ETL (Extract Load Transform) script that flatten the device Json messages to be queried by Amazon Athena and Amazon Quicksight. This ETL script is run by a Glue Job
* #### iot-onboarding-infra
This folder contains a CDK project that builds most of the infrastructure components described above except the Quicksight and Sitewise Dahsboards which are not yet supported by AWS CloudFormation at the tiome of construyction of this quickstart.
* #### iot-onboarding-quicksight
This folder contains a linux shell script that automate the creation of an AWS Qhicksight Dahsboard based on a public template. Note that this requires for the target account to have activated Amazon QuickSight (add linik here). Also, the example dashboard assumes the following structure for MQTT messages from the device (based on AWS Partner Rigado). 

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
The base topic can be configured as an input parametter from the CICD pipeline CloudFromation Stack and the IOT Datalake uses glue crawlers to dynamically iidentify the data structure of the incoming MQTT messages. This means that QuickStart users who use different device configuration can quickly adap the dashboard to their specific need.

Note that using AWS CLI comes with limitations compare to CloudFormation and some resources (such as dashbord, dataset and datasource) may need to be manually deleted to be updated or in case of failure during deployemnet. We hope that providing this autromated dahsbord allows you to move faster by relying on an example and will move this to a more robust infrastructure as code solution when available.

The example dashboard looks as follows:

![Alt text](images/quicksight.png?raw=true "Title")


* #### iot-onboarding-service
THis folder contain the Golang Code of the onboarding service lambda function. The function sits behind an AWS API gateway REST API exposing the following Services
```
POST {{baseUrl}}api/onboard/{{deviceName}}
GET {{baseUrl}}api/onboard/{{deviceName}}
DELET {{baseUrl}}api/onboard/{{deviceName}}
```
These endpoint respetively create, retreive and delete a device or gateway including the following AWS IOT resources:
* A Device Certificate
* An IOT Thing and Associated policy to publish on the base topic provided as parameter to the quickstart CICD cloudformatioin template

THe service Create and Retreive enpoints return all the data needed to setup the device for AWS connectivity and the message structure is as follow:
```json
{
    "serialNumber": "<device serial number>",
    "deviceName": "<device name = device serial number>",
    "thingId": "<ID of teh AWS IOT Core THing>",
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
The service is secured by Amazon Cognito and a random user is created during infrastructure deployment along with a refresh token. To access te service, the quickStart owner wiill need to access credentials sorted in the S3 architect bucket following the stack successful buidl and generate temporary credentials in teh for of Cognito Access token. MOre information on this flow is porvided blow.

Note that, as part of the partnership with Rigado on this quickstart, the Rigado team created a Web Wizard for Alegro Kit user that takes care of generating the temporary credentials and setiing up the devices based on teh credentials generated by this Microservice. More information at [Rigado.com](add rigado kit url)

* #### iot-onboarding-sitewise
This forlder coontains a linux shell script that builds AWS IOT SiteWise resourtces needed too build a real time dashboard. THese resources include:
* a Device model hierarchy, composed of a root device and 4 child devices (based on the Rigado Alegro Kit content)
* A sitewise project and portal

When working with Rigado devices A few manual steps are required to create the assets and addes the to a dashoard and obtain the following result:
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

To get started with this quickstart, download the AWS CloudFormation template

Not that you can also fork this repository and use it as a base for your own IOT project.

### Prerequisites

This is an example of how to list things you need to use the software and how to install them.
* npm
  ```sh
  npm install npm@latest -g
  ```

### Installation

1. Get a free API Key at [https://example.com](https://example.com)
2. Clone the repo
   ```sh
   git clone https://github.com/your_username_/Project-Name.git
   ```
3. Install NPM packages
   ```sh
   npm install
   ```
4. Enter your API in `config.js`
   ```JS
   const API_KEY = 'ENTER YOUR API';
   ```



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


