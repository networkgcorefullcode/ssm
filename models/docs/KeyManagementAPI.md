# \KeyManagementAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GenerateAESKey**](KeyManagementAPI.md#GenerateAESKey) | **Post** /generate-aes-key | Generar nueva clave AES
[**StoreKey**](KeyManagementAPI.md#StoreKey) | **Post** /store-key | Almacenar clave existente



## GenerateAESKey

> GenAESKeyResponse GenerateAESKey(ctx).GenAESKeyRequest(genAESKeyRequest).Execute()

Generar nueva clave AES



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
	genAESKeyRequest := *openapiclient.NewGenAESKeyRequest("MySecretKey", "key001", int32(256)) // GenAESKeyRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.KeyManagementAPI.GenerateAESKey(context.Background()).GenAESKeyRequest(genAESKeyRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `KeyManagementAPI.GenerateAESKey``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GenerateAESKey`: GenAESKeyResponse
	fmt.Fprintf(os.Stdout, "Response from `KeyManagementAPI.GenerateAESKey`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiGenerateAESKeyRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **genAESKeyRequest** | [**GenAESKeyRequest**](GenAESKeyRequest.md) |  | 

### Return type

[**GenAESKeyResponse**](GenAESKeyResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## StoreKey

> StoreKeyResponse StoreKey(ctx).StoreKeyRequest(storeKeyRequest).Execute()

Almacenar clave existente



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
	storeKeyRequest := *openapiclient.NewStoreKeyRequest("ImportedKey", "imported001", string([B@1e411d81)) // StoreKeyRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.KeyManagementAPI.StoreKey(context.Background()).StoreKeyRequest(storeKeyRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `KeyManagementAPI.StoreKey``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `StoreKey`: StoreKeyResponse
	fmt.Fprintf(os.Stdout, "Response from `KeyManagementAPI.StoreKey`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiStoreKeyRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **storeKeyRequest** | [**StoreKeyRequest**](StoreKeyRequest.md) |  | 

### Return type

[**StoreKeyResponse**](StoreKeyResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

