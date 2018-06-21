# Running a lambda on localstack

## Dependencies

1. [Rialto-api / localstack](https://github.com/sul-dlss-labs/rialto-api/)

## Create the lambda zip file, upload and subscribe

```
make
```

2. Start localstack. If you're on a Mac, ensure you are running the docker daemon.
```
SERVICES=lambda,sns LAMBDA_EXECUTOR=docker localstack start
```

3. Upload zip and create a function definition
```
AWS_ACCESS_KEY_ID=999999 AWS_SECRET_ACCESS_KEY=1231 aws \
--endpoint-url http://localhost:4574 lambda create-function \
--function-name f1 \
--runtime go1.x \
--role r1 \
--handler main \
--environment "Variables={SOLR_HOST=http://127.0.0.1:8983/solr,SOLR_COLLECTION=collection1}" \
--zip-file fileb://lambda.zip
```

4. Create SNS topic
```
AWS_ACCESS_KEY_ID=999999 AWS_SECRET_ACCESS_KEY=1231 aws sns \
--endpoint-url=http://localhost:4575 create-topic \
--name data-update
```

5. Subscribe to SNS events
```
AWS_ACCESS_KEY_ID=999999 AWS_SECRET_ACCESS_KEY=1231 aws sns \
--endpoint-url=http://localhost:4575 subscribe \
--topic-arn arn:aws:sns:us-east-1:123456789012:data-update \
--protocol lambda \
--notification-endpoint arn:aws:lambda:us-east-1:000000000000:function:f1
```

6. Start Solr and create a collection
```
gem install solr_wrapper
solr_wrapper

```

7. Publish a Message
```
AWS_ACCESS_KEY_ID=999999 AWS_SECRET_ACCESS_KEY=1231 aws sns \
--endpoint-url=http://localhost:4575 publish \
--topic-arn arn:aws:sns:us-east-1:123456789012:data-update \
--message '{"Records": [{"EventSource": "foo", "Sns": { "Timestamp": "2014-05-16T08:28:06.801Z",
"Message": "{ \"foo_si\": \"Hello world!\" }" }}]}'
```

8. View output
When you go to http://127.0.0.1:8983/solr/collection1/select?q=*:*

You should see an item record with:
```
"_source":{"foo": "barfoo"}
```

9. Cleanup (necessary before you upload a newer version of the function)

```
AWS_ACCESS_KEY_ID=999999 AWS_SECRET_ACCESS_KEY=1231 aws lambda \
--endpoint-url=http://localhost:4574 delete-function \
--function-name f1
```
