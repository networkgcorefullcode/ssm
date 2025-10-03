# EncryptResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CipherB64** | Pointer to **string** | Datos cifrados en Base64 | [optional] 
**IvB64** | Pointer to **string** | Vector de inicialización en Base64 | [optional] 
**Ok** | Pointer to **bool** | Indica si la operación fue exitosa | [optional] 
**TimeCreated** | Pointer to **time.Time** | Timestamp de creación en RFC3339 | [optional] 
**TimeUpdated** | Pointer to **time.Time** | Timestamp de actualización en RFC3339 | [optional] 

## Methods

### NewEncryptResponse

`func NewEncryptResponse() *EncryptResponse`

NewEncryptResponse instantiates a new EncryptResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewEncryptResponseWithDefaults

`func NewEncryptResponseWithDefaults() *EncryptResponse`

NewEncryptResponseWithDefaults instantiates a new EncryptResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCipherB64

`func (o *EncryptResponse) GetCipherB64() string`

GetCipherB64 returns the CipherB64 field if non-nil, zero value otherwise.

### GetCipherB64Ok

`func (o *EncryptResponse) GetCipherB64Ok() (*string, bool)`

GetCipherB64Ok returns a tuple with the CipherB64 field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCipherB64

`func (o *EncryptResponse) SetCipherB64(v string)`

SetCipherB64 sets CipherB64 field to given value.

### HasCipherB64

`func (o *EncryptResponse) HasCipherB64() bool`

HasCipherB64 returns a boolean if a field has been set.

### GetIvB64

`func (o *EncryptResponse) GetIvB64() string`

GetIvB64 returns the IvB64 field if non-nil, zero value otherwise.

### GetIvB64Ok

`func (o *EncryptResponse) GetIvB64Ok() (*string, bool)`

GetIvB64Ok returns a tuple with the IvB64 field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIvB64

`func (o *EncryptResponse) SetIvB64(v string)`

SetIvB64 sets IvB64 field to given value.

### HasIvB64

`func (o *EncryptResponse) HasIvB64() bool`

HasIvB64 returns a boolean if a field has been set.

### GetOk

`func (o *EncryptResponse) GetOk() bool`

GetOk returns the Ok field if non-nil, zero value otherwise.

### GetOkOk

`func (o *EncryptResponse) GetOkOk() (*bool, bool)`

GetOkOk returns a tuple with the Ok field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOk

`func (o *EncryptResponse) SetOk(v bool)`

SetOk sets Ok field to given value.

### HasOk

`func (o *EncryptResponse) HasOk() bool`

HasOk returns a boolean if a field has been set.

### GetTimeCreated

`func (o *EncryptResponse) GetTimeCreated() time.Time`

GetTimeCreated returns the TimeCreated field if non-nil, zero value otherwise.

### GetTimeCreatedOk

`func (o *EncryptResponse) GetTimeCreatedOk() (*time.Time, bool)`

GetTimeCreatedOk returns a tuple with the TimeCreated field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTimeCreated

`func (o *EncryptResponse) SetTimeCreated(v time.Time)`

SetTimeCreated sets TimeCreated field to given value.

### HasTimeCreated

`func (o *EncryptResponse) HasTimeCreated() bool`

HasTimeCreated returns a boolean if a field has been set.

### GetTimeUpdated

`func (o *EncryptResponse) GetTimeUpdated() time.Time`

GetTimeUpdated returns the TimeUpdated field if non-nil, zero value otherwise.

### GetTimeUpdatedOk

`func (o *EncryptResponse) GetTimeUpdatedOk() (*time.Time, bool)`

GetTimeUpdatedOk returns a tuple with the TimeUpdated field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTimeUpdated

`func (o *EncryptResponse) SetTimeUpdated(v time.Time)`

SetTimeUpdated sets TimeUpdated field to given value.

### HasTimeUpdated

`func (o *EncryptResponse) HasTimeUpdated() bool`

HasTimeUpdated returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


