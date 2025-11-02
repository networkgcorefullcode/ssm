
```bash
java -jar openapi-generator-cli-7.6.0.jar generate -i handlers/api/ssm-openapi.yml -g go -o ./go-structs-only -t handlers/api/templates/ --package-name=models --additional-properties=packageName=models
```