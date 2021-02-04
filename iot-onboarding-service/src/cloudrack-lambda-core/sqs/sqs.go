package sqs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type Config struct {
	Client *sqs.SQS
}

func Init() Config {
	mySession := session.Must(session.NewSession())
	svc := sqs.New(mySession)
	return Config{
		Client: svc,
	}
}

func (c Config) Create(queueName string) (string, error) {
	input := &sqs.CreateQueueInput{
		QueueName: aws.String(queueName),
	}
	out, err := c.Client.CreateQueue(input)
	if err != nil {
		return "", err
	}
	return *out.QueueUrl, nil
}

func (c Config) Delete(queueURL string) error {
	input := &sqs.DeleteQueueInput{
		QueueUrl: aws.String(queueURL),
	}
	_, err := c.Client.DeleteQueue(input)
	if err != nil {
		return err
	}
	return nil
}

func (c Config) Send(queueUrl string, dedupID string, msg string) error {
	input := &sqs.SendMessageInput{
		MessageBody:            aws.String(msg),
		MessageDeduplicationId: aws.String(dedupID),
		QueueUrl:               aws.String(queueUrl),
	}
	_, err := c.Client.SendMessage(input)
	return err
}

func (c Config) Get(queueUrl string) ([]string, error) {
	input := &sqs.ReceiveMessageInput{
		MaxNumberOfMessages: aws.Int64(10),
		VisibilityTimeout:   aws.Int64(1),
		WaitTimeSeconds:     aws.Int64(1),
		QueueUrl:            aws.String(queueUrl),
	}
	res, err := c.Client.ReceiveMessage(input)
	if err != nil {
		return []string{}, err
	}
	response := []string{}
	for _, msg := range res.Messages {
		response = append(response, *msg.Body)
	}
	return response, err
}
