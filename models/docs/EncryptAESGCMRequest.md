# EncryptAESGCMRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**KeyLabel** | **string** | Label of the AES key to use for GCM encryption | 
**Plain** | **string** | Plaintext data to encrypt encoded in hexadecimal | 
**Aad** | Pointer to **string** | Additional Authenticated Data (AAD) in hexadecimal (optional). This data is authenticated but not encrypted. | [optional] 

## Methods

### NewEncryptAESGCMRequest

`func NewEncryptAESGCMRequest(keyLabel string, plain string, ) *EncryptAESGCMRequest`

NewEncryptAESGCMRequest instantiates a new EncryptAESGCMRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewEncryptAESGCMRequestWithDefaults

`func NewEncryptAESGCMRequestWithDefaults() *EncryptAESGCMRequest`

NewEncryptAESGCMRequestWithDefaults instantiates a new EncryptAESGCMRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetKeyLabel

`func (o *EncryptAESGCMRequest) GetKeyLabel() string`

GetKeyLabel returns the KeyLabel field if non-nil, zero value otherwise.

### GetKeyLabelOk

`func (o *EncryptAESGCMRequest) GetKeyLabelOk() (*string, bool)`

GetKeyLabelOk returns a tuple with the KeyLabel field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKeyLabel

`func (o *EncryptAESGCMRequest) SetKeyLabel(v string)`

SetKeyLabel sets KeyLabel field to given value.


### GetPlain

`func (o *EncryptAESGCMRequest) GetPlain() string`

GetPlain returns the Plain field if non-nil, zero value otherwise.

### GetPlainOk

`func (o *EncryptAESGCMRequest) GetPlainOk() (*string, bool)`

GetPlainOk returns a tuple with the Plain field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPlain

`func (o *EncryptAESGCMRequest) SetPlain(v string)`

SetPlain sets Plain field to given value.


### GetAad

`func (o *EncryptAESGCMRequest) GetAad() string`

GetAad returns the Aad field if non-nil, zero value otherwise.

### GetAadOk

`func (o *EncryptAESGCMRequest) GetAadOk() (*string, bool)`

GetAadOk returns a tuple with the Aad field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAad

`func (o *EncryptAESGCMRequest) SetAad(v string)`

SetAad sets Aad field to given value.

### HasAad

`func (o *EncryptAESGCMRequest) HasAad() bool`

HasAad returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


