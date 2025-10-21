# \EncryptionAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**DecryptData**](EncryptionAPI.md#DecryptData) | **Post** /decrypt | Decrypt data
[**EncryptData**](EncryptionAPI.md#EncryptData) | **Post** /encrypt | Encrypt data



## DecryptData

> DecryptResponse DecryptData(ctx).DecryptRequest(decryptRequest).Execute()

Decrypt data



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	decryptRequest := *openapiclient.NewDecryptRequest("MySecretKey", string([B@17271176), string([B@2e3cdec2), int32(12345), int32(1)) // DecryptRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.EncryptionAPI.DecryptData(context.Background()).DecryptRequest(decryptRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `EncryptionAPI.DecryptData``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DecryptData`: DecryptResponse
	fmt.Fprintf(os.Stdout, "Response from `EncryptionAPI.DecryptData`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDecryptDataRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **decryptRequest** | [**DecryptRequest**](DecryptRequest.md) |  | 

### Return type

[**DecryptResponse**](DecryptResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## EncryptData

> EncryptResponse EncryptData(ctx).EncryptRequest(encryptRequest).Execute()

Encrypt data



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	encryptRequest := *openapiclient.NewEncryptRequest("MySecretKey", string([B@451f87af), int32(1)) // EncryptRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.EncryptionAPI.EncryptData(context.Background()).EncryptRequest(encryptRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `EncryptionAPI.EncryptData``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `EncryptData`: EncryptResponse
	fmt.Fprintf(os.Stdout, "Response from `EncryptionAPI.EncryptData`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiEncryptDataRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **encryptRequest** | [**EncryptRequest**](EncryptRequest.md) |  | 

### Return type

[**EncryptResponse**](EncryptResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

