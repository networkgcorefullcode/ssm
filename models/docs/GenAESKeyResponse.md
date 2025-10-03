# GenAESKeyResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Handle** | Pointer to **int32** | Handle de la clave en el HSM | [optional] 
**Label** | Pointer to **string** | Etiqueta de la clave generada | [optional] 
**Id** | Pointer to **string** | ID de la clave generada | [optional] 
**Bits** | Pointer to **int32** | Tamaño de la clave generada | [optional] 

## Methods

### NewGenAESKeyResponse

`func NewGenAESKeyResponse() *GenAESKeyResponse`

NewGenAESKeyResponse instantiates a new GenAESKeyResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewGenAESKeyResponseWithDefaults

`func NewGenAESKeyResponseWithDefaults() *GenAESKeyResponse`

NewGenAESKeyResponseWithDefaults instantiates a new GenAESKeyResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetHandle

`func (o *GenAESKeyResponse) GetHandle() int32`

GetHandle returns the Handle field if non-nil, zero value otherwise.

### GetHandleOk

`func (o *GenAESKeyResponse) GetHandleOk() (*int32, bool)`

GetHandleOk returns a tuple with the Handle field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHandle

`func (o *GenAESKeyResponse) SetHandle(v int32)`

SetHandle sets Handle field to given value.

### HasHandle

`func (o *GenAESKeyResponse) HasHandle() bool`

HasHandle returns a boolean if a field has been set.

### GetLabel

`func (o *GenAESKeyResponse) GetLabel() string`

GetLabel returns the Label field if non-nil, zero value otherwise.

### GetLabelOk

`func (o *GenAESKeyResponse) GetLabelOk() (*string, bool)`

GetLabelOk returns a tuple with the Label field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabel

`func (o *GenAESKeyResponse) SetLabel(v string)`

SetLabel sets Label field to given value.

### HasLabel

`func (o *GenAESKeyResponse) HasLabel() bool`

HasLabel returns a boolean if a field has been set.

### GetId

`func (o *GenAESKeyResponse) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *GenAESKeyResponse) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *GenAESKeyResponse) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *GenAESKeyResponse) HasId() bool`

HasId returns a boolean if a field has been set.

### GetBits

`func (o *GenAESKeyResponse) GetBits() int32`

GetBits returns the Bits field if non-nil, zero value otherwise.

### GetBitsOk

`func (o *GenAESKeyResponse) GetBitsOk() (*int32, bool)`

GetBitsOk returns a tuple with the Bits field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBits

`func (o *GenAESKeyResponse) SetBits(v int32)`

SetBits sets Bits field to given value.

### HasBits

`func (o *GenAESKeyResponse) HasBits() bool`

HasBits returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


