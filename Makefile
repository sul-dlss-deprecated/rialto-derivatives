default: package

package: solr-sqs postgres-sqs solr-sns postgres-sns

solr-sqs:
	GOOS=linux go build -o solr_derivative cmd/solr-sqs/main.go
	zip solr_derivative.zip solr_derivative

postgres-sqs:
	GOOS=linux go build -o postgres_derivative cmd/postgres-sqs/main.go
	zip postgres_derivative.zip postgres_derivative

solr-sns:
	GOOS=linux go build -o solr_derivative_sns cmd/solr-sns/main.go
	zip solr_derivative_sns.zip solr_derivative_sns

postgres-sns:
	GOOS=linux go build -o postgres_derivative_sns cmd/postgres-sns/main.go
	zip postgres_derivative_sns.zip postgres_derivative_sns


local-delete-solr:
	-AWS_ACCESS_KEY_ID=999999 AWS_SECRET_ACCESS_KEY=1231 aws lambda \
	--region us-east-1 \
	--endpoint-url=http://localhost:4574 delete-function \
	--function-name rialto-derivatives-solr-development

local-delete-postgres:
	-AWS_ACCESS_KEY_ID=999999 AWS_SECRET_ACCESS_KEY=1231 aws lambda \
	--region us-east-1 \
	--endpoint-url=http://localhost:4574 delete-function \
	--function-name rialto-derivatives-postgres-development

local-create-solr: solr-sns local-delete-solr
	AWS_ACCESS_KEY_ID=999999 AWS_SECRET_ACCESS_KEY=1231 aws \
	--endpoint-url http://localhost:4574 lambda create-function \
	--function-name rialto-derivatives-solr-development \
	--runtime go1.x \
	--role r1 \
	--handler solr_derivative_sns \
	--region us-east-1 \
	--environment "Variables={SOLR_HOST=http://solr:8983/solr,SOLR_COLLECTION=collection1,\
	  SPARQL_ENDPOINT=http://triplestore:9999/blazegraph/namespace/kb/sparql, \
		SPARQL_RETRIES=300}" \
	--zip-file fileb://solr_derivative_sns.zip

local-create-postgres: postgres-sns local-delete-postgres
	AWS_ACCESS_KEY_ID=999999 AWS_SECRET_ACCESS_KEY=1231 aws \
	--endpoint-url http://localhost:4574 lambda create-function \
	--function-name rialto-derivatives-postgres-development \
	--runtime go1.x \
	--role r1 \
	--handler postgres_derivative_sns \
	--region us-east-1 \
	--environment "Variables={SPARQL_ENDPOINT=http://triplestore:9999/blazegraph/namespace/kb/sparql, \
	  SPARQL_RETRIES=300, \
	  RDS_DB_NAME=rialto_development, \
	  RDS_USERNAME=postgres, \
	  RDS_HOSTNAME=db, \
	  RDS_PORT=5432, \
	  RDS_SSL=false, \
	  RDS_PASSWORD=sekret}" \
	--zip-file fileb://postgres_derivative_sns.zip

local-create-topic:
	AWS_ACCESS_KEY_ID=999999 AWS_SECRET_ACCESS_KEY=1231 aws sns \
	--endpoint-url=http://localhost:4575 create-topic \
	--region us-east-1 \
	--name data-update

local-subscribe-solr: local-create-topic
	AWS_ACCESS_KEY_ID=999999 AWS_SECRET_ACCESS_KEY=1231 aws sns \
	--endpoint-url=http://localhost:4575 subscribe \
	--topic-arn arn:aws:sns:us-east-1:123456789012:data-update \
	--protocol lambda \
	--region us-east-1 \
	--notification-endpoint arn:aws:lambda:us-east-1:000000000000:function:rialto-derivatives-solr-development

local-subscribe-postgres: local-create-topic
	AWS_ACCESS_KEY_ID=999999 AWS_SECRET_ACCESS_KEY=1231 aws sns \
	--endpoint-url=http://localhost:4575 subscribe \
	--topic-arn arn:aws:sns:us-east-1:123456789012:data-update \
	--protocol lambda \
	--region us-east-1 \
	--notification-endpoint arn:aws:lambda:us-east-1:000000000000:function:rialto-derivatives-postgres-development

local-deploy-solr: local-create-solr local-subscribe-solr

local-deploy-postgres: local-create-postgres local-subscribe-postgres

local-deploy: local-deploy-solr local-deploy-postgres
