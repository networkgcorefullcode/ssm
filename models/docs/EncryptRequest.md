# EncryptRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**KeyLabel** | **string** | Etiqueta de la clave para cifrar | 
**PlainB64** | **string** | Datos a cifrar codificados en Base64 | 

## Methods

### NewEncryptRequest

`func NewEncryptRequest(keyLabel string, plainB64 string, ) *EncryptRequest`

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


### GetPlainB64

`func (o *EncryptRequest) GetPlainB64() string`

GetPlainB64 returns the PlainB64 field if non-nil, zero value otherwise.

### GetPlainB64Ok

`func (o *EncryptRequest) GetPlainB64Ok() (*string, bool)`

GetPlainB64Ok returns a tuple with the PlainB64 field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPlainB64

`func (o *EncryptRequest) SetPlainB64(v string)`

SetPlainB64 sets PlainB64 field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


