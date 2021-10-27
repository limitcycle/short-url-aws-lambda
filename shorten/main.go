package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/teris-io/shortid"
)

type Request struct {
	URL string `json:"url"`
}

type Response struct {
	ShortURL string `json:"short_url"`
}

type Link struct {
	ShortURL string `json:"short_url"`
	LongURL  string `json:"long_url"`
}

type deps struct {
	ddb   dynamodbiface.DynamoDBAPI
	table string
}

func (d *deps) handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	if d.ddb == nil {
		sess := session.Must(session.NewSession())

		ddb := dynamodb.New(sess)
		d.ddb = ddb
		d.table = os.Getenv("TABLE_NAME")
	}
	rb := Request{}
	if err := json.Unmarshal([]byte(request.Body), &rb); err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	// Generate a short URL
	shortURL := shortid.MustGenerate()
	link := &Link{
		ShortURL: shortURL,
		LongURL:  rb.URL,
	}
	// Marshal the link into a attribute value map
	item, err := dynamodbattribute.MarshalMap(link)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	// Insert link into DynamoDB table
	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(d.table),
	}
	if _, err := d.ddb.PutItem(input); err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	// Return short URL in the response body.
	response, err := json.Marshal(Response{ShortURL: shortURL})
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(response),
	}, nil
}

func main() {
	d := deps{}

	lambda.Start(d.handler)
}
