# EncryptRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**KeyLabel** | **string** | Label of the key to encrypt | 
**Plain** | **string** | Data to encrypt encoded in hexadecimal | 
**EncryptionAlgorithm** | **int32** | Encryption algorithm to use (1: AES, 2: AES, 3: DES, 4: DES3) | 

## Methods

### NewEncryptRequest

`func NewEncryptRequest(keyLabel string, plain string, encryptionAlgorithm int32, ) *EncryptRequest`

NewEncryptRequest instantiates a new EncryptRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewEncryptRequestWithDefaults

`func NewEncryptRequestWithDefaults() *EncryptRequest`

NewEncryptRequestWithDefaults instantiates a new EncryptRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetKeyLabel

`func (o *EncryptRequest) GetKeyLabel() string`

GetKeyLabel returns the KeyLabel field if non-nil, zero value otherwise.

### GetKeyLabelOk

`func (o *EncryptRequest) GetKeyLabelOk() (*string, bool)`

GetKeyLabelOk returns a tuple with the KeyLabel field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKeyLabel

`func (o *EncryptRequest) SetKeyLabel(v string)`

SetKeyLabel sets KeyLabel field to given value.


### GetPlain

`func (o *EncryptRequest) GetPlain() string`

GetPlain returns the Plain field if non-nil, zero value otherwise.

### GetPlainOk

`func (o *EncryptRequest) GetPlainOk() (*string, bool)`

GetPlainOk returns a tuple with the Plain field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPlain

`func (o *EncryptRequest) SetPlain(v string)`

SetPlain sets Plain field to given value.


### GetEncryptionAlgorithm

`func (o *EncryptRequest) GetEncryptionAlgorithm() int32`

GetEncryptionAlgorithm returns the EncryptionAlgorithm field if non-nil, zero value otherwise.

### GetEncryptionAlgorithmOk

`func (o *EncryptRequest) GetEncryptionAlgorithmOk() (*int32, bool)`

GetEncryptionAlgorithmOk returns a tuple with the EncryptionAlgorithm field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEncryptionAlgorithm

`func (o *EncryptRequest) SetEncryptionAlgorithm(v int32)`

SetEncryptionAlgorithm sets EncryptionAlgorithm field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


