package lambda

import (
	core "cloudrack-lambda-core/core"
	"encoding/json"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	awsLambda "github.com/aws/aws-sdk-go/service/lambda"
)

type Config struct {
	Client       *awsLambda.Lambda //sdk client to make call to the AWS API
	FunctionName string
}

func Init(function string) Config {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	return Config{
		Client:       awsLambda.New(sess),
		FunctionName: function,
	}
}

func (c Config) Invoke(payload interface{}, res core.CloudrackObject) (interface{}, error) {
	marshalledPayload, err := json.Marshal(payload)
	if err != nil {
		log.Printf("[LAMBDA] Error marshaling request for lambda invoke: %+v", err)
		return res, err
	}
	input := &awsLambda.InvokeInput{
		Payload:      marshalledPayload,
		FunctionName: aws.String(c.FunctionName),
	}
	log.Printf("[LAMBDA] Lambda Request: %+v", input)
	output, err2 := c.Client.Invoke(input)
	if err2 != nil {
		log.Printf("[LAMBDA] Error invoking lambda: %+v", err2)
		return res, err2
	}
	log.Printf("[LAMBDA] Lambda Response: %+v", output)
	err, res = res.Decode(*json.NewDecoder(strings.NewReader(string(output.Payload))))
	if err != nil {
		log.Printf("[LAMBDA] Error unmarshaling response for lambda invoke: %+v", err)
		return res, err
	}
	log.Printf("[LAMBDA] DecodedLambda Response: %+v", res)
	return res, err
}
