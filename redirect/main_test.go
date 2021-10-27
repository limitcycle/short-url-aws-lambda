package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type mockGetItem struct {
	dynamodbiface.DynamoDBAPI
	Response dynamodb.GetItemOutput
}

func (d mockGetItem) GetItem(in *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	return &d.Response, nil
}

func TestHandler(t *testing.T) {
	t.Run("Successful Request", func(t *testing.T) {
		m := mockGetItem{
			Response: dynamodb.GetItemOutput{},
		}

		d := deps{
			ddb:   m,
			table: "test_table",
		}

		mockReq := events.APIGatewayProxyRequest{
			Path: "mock_path",
		}

		_, err := d.Handler(mockReq)
		if err != nil {
			t.Fatal(err)
		}
	})
}
