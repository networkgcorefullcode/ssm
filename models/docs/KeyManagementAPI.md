# \KeyManagementAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**DeleteKey**](KeyManagementAPI.md#DeleteKey) | **Delete** /store-key | Delete key
[**GenerateAESKey**](KeyManagementAPI.md#GenerateAESKey) | **Post** /generate-aes-key | Generate new AES key
[**GenerateDES3Key**](KeyManagementAPI.md#GenerateDES3Key) | **Post** /generate-des3-key | Generate new DES3 key
[**GenerateDESKey**](KeyManagementAPI.md#GenerateDESKey) | **Post** /generate-des-key | Generate new DES key
[**StoreKey**](KeyManagementAPI.md#StoreKey) | **Post** /store-key | Store existing key
[**UpdateKey**](KeyManagementAPI.md#UpdateKey) | **Put** /store-key | Update key



## DeleteKey

> DeleteKeyResponse DeleteKey(ctx).DeleteKeyRequest(deleteKeyRequest).Execute()

Delete key



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
	deleteKeyRequest := *openapiclient.NewDeleteKeyRequest("ImportedKey", int32(1)) // DeleteKeyRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.KeyManagementAPI.DeleteKey(context.Background()).DeleteKeyRequest(deleteKeyRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `KeyManagementAPI.DeleteKey``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DeleteKey`: DeleteKeyResponse
	fmt.Fprintf(os.Stdout, "Response from `KeyManagementAPI.DeleteKey`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDeleteKeyRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **deleteKeyRequest** | [**DeleteKeyRequest**](DeleteKeyRequest.md) |  | 

### Return type

[**DeleteKeyResponse**](DeleteKeyResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GenerateAESKey

> GenAESKeyResponse GenerateAESKey(ctx).GenAESKeyRequest(genAESKeyRequest).Execute()

Generate new AES key



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
	genAESKeyRequest := *openapiclient.NewGenAESKeyRequest(int32(123), int32(256)) // GenAESKeyRequest | 

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


## GenerateDES3Key

> GenDES3KeyResponse GenerateDES3Key(ctx).GenDES3KeyRequest(genDES3KeyRequest).Execute()

Generate new DES3 key



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
	genDES3KeyRequest := *openapiclient.NewGenDES3KeyRequest(int32(1)) // GenDES3KeyRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.KeyManagementAPI.GenerateDES3Key(context.Background()).GenDES3KeyRequest(genDES3KeyRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `KeyManagementAPI.GenerateDES3Key``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GenerateDES3Key`: GenDES3KeyResponse
	fmt.Fprintf(os.Stdout, "Response from `KeyManagementAPI.GenerateDES3Key`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiGenerateDES3KeyRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **genDES3KeyRequest** | [**GenDES3KeyRequest**](GenDES3KeyRequest.md) |  | 

### Return type

[**GenDES3KeyResponse**](GenDES3KeyResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GenerateDESKey

> GenDESKeyResponse GenerateDESKey(ctx).GenDESKeyRequest(genDESKeyRequest).Execute()

Generate new DES key



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
	genDESKeyRequest := *openapiclient.NewGenDESKeyRequest(int32(1)) // GenDESKeyRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.KeyManagementAPI.GenerateDESKey(context.Background()).GenDESKeyRequest(genDESKeyRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `KeyManagementAPI.GenerateDESKey``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `GenerateDESKey`: GenDESKeyResponse
	fmt.Fprintf(os.Stdout, "Response from `KeyManagementAPI.GenerateDESKey`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiGenerateDESKeyRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **genDESKeyRequest** | [**GenDESKeyRequest**](GenDESKeyRequest.md) |  | 

### Return type

[**GenDESKeyResponse**](GenDESKeyResponse.md)

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

Store existing key



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
	storeKeyRequest := *openapiclient.NewStoreKeyRequest("ImportedKey", int32(2), string([B@3051e0b2), "AES") // StoreKeyRequest | 

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


## UpdateKey

> UpdateKeyResponse UpdateKey(ctx).UpdateKeyRequest(updateKeyRequest).Execute()

Update key



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
	updateKeyRequest := *openapiclient.NewUpdateKeyRequest("ImportedKey", int32(2), "bmV3X2tleV92YWx1ZQ==", "AES") // UpdateKeyRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.KeyManagementAPI.UpdateKey(context.Background()).UpdateKeyRequest(updateKeyRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `KeyManagementAPI.UpdateKey``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UpdateKey`: UpdateKeyResponse
	fmt.Fprintf(os.Stdout, "Response from `KeyManagementAPI.UpdateKey`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiUpdateKeyRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **updateKeyRequest** | [**UpdateKeyRequest**](UpdateKeyRequest.md) |  | 

### Return type

[**UpdateKeyResponse**](UpdateKeyResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

