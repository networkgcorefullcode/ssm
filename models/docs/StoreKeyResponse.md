# StoreKeyResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Handle** | Pointer to **int32** | Handle of the stored key | [optional] 
**CipherKey** | Pointer to **string** | Stored encrypted key | [optional] 

## Methods

### NewStoreKeyResponse

`func NewStoreKeyResponse() *StoreKeyResponse`

NewStoreKeyResponse instantiates a new StoreKeyResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewStoreKeyResponseWithDefaults

`func NewStoreKeyResponseWithDefaults() *StoreKeyResponse`

NewStoreKeyResponseWithDefaults instantiates a new StoreKeyResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetHandle

`func (o *StoreKeyResponse) GetHandle() int32`

GetHandle returns the Handle field if non-nil, zero value otherwise.

### GetHandleOk

`func (o *StoreKeyResponse) GetHandleOk() (*int32, bool)`

GetHandleOk returns a tuple with the Handle field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHandle

`func (o *StoreKeyResponse) SetHandle(v int32)`

SetHandle sets Handle field to given value.

### HasHandle

`func (o *StoreKeyResponse) HasHandle() bool`

HasHandle returns a boolean if a field has been set.

### GetCipherKey

`func (o *StoreKeyResponse) GetCipherKey() string`

GetCipherKey returns the CipherKey field if non-nil, zero value otherwise.

### GetCipherKeyOk

`func (o *StoreKeyResponse) GetCipherKeyOk() (*string, bool)`

GetCipherKeyOk returns a tuple with the CipherKey field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCipherKey

`func (o *StoreKeyResponse) SetCipherKey(v string)`

SetCipherKey sets CipherKey field to given value.

### HasCipherKey

`func (o *StoreKeyResponse) HasCipherKey() bool`

HasCipherKey returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


