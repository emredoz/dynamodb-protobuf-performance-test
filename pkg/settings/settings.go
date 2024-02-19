package settings

import (
	"os"

	"github.com/joho/godotenv"
)

type AwsDynamoDb struct {
	Region    string
	TableName string
}

var AwsDynamoDbSettings = &AwsDynamoDb{}

func Setup() {
	_ = godotenv.Load()
	AwsDynamoDbSettings.TableName = os.Getenv("DYNAMO_TABLE_NAME")
	AwsDynamoDbSettings.Region = os.Getenv("AWS_REGION")
}
