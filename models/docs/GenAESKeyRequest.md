# GenAESKeyRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Label** | **string** | Etiqueta única para identificar la clave | 
**Id** | **string** | Identificador único de la clave | 
**Bits** | **int32** | Tamaño de la clave en bits | 

## Methods

### NewGenAESKeyRequest

`func NewGenAESKeyRequest(label string, id string, bits int32, ) *GenAESKeyRequest`

NewGenAESKeyRequest instantiates a new GenAESKeyRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewGenAESKeyRequestWithDefaults

`func NewGenAESKeyRequestWithDefaults() *GenAESKeyRequest`

NewGenAESKeyRequestWithDefaults instantiates a new GenAESKeyRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetLabel

`func (o *GenAESKeyRequest) GetLabel() string`

GetLabel returns the Label field if non-nil, zero value otherwise.

### GetLabelOk

`func (o *GenAESKeyRequest) GetLabelOk() (*string, bool)`

GetLabelOk returns a tuple with the Label field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabel

`func (o *GenAESKeyRequest) SetLabel(v string)`

SetLabel sets Label field to given value.


### GetId

`func (o *GenAESKeyRequest) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *GenAESKeyRequest) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *GenAESKeyRequest) SetId(v string)`

SetId sets Id field to given value.


### GetBits

`func (o *GenAESKeyRequest) GetBits() int32`

GetBits returns the Bits field if non-nil, zero value otherwise.

### GetBitsOk

`func (o *GenAESKeyRequest) GetBitsOk() (*int32, bool)`

GetBitsOk returns a tuple with the Bits field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBits

`func (o *GenAESKeyRequest) SetBits(v int32)`

SetBits sets Bits field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


