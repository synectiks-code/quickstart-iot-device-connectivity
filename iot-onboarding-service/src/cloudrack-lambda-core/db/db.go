package db

import (
	core "cloudrack-lambda-core/core"
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

type Notification struct {
	MessageType string `json:"messageType"`
}

type DynamoRecord struct {
	Pk string `json:"pk"`
	Sk string `json:"sk"`
}

type DynamoUserConnection struct {
	UserId      string `json:"userId"`
	ConnetionId string `json:"connectionId"`
}

type DynamoEventChange struct {
	NewImage map[string]*dynamodb.AttributeValue `json:"NewImage"`
	OldImage map[string]*dynamodb.AttributeValue `json:"OldImage"`
}

type DynamoEventRecord struct {
	Change    DynamoEventChange `json:"dynamodb"`
	EventName string            `json:"eventName"`
	EventID   string            `json:"eventID"`
	// ... more fields if needed: https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_streams_GetRecords.html
}

type DynamoEvent struct {
	Records []DynamoEventRecord `json:"records"`
}

type DBConfig struct {
	DbService     *dynamodb.DynamoDB
	PrimaryKey    string
	SortKey       string
	TableName     string
	LambdaContext context.Context
}

type DynamoFilterExpression struct {
	BeginsWith []DynamoFilterCondition
	Contains   []DynamoFilterCondition
}

type DynamoFilterCondition struct {
	Key   string
	Value string
}

func (exp DynamoFilterExpression) HasFilter() bool {
	return exp.HasBeginsWith() || exp.HasContains()
}
func (exp DynamoFilterExpression) HasBeginsWith() bool {
	if len(exp.BeginsWith) > 0 && exp.BeginsWith[0].Key != "" {
		return true
	}
	return false
}
func (exp DynamoFilterExpression) HasContains() bool {
	if len(exp.Contains) > 0 && exp.Contains[0].Key != "" {
		return true
	}
	return false
}

//init setup teh session and define table name, primary key and sort key
func Init(tn string, pk string, sk string) DBConfig {
	if pk == "" {
		fmt.Println("[CORE][DB] WARNING: empty PK provided")
	}
	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	dbSession := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	// Create DynamoDB client
	return DBConfig{
		DbService:  dynamodb.New(dbSession),
		PrimaryKey: pk,
		SortKey:    sk,
		TableName:  tn,
	}

}

//add lambda execution context to db connection structto be able to be used for tracing
//this is not done in the init function to be able to initialize the connectino
//outside the lambda function handler
func (dbc *DBConfig) AddContext(ctx context.Context) error {
	dbc.LambdaContext = ctx
	return nil
}

func ValidateConfig(dbc DBConfig) error {
	if dbc.LambdaContext == nil {
		return errors.New("Lambda Context is Empty. Please call AddContext function after Init")
	}
	if dbc.PrimaryKey == "" {
		return errors.New("Cannot have an empty PK in DB config")
	}
	return nil
}

func (dbc DBConfig) Save(prop interface{}) (interface{}, error) {
	err := ValidateConfig(dbc)
	if err != nil {
		return nil, err
	}
	av, err := dynamodbattribute.MarshalMap(prop)
	if err != nil {
		fmt.Println("Got error marshalling new property item:")
		fmt.Println(err.Error())
	}
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(dbc.TableName),
	}

	_, err = dbc.DbService.PutItemWithContext(dbc.LambdaContext, input)
	if err != nil {
		fmt.Println("Got error calling PutItem:")
		fmt.Println(err.Error())
	}
	return prop, err
}

func (dbc DBConfig) Delete(prop interface{}) (interface{}, error) {
	av, err := dynamodbattribute.MarshalMap(prop)
	if err != nil {
		fmt.Println("Got error marshalling new property item:")
		fmt.Println(err.Error())
	}

	input := &dynamodb.DeleteItemInput{
		Key:       av,
		TableName: aws.String(dbc.TableName),
	}

	_, err = dbc.DbService.DeleteItemWithContext(dbc.LambdaContext, input)
	if err != nil {
		fmt.Println("Got error calling DeetItem:")
		fmt.Println(err.Error())
	}
	return prop, err
}

//TODO: to evaluate th value of this tradeoff: this is probably a little slow but abstract the complexity for all uses of
//the save many function(and actually any core operation on array of interface)
func InterfaceSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}

//Writtes many items to a single table
func (dbc DBConfig) SaveMany(data interface{}) error {
	//Dynamo db currently limits batches to 25 items
	batches := core.Chunk(InterfaceSlice(data), 25)
	for i, dataArray := range batches {

		log.Printf("DB> Batch %i inserting: %+v", i, dataArray)
		items := make([]*dynamodb.WriteRequest, len(dataArray), len(dataArray))
		for i, item := range dataArray {
			av, err := dynamodbattribute.MarshalMap(item)
			if err != nil {
				fmt.Println("Got error marshalling new property item:")
				fmt.Println(err.Error())
			}
			items[i] = &dynamodb.WriteRequest{
				PutRequest: &dynamodb.PutRequest{
					Item: av,
				},
			}
		}

		bwii := &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]*dynamodb.WriteRequest{
				dbc.TableName: items,
			},
		}

		_, err := dbc.DbService.BatchWriteItem(bwii)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case dynamodb.ErrCodeProvisionedThroughputExceededException:
					fmt.Println(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
				case dynamodb.ErrCodeResourceNotFoundException:
					fmt.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
				case dynamodb.ErrCodeItemCollectionSizeLimitExceededException:
					fmt.Println(dynamodb.ErrCodeItemCollectionSizeLimitExceededException, aerr.Error())
				case dynamodb.ErrCodeRequestLimitExceeded:
					fmt.Println(dynamodb.ErrCodeRequestLimitExceeded, aerr.Error())
				case dynamodb.ErrCodeInternalServerError:
					fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				fmt.Println(err.Error())
			}
			return err
		}
	}
	return nil
}

//Deletes many items to a single table
func (dbc DBConfig) DeleteMany(data interface{}) error {
	//Dynamo db currently limits batches to 25 items
	batches := core.Chunk(InterfaceSlice(data), 25)
	for i, dataArray := range batches {

		log.Printf("DB> Batch %i deleting: %+v", i, dataArray)
		items := make([]*dynamodb.WriteRequest, len(dataArray), len(dataArray))
		for i, item := range dataArray {
			av, err := dynamodbattribute.MarshalMap(item)
			if err != nil {
				fmt.Println("Got error marshalling new property item:")
				fmt.Println(err.Error())
			}
			items[i] = &dynamodb.WriteRequest{
				DeleteRequest: &dynamodb.DeleteRequest{
					Key: av,
				},
			}
		}

		bwii := &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]*dynamodb.WriteRequest{
				dbc.TableName: items,
			},
		}

		_, err := dbc.DbService.BatchWriteItem(bwii)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case dynamodb.ErrCodeProvisionedThroughputExceededException:
					fmt.Println(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
				case dynamodb.ErrCodeResourceNotFoundException:
					fmt.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
				case dynamodb.ErrCodeItemCollectionSizeLimitExceededException:
					fmt.Println(dynamodb.ErrCodeItemCollectionSizeLimitExceededException, aerr.Error())
				case dynamodb.ErrCodeRequestLimitExceeded:
					fmt.Println(dynamodb.ErrCodeRequestLimitExceeded, aerr.Error())
				case dynamodb.ErrCodeInternalServerError:
					fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				fmt.Println(err.Error())
			}
			return err
		}
	}
	return nil
}

func (dbc DBConfig) GetMany(records []DynamoRecord, data interface{}) error {
	//Dynamo db currently limits read batches to 100 items OR up to 16 MB of data,
	//https: //docs.aws.amazon.com/sdk-for-go/api/service/dynamodb/#DynamoDB.BatchGetItem
	batches := core.Chunk(InterfaceSlice(records), 100)
	for i, dataArray := range batches {

		log.Printf("DB> Batch %v fetching: %+v", i, dataArray)
		items := make([]map[string]*dynamodb.AttributeValue, len(dataArray), len(dataArray))
		for i, item := range dataArray {
			rec := item.(DynamoRecord)
			av := map[string]*dynamodb.AttributeValue{
				dbc.PrimaryKey: {
					S: aws.String(rec.Pk),
				},
			}
			if rec.Sk != "" {
				av[dbc.SortKey] = &dynamodb.AttributeValue{
					S: aws.String(rec.Sk),
				}
			}
			items[i] = av
		}

		bgii := &dynamodb.BatchGetItemInput{
			RequestItems: map[string]*dynamodb.KeysAndAttributes{
				dbc.TableName: &dynamodb.KeysAndAttributes{
					Keys: items,
				},
			},
		}
		log.Printf("DB> DynamoDB Batch Get rq %+v", bgii)

		bgio, err := dbc.DbService.BatchGetItem(bgii)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case dynamodb.ErrCodeProvisionedThroughputExceededException:
					fmt.Println(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
				case dynamodb.ErrCodeResourceNotFoundException:
					fmt.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
				case dynamodb.ErrCodeItemCollectionSizeLimitExceededException:
					fmt.Println(dynamodb.ErrCodeItemCollectionSizeLimitExceededException, aerr.Error())
				case dynamodb.ErrCodeRequestLimitExceeded:
					fmt.Println(dynamodb.ErrCodeRequestLimitExceeded, aerr.Error())
				case dynamodb.ErrCodeInternalServerError:
					fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				fmt.Println(err.Error())
			}
			return err
		}
		//log.Printf("DB> DynamoDB Batch Get rs %+v", bgio)
		err = dynamodbattribute.UnmarshalListOfMaps(bgio.Responses[dbc.TableName], data)
		if err != nil {
			fmt.Println("DB:GetMany> Error while unmarshalling")
			fmt.Println(err.Error())
			return err
		}
		//log.Printf("DB> DynamoDB Batch Get unmarchalled rs %+v", data)
	}
	return nil
}

func (dbc DBConfig) Get(pk string, sk string, data interface{}) error {
	err := ValidateConfig(dbc)
	if err != nil {
		return err
	}
	av := map[string]*dynamodb.AttributeValue{
		dbc.PrimaryKey: {
			S: aws.String(pk),
		},
	}
	if sk != "" {
		av[dbc.SortKey] = &dynamodb.AttributeValue{
			S: aws.String(sk),
		}
	}

	gii := &dynamodb.GetItemInput{
		TableName: aws.String(dbc.TableName),
		Key:       av,
	}
	log.Printf("DB> DynamoDB  Get rq %+v", gii)
	result, err := dbc.DbService.GetItemWithContext(dbc.LambdaContext, gii)
	if err != nil {
		fmt.Println("NOT FOUND")
		fmt.Println(err.Error())
		return err
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, data)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}
	return err
}

func (dbc DBConfig) FindStartingWith(pk string, value string, data interface{}) error {
	var queryInput = &dynamodb.QueryInput{
		TableName: aws.String(dbc.TableName),
		KeyConditions: map[string]*dynamodb.Condition{
			dbc.PrimaryKey: {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(pk),
					},
				},
			},
			dbc.SortKey: {
				ComparisonOperator: aws.String("BEGINS_WITH"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(value),
					},
				},
			},
		},
	}
	fmt.Println("DB:FindStartingWith> rq: %+v", queryInput)

	var result, err = dbc.DbService.QueryWithContext(dbc.LambdaContext, queryInput)
	if err != nil {
		fmt.Println("DB:FindStartingWith> NOT FOUND")
		fmt.Println(err.Error())
		return err
	}

	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, data)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}
	return err
}

func (dbc DBConfig) FindStartingWithAndFilter(pk string, value string, data interface{}, filter DynamoFilterExpression) error {

	var queryInput = &dynamodb.QueryInput{
		TableName: aws.String(dbc.TableName),
		KeyConditions: map[string]*dynamodb.Condition{
			dbc.PrimaryKey: {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(pk),
					},
				},
			},
			dbc.SortKey: {
				ComparisonOperator: aws.String("BEGINS_WITH"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(value),
					},
				},
			},
		},
	}
	//Building Filter expression
	if filter.HasFilter() {
		builder := expression.NewBuilder()
		if filter.HasBeginsWith() {
			fmt.Println("DB:FindStartingWithAndFilter> Adding Filter expression: %+v", filter.BeginsWith[0])
			condBuilder := expression.Name(filter.BeginsWith[0].Key).BeginsWith(filter.BeginsWith[0].Value)
			builder = builder.WithFilter(condBuilder)
		}
		if filter.HasContains() {
			fmt.Println("DB:FindStartingWithAndFilter> Adding Filter expression: %+v", filter.Contains[0])
			condBuilder := expression.Name(filter.Contains[0].Key).Contains(filter.Contains[0].Value)
			builder = builder.WithFilter(condBuilder)
		}
		expr, err := builder.Build()
		if err != nil {
			fmt.Println("DB:FindStartingWithAndFilter> Error in filter expression: %+v", err)
		}
		queryInput.ExpressionAttributeNames = expr.Names()
		queryInput.ExpressionAttributeValues = expr.Values()
		queryInput.FilterExpression = expr.Filter()
	}

	fmt.Println("DB:FindStartingWith> rq: %+v", queryInput)

	//list for all paginated items
	allItems := []map[string]*dynamodb.AttributeValue{}
	var result, err = dbc.DbService.QueryWithContext(dbc.LambdaContext, queryInput)
	if err != nil {
		fmt.Println("DB:FindStartingWith> NOT FOUND")
		fmt.Println(err.Error())
		return err
	}
	for _, item := range result.Items {
		allItems = append(allItems, item)
	}
	//if paginated response
	for result.LastEvaluatedKey != nil {
		fmt.Println("DB:FindStartingWith> Response is paginated. LastEvaluatedKey: %v", result.LastEvaluatedKey)
		queryInput.ExclusiveStartKey = result.LastEvaluatedKey
		result, err = dbc.DbService.QueryWithContext(dbc.LambdaContext, queryInput)
		if err != nil {
			fmt.Println("DB:FindStartingWith> NOT FOUND")
			fmt.Println(err.Error())
			return err
		}
		for _, item := range result.Items {
			allItems = append(allItems, item)
		}
	}
	err = dynamodbattribute.UnmarshalListOfMaps(allItems, data)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}
	return err
}

func (dbc DBConfig) FindByGsi(value string, indexName string, indexPk string, data interface{}) error {
	var queryInput = &dynamodb.QueryInput{
		TableName: aws.String(dbc.TableName),
		IndexName: aws.String(indexName),
		KeyConditions: map[string]*dynamodb.Condition{
			indexPk: {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(value),
					},
				},
			},
		},
	}

	var result, err = dbc.DbService.QueryWithContext(dbc.LambdaContext, queryInput)
	if err != nil {
		fmt.Println("NOT FOUND")
		fmt.Println(err.Error())
		return err
	}

	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, data)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}
	return err
}
