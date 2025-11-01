# GenAESKeyResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Handle** | Pointer to **int32** | HSM key handle | [optional] 
**Id** | Pointer to **int32** | Generated key identifier | [optional] 
**Bits** | Pointer to **int32** | Size of the generated key in bits | [optional] 

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

### GetId

`func (o *GenAESKeyResponse) GetId() int32`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *GenAESKeyResponse) GetIdOk() (*int32, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *GenAESKeyResponse) SetId(v int32)`

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


