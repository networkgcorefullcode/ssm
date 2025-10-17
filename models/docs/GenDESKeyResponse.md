# GenDESKeyResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Handle** | Pointer to **int32** | HSM key handle | [optional] 
**Id** | Pointer to **int32** | Key identifier | [optional] 

## Methods

### NewGenDESKeyResponse

`func NewGenDESKeyResponse() *GenDESKeyResponse`

NewGenDESKeyResponse instantiates a new GenDESKeyResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewGenDESKeyResponseWithDefaults

`func NewGenDESKeyResponseWithDefaults() *GenDESKeyResponse`

NewGenDESKeyResponseWithDefaults instantiates a new GenDESKeyResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetHandle

`func (o *GenDESKeyResponse) GetHandle() int32`

GetHandle returns the Handle field if non-nil, zero value otherwise.

### GetHandleOk

`func (o *GenDESKeyResponse) GetHandleOk() (*int32, bool)`

GetHandleOk returns a tuple with the Handle field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHandle

`func (o *GenDESKeyResponse) SetHandle(v int32)`

SetHandle sets Handle field to given value.

### HasHandle

`func (o *GenDESKeyResponse) HasHandle() bool`

HasHandle returns a boolean if a field has been set.

### GetId

`func (o *GenDESKeyResponse) GetId() int32`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *GenDESKeyResponse) GetIdOk() (*int32, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *GenDESKeyResponse) SetId(v int32)`

SetId sets Id field to given value.

### HasId

`func (o *GenDESKeyResponse) HasId() bool`

HasId returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


