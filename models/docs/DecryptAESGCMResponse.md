# DecryptAESGCMResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Plain** | Pointer to **string** | Decrypted plaintext data in hexadecimal | [optional] 
**Id** | Pointer to **int32** | ID of the key that was used to decrypt the data | [optional] 
**Ok** | Pointer to **bool** | Indicates if the operation was successful | [optional] 
**TimeCreated** | Pointer to **time.Time** | Creation timestamp in RFC3339 format | [optional] 
**TimeUpdated** | Pointer to **time.Time** | Update timestamp in RFC3339 format | [optional] 

## Methods

### NewDecryptAESGCMResponse

`func NewDecryptAESGCMResponse() *DecryptAESGCMResponse`

NewDecryptAESGCMResponse instantiates a new DecryptAESGCMResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDecryptAESGCMResponseWithDefaults

`func NewDecryptAESGCMResponseWithDefaults() *DecryptAESGCMResponse`

NewDecryptAESGCMResponseWithDefaults instantiates a new DecryptAESGCMResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetPlain

`func (o *DecryptAESGCMResponse) GetPlain() string`

GetPlain returns the Plain field if non-nil, zero value otherwise.

### GetPlainOk

`func (o *DecryptAESGCMResponse) GetPlainOk() (*string, bool)`

GetPlainOk returns a tuple with the Plain field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPlain

`func (o *DecryptAESGCMResponse) SetPlain(v string)`

SetPlain sets Plain field to given value.

### HasPlain

`func (o *DecryptAESGCMResponse) HasPlain() bool`

HasPlain returns a boolean if a field has been set.

### GetId

`func (o *DecryptAESGCMResponse) GetId() int32`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *DecryptAESGCMResponse) GetIdOk() (*int32, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *DecryptAESGCMResponse) SetId(v int32)`

SetId sets Id field to given value.

### HasId

`func (o *DecryptAESGCMResponse) HasId() bool`

HasId returns a boolean if a field has been set.

### GetOk

`func (o *DecryptAESGCMResponse) GetOk() bool`

GetOk returns the Ok field if non-nil, zero value otherwise.

### GetOkOk

`func (o *DecryptAESGCMResponse) GetOkOk() (*bool, bool)`

GetOkOk returns a tuple with the Ok field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOk

`func (o *DecryptAESGCMResponse) SetOk(v bool)`

SetOk sets Ok field to given value.

### HasOk

`func (o *DecryptAESGCMResponse) HasOk() bool`

HasOk returns a boolean if a field has been set.

### GetTimeCreated

`func (o *DecryptAESGCMResponse) GetTimeCreated() time.Time`

GetTimeCreated returns the TimeCreated field if non-nil, zero value otherwise.

### GetTimeCreatedOk

`func (o *DecryptAESGCMResponse) GetTimeCreatedOk() (*time.Time, bool)`

GetTimeCreatedOk returns a tuple with the TimeCreated field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTimeCreated

`func (o *DecryptAESGCMResponse) SetTimeCreated(v time.Time)`

SetTimeCreated sets TimeCreated field to given value.

### HasTimeCreated

`func (o *DecryptAESGCMResponse) HasTimeCreated() bool`

HasTimeCreated returns a boolean if a field has been set.

### GetTimeUpdated

`func (o *DecryptAESGCMResponse) GetTimeUpdated() time.Time`

GetTimeUpdated returns the TimeUpdated field if non-nil, zero value otherwise.

### GetTimeUpdatedOk

`func (o *DecryptAESGCMResponse) GetTimeUpdatedOk() (*time.Time, bool)`

GetTimeUpdatedOk returns a tuple with the TimeUpdated field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTimeUpdated

`func (o *DecryptAESGCMResponse) SetTimeUpdated(v time.Time)`

SetTimeUpdated sets TimeUpdated field to given value.

### HasTimeUpdated

`func (o *DecryptAESGCMResponse) HasTimeUpdated() bool`

HasTimeUpdated returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


