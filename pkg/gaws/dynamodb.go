package gaws

import (
	"dynamodb-protobuf-performance-test/pkg/settings"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var db dynamodbiface.DynamoDBAPI

func GetDb() dynamodbiface.DynamoDBAPI {
	return db
}

func Setup(config *settings.AwsDynamoDb) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(config.Region),
	}))
	svc := dynamodb.New(sess)
	_, err := svc.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String(config.TableName),
	})

	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	db = svc
}
