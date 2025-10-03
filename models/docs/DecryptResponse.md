# DecryptResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**PlainB64** | Pointer to **string** | Datos descifrados en Base64 | [optional] 

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

### GetPlainB64

`func (o *DecryptResponse) GetPlainB64() string`

GetPlainB64 returns the PlainB64 field if non-nil, zero value otherwise.

### GetPlainB64Ok

`func (o *DecryptResponse) GetPlainB64Ok() (*string, bool)`

GetPlainB64Ok returns a tuple with the PlainB64 field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPlainB64

`func (o *DecryptResponse) SetPlainB64(v string)`

SetPlainB64 sets PlainB64 field to given value.

### HasPlainB64

`func (o *DecryptResponse) HasPlainB64() bool`

HasPlainB64 returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


