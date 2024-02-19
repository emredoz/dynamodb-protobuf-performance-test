package main

import (
	"context"
	"dynamodb-protobuf-performance-test/model"
	"dynamodb-protobuf-performance-test/pkg/gaws"
	"dynamodb-protobuf-performance-test/pkg/settings"
	"encoding/json"
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
	//saveData()
	readData()
}

type DynamoDataItem struct {
	PK       string
	SK       int
	Messages *model.Messages
}

func saveData() {
	db := gaws.GetDb()
	dynamoDataItem := DynamoDataItem{
		PK:       "pkey2",
		SK:       1,
		Messages: model.CreateMessages(),
	}
	
	marshalStartTime := time.Now()
	putItem, _ := dynamodbattribute.MarshalMap(dynamoDataItem)
	marshalFinishTime := time.Now()
	fmt.Printf("Marshal süresi: %s\n", marshalFinishTime.Sub(marshalStartTime))

	input := &dynamodb.PutItemInput{
		TableName: aws.String(settings.AwsDynamoDbSettings.TableName),
		Item:      putItem,
	}
	dynamoSaveStart := time.Now()
	_, _ = db.PutItemWithContext(context.Background(), input)
	dynamoSaveEnd := time.Now()
	fmt.Printf("DynamoDB save süresi: %s\n", dynamoSaveEnd.Sub(dynamoSaveStart))
	fmt.Printf("Toplam süre: %s\n", dynamoSaveEnd.Sub(marshalStartTime))
}

func readData() {
	db := gaws.GetDb()
	getItemInput := &dynamodb.GetItemInput{
		TableName: aws.String(settings.AwsDynamoDbSettings.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {S: aws.String("pkey2")},
			"SK": {N: aws.String("1")},
		},
	}
	dynamoDataItem := &DynamoDataItem{}

	dynamoReadStart := time.Now()
	result, _ := db.GetItemWithContext(context.Background(), getItemInput)
	dynamoReadEnd := time.Now()
	fmt.Printf("DynamoDB read süresi: %s\n", dynamoReadEnd.Sub(dynamoReadStart))
	
	unmarshalStartTime := time.Now()
	_ = dynamodbattribute.UnmarshalMap(result.Item, &dynamoDataItem)
	unmarshalFinishTime := time.Now()
	
	fmt.Printf("Unmarshal süresi: %s\n", unmarshalFinishTime.Sub(unmarshalStartTime))
	fmt.Printf("Toplam süre: %s\n", unmarshalFinishTime.Sub(dynamoReadStart))
	fmt.Println(dynamoDataItem.Messages.MessageList[0].To)
}

func printDataSize() {
	messages := model.CreateMessages()
	jsonData, _ := json.Marshal(messages)
	sizeKB := float64(len(jsonData)) / 1024.0
	fmt.Printf("json veri Boyutu: %.2f KB\n", sizeKB)

	protoData, _ := proto.Marshal(messages)
	sizeKB = float64(len(protoData)) / 1024.0
	fmt.Printf("protobuf veri Boyutu: %.2f KB\n", sizeKB)
	fmt.Println(string(jsonData))
}
