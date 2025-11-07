# DecryptAESGCMRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**KeyLabel** | **string** | Label of the AES key to use for GCM decryption | 
**Cipher** | **string** | Ciphertext data to decrypt encoded in hexadecimal | 
**Iv** | **string** | Initialization vector (nonce) in hexadecimal used during encryption (12 bytes recommended) | 
**Tag** | **string** | Authentication tag in hexadecimal (16 bytes for 128-bit tag) | 
**Aad** | Pointer to **string** | Additional Authenticated Data (AAD) in hexadecimal (must match the AAD used during encryption) | [optional] 

## Methods

### NewDecryptAESGCMRequest

`func NewDecryptAESGCMRequest(keyLabel string, cipher string, iv string, tag string, ) *DecryptAESGCMRequest`

NewDecryptAESGCMRequest instantiates a new DecryptAESGCMRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDecryptAESGCMRequestWithDefaults

`func NewDecryptAESGCMRequestWithDefaults() *DecryptAESGCMRequest`

NewDecryptAESGCMRequestWithDefaults instantiates a new DecryptAESGCMRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetKeyLabel

`func (o *DecryptAESGCMRequest) GetKeyLabel() string`

GetKeyLabel returns the KeyLabel field if non-nil, zero value otherwise.

### GetKeyLabelOk

`func (o *DecryptAESGCMRequest) GetKeyLabelOk() (*string, bool)`

GetKeyLabelOk returns a tuple with the KeyLabel field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKeyLabel

`func (o *DecryptAESGCMRequest) SetKeyLabel(v string)`

SetKeyLabel sets KeyLabel field to given value.


### GetCipher

`func (o *DecryptAESGCMRequest) GetCipher() string`

GetCipher returns the Cipher field if non-nil, zero value otherwise.

### GetCipherOk

`func (o *DecryptAESGCMRequest) GetCipherOk() (*string, bool)`

GetCipherOk returns a tuple with the Cipher field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCipher

`func (o *DecryptAESGCMRequest) SetCipher(v string)`

SetCipher sets Cipher field to given value.


### GetIv

`func (o *DecryptAESGCMRequest) GetIv() string`

GetIv returns the Iv field if non-nil, zero value otherwise.

### GetIvOk

`func (o *DecryptAESGCMRequest) GetIvOk() (*string, bool)`

GetIvOk returns a tuple with the Iv field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIv

`func (o *DecryptAESGCMRequest) SetIv(v string)`

SetIv sets Iv field to given value.


### GetTag

`func (o *DecryptAESGCMRequest) GetTag() string`

GetTag returns the Tag field if non-nil, zero value otherwise.

### GetTagOk

`func (o *DecryptAESGCMRequest) GetTagOk() (*string, bool)`

GetTagOk returns a tuple with the Tag field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTag

`func (o *DecryptAESGCMRequest) SetTag(v string)`

SetTag sets Tag field to given value.


### GetAad

`func (o *DecryptAESGCMRequest) GetAad() string`

GetAad returns the Aad field if non-nil, zero value otherwise.

### GetAadOk

`func (o *DecryptAESGCMRequest) GetAadOk() (*string, bool)`

GetAadOk returns a tuple with the Aad field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAad

`func (o *DecryptAESGCMRequest) SetAad(v string)`

SetAad sets Aad field to given value.

### HasAad

`func (o *DecryptAESGCMRequest) HasAad() bool`

HasAad returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


