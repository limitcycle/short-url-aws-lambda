package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type mockedPutItem struct {
	dynamodbiface.DynamoDBAPI
	Response dynamodb.PutItemOutput
}

func (d mockedPutItem) PutItem(in *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return &d.Response, nil
}

func TestHandler(t *testing.T) {

	t.Run("Successful Request", func(t *testing.T) {
		m := mockedPutItem{
			Response: dynamodb.PutItemOutput{},
		}

		d := deps{
			ddb:   m,
			table: "test_table",
		}

		mockRequest := events.APIGatewayProxyRequest{
			Body: `{"URL": "mockURL"}`,
		}
		_, err := d.handler(mockRequest)
		if err != nil {
			t.Fatal(err)
		}
	})
}
