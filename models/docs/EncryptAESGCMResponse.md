# EncryptAESGCMResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Cipher** | Pointer to **string** | Encrypted data (ciphertext) in hexadecimal | [optional] 
**Iv** | Pointer to **string** | Initialization vector (nonce) in hexadecimal used for GCM encryption (12 bytes recommended) | [optional] 
**Tag** | Pointer to **string** | Authentication tag in hexadecimal (16 bytes for 128-bit tag) | [optional] 
**Id** | Pointer to **int32** | ID of the key that was used to encrypt the data | [optional] 
**Ok** | Pointer to **bool** | Indicates if the operation was successful | [optional] 
**TimeCreated** | Pointer to **time.Time** | Creation timestamp in RFC3339 format | [optional] 
**TimeUpdated** | Pointer to **time.Time** | Update timestamp in RFC3339 format | [optional] 

## Methods

### NewEncryptAESGCMResponse

`func NewEncryptAESGCMResponse() *EncryptAESGCMResponse`

NewEncryptAESGCMResponse instantiates a new EncryptAESGCMResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewEncryptAESGCMResponseWithDefaults

`func NewEncryptAESGCMResponseWithDefaults() *EncryptAESGCMResponse`

NewEncryptAESGCMResponseWithDefaults instantiates a new EncryptAESGCMResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCipher

`func (o *EncryptAESGCMResponse) GetCipher() string`

GetCipher returns the Cipher field if non-nil, zero value otherwise.

### GetCipherOk

`func (o *EncryptAESGCMResponse) GetCipherOk() (*string, bool)`

GetCipherOk returns a tuple with the Cipher field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCipher

`func (o *EncryptAESGCMResponse) SetCipher(v string)`

SetCipher sets Cipher field to given value.

### HasCipher

`func (o *EncryptAESGCMResponse) HasCipher() bool`

HasCipher returns a boolean if a field has been set.

### GetIv

`func (o *EncryptAESGCMResponse) GetIv() string`

GetIv returns the Iv field if non-nil, zero value otherwise.

### GetIvOk

`func (o *EncryptAESGCMResponse) GetIvOk() (*string, bool)`

GetIvOk returns a tuple with the Iv field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIv

`func (o *EncryptAESGCMResponse) SetIv(v string)`

SetIv sets Iv field to given value.

### HasIv

`func (o *EncryptAESGCMResponse) HasIv() bool`

HasIv returns a boolean if a field has been set.

### GetTag

`func (o *EncryptAESGCMResponse) GetTag() string`

GetTag returns the Tag field if non-nil, zero value otherwise.

### GetTagOk

`func (o *EncryptAESGCMResponse) GetTagOk() (*string, bool)`

GetTagOk returns a tuple with the Tag field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTag

`func (o *EncryptAESGCMResponse) SetTag(v string)`

SetTag sets Tag field to given value.

### HasTag

`func (o *EncryptAESGCMResponse) HasTag() bool`

HasTag returns a boolean if a field has been set.

### GetId

`func (o *EncryptAESGCMResponse) GetId() int32`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *EncryptAESGCMResponse) GetIdOk() (*int32, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *EncryptAESGCMResponse) SetId(v int32)`

SetId sets Id field to given value.

### HasId

`func (o *EncryptAESGCMResponse) HasId() bool`

HasId returns a boolean if a field has been set.

### GetOk

`func (o *EncryptAESGCMResponse) GetOk() bool`

GetOk returns the Ok field if non-nil, zero value otherwise.

### GetOkOk

`func (o *EncryptAESGCMResponse) GetOkOk() (*bool, bool)`

GetOkOk returns a tuple with the Ok field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOk

`func (o *EncryptAESGCMResponse) SetOk(v bool)`

SetOk sets Ok field to given value.

### HasOk

`func (o *EncryptAESGCMResponse) HasOk() bool`

HasOk returns a boolean if a field has been set.

### GetTimeCreated

`func (o *EncryptAESGCMResponse) GetTimeCreated() time.Time`

GetTimeCreated returns the TimeCreated field if non-nil, zero value otherwise.

### GetTimeCreatedOk

`func (o *EncryptAESGCMResponse) GetTimeCreatedOk() (*time.Time, bool)`

GetTimeCreatedOk returns a tuple with the TimeCreated field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTimeCreated

`func (o *EncryptAESGCMResponse) SetTimeCreated(v time.Time)`

SetTimeCreated sets TimeCreated field to given value.

### HasTimeCreated

`func (o *EncryptAESGCMResponse) HasTimeCreated() bool`

HasTimeCreated returns a boolean if a field has been set.

### GetTimeUpdated

`func (o *EncryptAESGCMResponse) GetTimeUpdated() time.Time`

GetTimeUpdated returns the TimeUpdated field if non-nil, zero value otherwise.

### GetTimeUpdatedOk

`func (o *EncryptAESGCMResponse) GetTimeUpdatedOk() (*time.Time, bool)`

GetTimeUpdatedOk returns a tuple with the TimeUpdated field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTimeUpdated

`func (o *EncryptAESGCMResponse) SetTimeUpdated(v time.Time)`

SetTimeUpdated sets TimeUpdated field to given value.

### HasTimeUpdated

`func (o *EncryptAESGCMResponse) HasTimeUpdated() bool`

HasTimeUpdated returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


