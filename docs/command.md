# Generate Go structs from OpenAPI spec

```bash
java -jar openapi-generator-cli-7.6.0.jar generate -i handlers/api/ssm-openapi.yml -g go -o ./go-structs-only -t handlers/api/templates/ --package-name=models --additional-properties=packageName=models

sudo apt install golang-golang-x-tools
goimports -w .

cd go-structs-only

rm -rf go.mod go.sum

cd test
for file in *.go; do
    sed -i 's|openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"|openapiclient "github.com/networkgcorefullcode/ssm/models"|g' "$file"
done
cd ..

go mod tidy

cd ../

rm -rf models

mkdir -p models
mv go-structs-only/* models/
rm -rf go-structs-only

go mod tidy
```

```bash
sudo env GODEBUG=tls13=1,tlsdebug=1 go run ./ssm.go --cfg factory/ssmConfig.yml
```
