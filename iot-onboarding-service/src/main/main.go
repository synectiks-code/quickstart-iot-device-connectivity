package main

import (
	core "cloudrack-lambda-core/core"
	db "cloudrack-lambda-core/db"
	iot "cloudrack-lambda-core/iot"
	model "cloudrack-lambda-fn/model"
	usecase "cloudrack-lambda-fn/usecase"
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

////////////////////////////////////////////////////////////////////////////////////////////////////
// This Lambda function receives requests to the onboarding service throught the API Gateway
// Endpoint. Requests can be:
// - POST /api/onboard/{id}: create a new IOT thing, a new certificate, keypair and stores the data in DynamoDB with the device serial number
// - GET  /api/onboard/{id} : retreives the created thing and securtitey credentials
// - DELETE /api/onboard/{id}: Delete the onboarded device and security credentials
//////////////////////////////////////////////////////////////////////////////////////////////

//environment Name passed from CDK
var LAMBDA_ENV = os.Getenv("LAMBDA_ENV")

//
var LAMBDA_REGION = os.Getenv("AWS_REGION")

var ONBOARDING_TABLE_NAME = os.Getenv("ONBOARDING_TABLE_NAME")
var ONBOARDING_TABLE_PK = os.Getenv("ONBOARDING_TABLE_PK")
var ONBOARDING_TABLE_SK = os.Getenv("ONBOARDING_TABLE_SK")

var ONBOARDING_TABLE_SK_PREFIX = "serial"
var FN_ONBOARD_DEVICE = "onboard_device"
var FN_RETRIEVE_ONBOARDED_DEVICE = "retrieve_onboarded_device"
var FN_DELETE_ONBOARDED_DEVICE = "delete_onboarded_device"

//https://docs.aws.amazon.com/iot/latest/developerguide/server-authentication.html
var CA_CERTTIFICATE_URL_CA1 = "https://www.amazontrust.com/repository/AmazonRootCA1.pem"
var CA_CERTTIFICATE_URL_CA2 = "https://www.amazontrust.com/repository/AmazonRootCA2.pem"
var CA_CERTTIFICATE_URL_CA3 = "https://www.amazontrust.com/repository/AmazonRootCA3.pem"
var CA_CERTTIFICATE_URL_CA4 = "https://www.amazontrust.com/repository/AmazonRootCA4.pem"
var CA_CROSSSIGNED_CERTTIFICATE_URL_CA1 = "https://www.amazontrust.com/repository/G2-RootCA1.pem"
var CA_CROSSSIGNED_CERTTIFICATE_URL_CA2 = "https://www.amazontrust.com/repository/G2-RootCA2.pem"
var CA_CROSSSIGNED_CERTTIFICATE_URL_CA3 = "https://www.amazontrust.com/repository/G2-RootCA3.pem"
var CA_CROSSSIGNED_CERTTIFICATE_URL_CA4 = "https://www.amazontrust.com/repository/G2-RootCA4.pem"
var CA_STARTFIELD_ROOT_CERTIFICATE_URL = "https://www.amazontrust.com/repository/SFSRootCAG2.pem"

var onboardingDB = db.Init(ONBOARDING_TABLE_NAME, ONBOARDING_TABLE_PK, ONBOARDING_TABLE_SK)
var iotCore = iot.Init("data")

func HandleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("[ONBOARDING] Received Request %+v with context %+v", req, ctx)
	//0-Adding content to dynamoDB client to be able to use Xray Training
	onboardingDB.AddContext(ctx)
	//1-Identifying use case from HTTP Method and path
	useCase := identifyUseCase(req.Resource, req.HTTPMethod)
	var err error
	var res model.Response
	//2-Validating request and creating request wrapper
	onboardingRequest, err := createOnBoardingRequest(req)
	if err != nil {
		log.Printf("[ONBOARDING] Use Case %s failed with error: %v", useCase, err)
		return builResponseError(err), nil
	}
	//3-Get IOT Data-ATS Enpoint for the account
	mqttEndpoint, err2 := iotCore.GetEndpoint()
	if err2 != nil {
		log.Printf("[ONBOARDING] Use Case %s failed while retreiving MQTT Endpoint with error: %v", useCase, err)
		return builResponseError(err2), nil
	}
	//3-Selecting apropriate use case
	if useCase == FN_RETRIEVE_ONBOARDED_DEVICE {
		log.Printf("[ONBOARDING]Selected Use Case %v", useCase)
		if err == nil {
			res, err = usecase.Retrieve(onboardingRequest, onboardingDB, mqttEndpoint)
			if err == nil {
				log.Printf("[ONBOARDING] Use Case %s successult with response: %v", useCase, res)
				return builResponse(res), nil
			}
		}
		return builResponseError(err), nil
	} else if useCase == FN_ONBOARD_DEVICE {
		log.Printf("[ONBOARDING] Selected Use Case %v", useCase)
		if err == nil {
			res, err = usecase.Create(onboardingRequest, onboardingDB, iotCore, mqttEndpoint)
			if err == nil {
				log.Printf("[ONBOARDING] Use Case %s successult with response: %v", useCase, res)
				return builResponse(res), nil
			}
		}
		return builResponseError(err), nil
	} else if useCase == FN_DELETE_ONBOARDED_DEVICE {
		log.Printf("[ONBOARDING] Selected Use Case %v", useCase)
		if err == nil {
			res, err = usecase.Delete(onboardingRequest, onboardingDB, iotCore)
			if err == nil {
				log.Printf("[ONBOARDING] Use Case %s successult with response: %v", useCase, res)
				return builResponse(res), nil
			}
		}
		return builResponseError(err), nil
	}
	//4-if no use case found, return an error
	log.Printf("No Use Case Found for %v", req.HTTPMethod+" "+req.Resource)
	err = errors.New("No Use Case Found for " + req.HTTPMethod + " " + req.Resource)
	return builResponseError(err), nil
}

//This function identifies the use case based on HTTP Method and path
//return an empty string id no use case if founf
func identifyUseCase(res string, meth string) string {
	if res == "/onboard/{id}" && meth == "POST" {
		return FN_ONBOARD_DEVICE
	}
	if res == "/onboard/{id}" && meth == "GET" {
		return FN_RETRIEVE_ONBOARDED_DEVICE
	}
	if res == "/onboard/{id}" && meth == "DELETE" {
		return FN_DELETE_ONBOARDED_DEVICE
	}
	return ""
}

//This function validates the request and creates a wrapper struct
//in our case the serialNumber ofthe device is mandatory
func createOnBoardingRequest(req events.APIGatewayProxyRequest) (model.Request, error) {
	serialNumber := req.PathParameters["id"]
	if serialNumber != "" {
		return model.Request{SerialNumber: serialNumber}, nil
	}
	return model.Request{}, errors.New("Invalid request: device serial number is missing")
}

//Building properly formated SUCCESS response for API gateway to process
func builResponse(resWrapper model.Response) events.APIGatewayProxyResponse {
	res := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	jsonRes, err := json.Marshal(resWrapper)
	if err != nil {
		log.Printf("[ONBOARDING] Error while unmarshalling response %v", err)
		return builResponseError(err)
	}
	res.Body = string(jsonRes)
	return res
}

//Building properly formated ERROR response for API gateway to process
func builResponseError(err error) events.APIGatewayProxyResponse {
	log.Printf("[ONBOARDING] Response Error: %v", err)
	res := events.APIGatewayProxyResponse{
		StatusCode: 400,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	json, _ := json.Marshal(model.Response{Error: core.BuildResError(err)})
	res.Body = string(json)
	return res
}

func main() {
	lambda.Start(HandleRequest)
}
