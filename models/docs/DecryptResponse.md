# DecryptResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Plain** | Pointer to **string** | Decrypted data in Base64 | [optional] 

## Methods

### NewDecryptResponse

`func NewDecryptResponse() *DecryptResponse`

NewDecryptResponse instantiates a new DecryptResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDecryptResponseWithDefaults

`func NewDecryptResponseWithDefaults() *DecryptResponse`

NewDecryptResponseWithDefaults instantiates a new DecryptResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetPlain

`func (o *DecryptResponse) GetPlain() string`

GetPlain returns the Plain field if non-nil, zero value otherwise.

### GetPlainOk

`func (o *DecryptResponse) GetPlainOk() (*string, bool)`

GetPlainOk returns a tuple with the Plain field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPlain

`func (o *DecryptResponse) SetPlain(v string)`

SetPlain sets Plain field to given value.

### HasPlain

`func (o *DecryptResponse) HasPlain() bool`

HasPlain returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


