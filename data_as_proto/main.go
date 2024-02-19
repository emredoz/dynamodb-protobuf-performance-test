package main

import (
	"context"
	"dynamodb-protobuf-performance-test/model"
	"dynamodb-protobuf-performance-test/pkg/gaws"
	"dynamodb-protobuf-performance-test/pkg/settings"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"google.golang.org/protobuf/proto"
)

func init() {
	settings.Setup()
	gaws.Setup(settings.AwsDynamoDbSettings)
}

func main() {
	//saveDataAsProto()
	readDataAsProto()
}

type DynamoDataItem struct {
	PK                string
	SK                int
	MessagesByteArray []byte
}

func saveDataAsProto() {
	db := gaws.GetDb()
	dynamoDataItem := DynamoDataItem{
		PK: "pkey1",
		SK: 1,
	}
	messages := model.CreateMessages()

	protoMarshalStart := time.Now()
	dynamoDataItem.MessagesByteArray, _ = proto.Marshal(messages)
	protoMarshalEnd := time.Now()
	fmt.Printf("Proto marshal süresi: %s\n", protoMarshalEnd.Sub(protoMarshalStart))

	dynamoMarshalStart := time.Now()
	item, _ := dynamodbattribute.MarshalMap(dynamoDataItem)
	dynamoMarshalEnd := time.Now()
	fmt.Printf("Dynamo marshal süresi: %s\n", dynamoMarshalEnd.Sub(dynamoMarshalStart))
	fmt.Printf("Toplam marshal süresi: %s\n", dynamoMarshalEnd.Sub(protoMarshalStart))

	input := &dynamodb.PutItemInput{
		TableName: aws.String(settings.AwsDynamoDbSettings.TableName),
		Item:      item,
	}
	dynamoSaveStart := time.Now()
	_, _ = db.PutItemWithContext(context.Background(), input)
	dynamoSaveEnd := time.Now()
	fmt.Printf("DynamoDB save süresi: %s\n", dynamoSaveEnd.Sub(dynamoSaveStart))
	fmt.Printf("Toplam süre: %s\n", dynamoSaveEnd.Sub(protoMarshalStart))
}

func readDataAsProto() {
	db := gaws.GetDb()
	input := &dynamodb.GetItemInput{
		TableName: aws.String(settings.AwsDynamoDbSettings.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {S: aws.String("pkey1")},
			"SK": {N: aws.String("1")},
		},
	}
	response := &DynamoDataItem{}
	messages := &model.Messages{}

	dynamoReadStart := time.Now()
	result, _ := db.GetItemWithContext(context.Background(), input)
	dynamoReadEnd := time.Now()
	fmt.Printf("DynamoDB read süresi: %s\n", dynamoReadEnd.Sub(dynamoReadStart))

	dynamoUnmarshalStart := time.Now()
	_ = dynamodbattribute.UnmarshalMap(result.Item, &response)
	dynamoUnmarshalEnd := time.Now()
	fmt.Printf("Dynamo unmarshal süresi: %s\n", dynamoUnmarshalEnd.Sub(dynamoUnmarshalStart))

	protoUnmarshalStart := time.Now()
	_ = proto.Unmarshal(response.MessagesByteArray, messages)
	protoUnmarshalEnd := time.Now()
	fmt.Printf("Proto unmarshal süresi: %s\n", protoUnmarshalEnd.Sub(protoUnmarshalStart))
	fmt.Printf("Toplam unmarshal süresi: %s\n", protoUnmarshalEnd.Sub(dynamoUnmarshalStart))
	fmt.Printf("Toplam süre: %s\n", protoUnmarshalEnd.Sub(dynamoReadStart))
	fmt.Println(messages.MessageList[0].To)
}
