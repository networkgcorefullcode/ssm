# DecryptRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**KeyLabel** | **string** | Label of the key to decrypt | 
**Cipher** | **string** | Encrypted data in Base64 | 
**Iv** | **string** | Initialization vector in Base64 (same used for encryption) | 
**Id** | **int32** | Optional ID for tracking | 
**EncryptionAlgorithm** | **int32** | Encryption algorithm to use (1: AES, 2: AES, 3: DES, 4: DES3) | 

## Methods

### NewDecryptRequest

`func NewDecryptRequest(keyLabel string, cipher string, iv string, id int32, encryptionAlgorithm int32, ) *DecryptRequest`

NewDecryptRequest instantiates a new DecryptRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDecryptRequestWithDefaults

`func NewDecryptRequestWithDefaults() *DecryptRequest`

NewDecryptRequestWithDefaults instantiates a new DecryptRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetKeyLabel

`func (o *DecryptRequest) GetKeyLabel() string`

GetKeyLabel returns the KeyLabel field if non-nil, zero value otherwise.

### GetKeyLabelOk

`func (o *DecryptRequest) GetKeyLabelOk() (*string, bool)`

GetKeyLabelOk returns a tuple with the KeyLabel field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKeyLabel

`func (o *DecryptRequest) SetKeyLabel(v string)`

SetKeyLabel sets KeyLabel field to given value.


### GetCipher

`func (o *DecryptRequest) GetCipher() string`

GetCipher returns the Cipher field if non-nil, zero value otherwise.

### GetCipherOk

`func (o *DecryptRequest) GetCipherOk() (*string, bool)`

GetCipherOk returns a tuple with the Cipher field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCipher

`func (o *DecryptRequest) SetCipher(v string)`

SetCipher sets Cipher field to given value.


### GetIv

`func (o *DecryptRequest) GetIv() string`

GetIv returns the Iv field if non-nil, zero value otherwise.

### GetIvOk

`func (o *DecryptRequest) GetIvOk() (*string, bool)`

GetIvOk returns a tuple with the Iv field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIv

`func (o *DecryptRequest) SetIv(v string)`

SetIv sets Iv field to given value.


### GetId

`func (o *DecryptRequest) GetId() int32`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *DecryptRequest) GetIdOk() (*int32, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *DecryptRequest) SetId(v int32)`

SetId sets Id field to given value.


### GetEncryptionAlgorithm

`func (o *DecryptRequest) GetEncryptionAlgorithm() int32`

GetEncryptionAlgorithm returns the EncryptionAlgorithm field if non-nil, zero value otherwise.

### GetEncryptionAlgorithmOk

`func (o *DecryptRequest) GetEncryptionAlgorithmOk() (*int32, bool)`

GetEncryptionAlgorithmOk returns a tuple with the EncryptionAlgorithm field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEncryptionAlgorithm

`func (o *DecryptRequest) SetEncryptionAlgorithm(v int32)`

SetEncryptionAlgorithm sets EncryptionAlgorithm field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


