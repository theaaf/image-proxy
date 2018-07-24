# aws-sam

This package allows you to run the image proxy using the AWS Serverless Application Model (SAM). 

## Developing

* Use `make build` to compile the package for the Lambda runtime environment.
* Use `sam local start-api` to run the package locally.

## Deploying

* Use `make build` to compile the package for the Lambda runtime environment.
* Package everything up: `aws --profile aaf-platform cloudformation package --template-file ./template.yaml --s3-bucket aaf-platform-image-proxy-packaging --output-template-file ./template-packaged.yaml`
* Deploy: `aws --profile aaf-platform cloudformation deploy --template-file ./template-packaged.yaml --stack-name image-proxy-dev --capabilities CAPABILITY_IAM`
