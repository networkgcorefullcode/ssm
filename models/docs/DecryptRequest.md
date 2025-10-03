# DecryptRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**KeyLabel** | **string** | Etiqueta de la clave para descifrar | 
**CipherB64** | **string** | Datos cifrados en Base64 | 
**IvB64** | **string** | Vector de inicialización en Base64 (mismo usado para cifrar) | 
**Id** | Pointer to **int32** | ID opcional para tracking | [optional] 

## Methods

### NewDecryptRequest

`func NewDecryptRequest(keyLabel string, cipherB64 string, ivB64 string, ) *DecryptRequest`

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


### GetCipherB64

`func (o *DecryptRequest) GetCipherB64() string`

GetCipherB64 returns the CipherB64 field if non-nil, zero value otherwise.

### GetCipherB64Ok

`func (o *DecryptRequest) GetCipherB64Ok() (*string, bool)`

GetCipherB64Ok returns a tuple with the CipherB64 field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCipherB64

`func (o *DecryptRequest) SetCipherB64(v string)`

SetCipherB64 sets CipherB64 field to given value.


### GetIvB64

`func (o *DecryptRequest) GetIvB64() string`

GetIvB64 returns the IvB64 field if non-nil, zero value otherwise.

### GetIvB64Ok

`func (o *DecryptRequest) GetIvB64Ok() (*string, bool)`

GetIvB64Ok returns a tuple with the IvB64 field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIvB64

`func (o *DecryptRequest) SetIvB64(v string)`

SetIvB64 sets IvB64 field to given value.


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

### HasId

`func (o *DecryptRequest) HasId() bool`

HasId returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


