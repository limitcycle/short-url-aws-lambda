# short-url-aws-lambda

This is a sample template for short-url-aws-lambda - Below is a brief explanation of what we have generated for you:

```bash
.
├── Makefile                    <-- Make to automate build
├── README.md                   <-- This instructions file
├── shorten                     <-- Source code for a lambda function
│   ├── main.go                 <-- Lambda function code
│   └── main_test.go            <-- Unit tests
│   └── Dockerfile              <-- Dockerfile
├── redirect                    <-- Source code for a lambda function
│   ├── main.go                 <-- Lambda function code
│   └── main_test.go            <-- Unit tests
│   └── Dockerfile              <-- Dockerfile
└── template.yaml
```

## Requirements

* AWS CLI already configured with Administrator permission
* [Docker installed](https://www.docker.com/community-edition)
* SAM CLI - [Install the SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html)

You may need the following for local testing.
* [Golang](https://golang.org)

## Building and deployment

### Create two new ECR repositories

```bash
aws ecr create-repository --repository-name lambda-shorten \
--image-tag-mutability IMMUTABLE --image-scanning-configuration scanOnPush=true
```

```bash
aws ecr create-repository --repository-name lambda-redirect \
--image-tag-mutability IMMUTABLE --image-scanning-configuration scanOnPush=true
```

### Create a new S3 bucket

```bash
aws s3 mb s3://${your_bucket_name}
```

### SAM Build

```bash
sam build
```

### SAM Deploy

1. Use AWS CLI with Parameters

```bash
sam deploy \
--stack-name ${your_stack_name} \
--s3-bucket ${your_bucket_name} \
--s3-prefix ${prefix_text} \
--region ${lambda_region} \
--confirm-changeset \
capabilities CAPABILITY_IAM \
--parameter-overrides TableName=${dynamodb_table_name} \
--image-repositories ShortenFunction=${your_aws_account_id}.dkr.ecr.${lambda_region}.amazonaws.com/lambda-shorten \
--image-repositories RedirectFunction=${your_aws_account_id}.dkr.ecr.${lambda_region}.amazonaws.com/lambda-redirect
```

2. Use AWS CLI with `samconfig.toml`

samconfig.toml

```toml
version = 0.1
[default]
[default.deploy]
[default.deploy.parameters]
stack_name = "${your_stack_name}"
s3_bucket = "${your_bucket_name}"
s3_prefix = "${prefix_text}"
region = "${lambda_region}"
confirm_changeset = true
capabilities = "CAPABILITY_IAM"
parameter_overrides = "TableName=${dynamodb_table_name}"
image_repositories = [
	"ShortenFunction=${your_aws_account_id}.dkr.ecr.${lambda_region}.amazonaws.com/sample-shorten",
	"RedirectFunction=${your_aws_account_id}.dkr.ecr.${lambda_region}.amazonaws.com/sample-redirect"
]
```

sam deploy
```bash
sam deploy
```

### Testing

1. Test `ShortenFunction`

**Request**

```shell
curl -X POST -H 'Content-Type: application/json' -d '{"url": "${test_url}"}' \
https://${random_string}.execute-api.${lambda_region}.amazonaws.com/Prod/shorten
```

**Response**

```json
{"short_url": "HAlqGWK7R"}
```

2. Test `RedirectFunction`

Open the browser and enter the Short URL

```bash
https://${random_string}.execute-api.${lambda_region}.amazonaws.com/Prod/HAlqGWK7R
```

## Destory Resources

```bash
sam delete --stack-name ${your_stack_name}
```
