# EncryptResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Cipher** | Pointer to **string** | Encrypted data in hexadecimal | [optional] 
**Iv** | Pointer to **string** | Initialization vector in hexadecimal | [optional] 
**Id** | Pointer to **int32** | Id key that used to encrypt plan data | [optional] 
**Ok** | Pointer to **bool** | Indicates if the operation was successful | [optional] 
**TimeCreated** | Pointer to **time.Time** | Creation timestamp in RFC3339 | [optional] 
**TimeUpdated** | Pointer to **time.Time** | Update timestamp in RFC3339 | [optional] 

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

### GetCipher

`func (o *EncryptResponse) GetCipher() string`

GetCipher returns the Cipher field if non-nil, zero value otherwise.

### GetCipherOk

`func (o *EncryptResponse) GetCipherOk() (*string, bool)`

GetCipherOk returns a tuple with the Cipher field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCipher

`func (o *EncryptResponse) SetCipher(v string)`

SetCipher sets Cipher field to given value.

### HasCipher

`func (o *EncryptResponse) HasCipher() bool`

HasCipher returns a boolean if a field has been set.

### GetIv

`func (o *EncryptResponse) GetIv() string`

GetIv returns the Iv field if non-nil, zero value otherwise.

### GetIvOk

`func (o *EncryptResponse) GetIvOk() (*string, bool)`

GetIvOk returns a tuple with the Iv field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIv

`func (o *EncryptResponse) SetIv(v string)`

SetIv sets Iv field to given value.

### HasIv

`func (o *EncryptResponse) HasIv() bool`

HasIv returns a boolean if a field has been set.

### GetId

`func (o *EncryptResponse) GetId() int32`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *EncryptResponse) GetIdOk() (*int32, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *EncryptResponse) SetId(v int32)`

SetId sets Id field to given value.

### HasId

`func (o *EncryptResponse) HasId() bool`

HasId returns a boolean if a field has been set.

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


