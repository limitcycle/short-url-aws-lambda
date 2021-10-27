package main

import (
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type Link struct {
	ShortURL string `json:"short_url"`
	LongURL  string `json:"long_url"`
}

type deps struct {
	ddb   dynamodbiface.DynamoDBAPI
	table string
}

func (d *deps) Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	shortURL := request.PathParameters["short_url"]

	// Read link item
	result, err := d.ddb.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(d.table),
		Key: map[string]*dynamodb.AttributeValue{
			"short_url": {
				S: aws.String(shortURL),
			},
		},
	})
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	// Parse link item into the Link struct
	link := Link{}
	if err := dynamodbattribute.UnmarshalMap(result.Item, &link); err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	// Redirect user to the long URL by specifying the location header.
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusPermanentRedirect,
		Headers: map[string]string{
			"location": link.LongURL,
		},
	}, nil
}

func main() {
	sess := session.Must(session.NewSession())
	ddb := dynamodb.New(sess)

	d := deps{
		ddb:   ddb,
		table: os.Getenv("TABLE_NAME"),
	}

	lambda.Start(d.Handler)
}
