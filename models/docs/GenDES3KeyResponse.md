# GenDES3KeyResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Handle** | Pointer to **int32** | HSM key handle | [optional] 
**Id** | Pointer to **int32** | Key identifier | [optional] 

## Methods

### NewGenDES3KeyResponse

`func NewGenDES3KeyResponse() *GenDES3KeyResponse`

NewGenDES3KeyResponse instantiates a new GenDES3KeyResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewGenDES3KeyResponseWithDefaults

`func NewGenDES3KeyResponseWithDefaults() *GenDES3KeyResponse`

NewGenDES3KeyResponseWithDefaults instantiates a new GenDES3KeyResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetHandle

`func (o *GenDES3KeyResponse) GetHandle() int32`

GetHandle returns the Handle field if non-nil, zero value otherwise.

### GetHandleOk

`func (o *GenDES3KeyResponse) GetHandleOk() (*int32, bool)`

GetHandleOk returns a tuple with the Handle field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHandle

`func (o *GenDES3KeyResponse) SetHandle(v int32)`

SetHandle sets Handle field to given value.

### HasHandle

`func (o *GenDES3KeyResponse) HasHandle() bool`

HasHandle returns a boolean if a field has been set.

### GetId

`func (o *GenDES3KeyResponse) GetId() int32`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *GenDES3KeyResponse) GetIdOk() (*int32, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *GenDES3KeyResponse) SetId(v int32)`

SetId sets Id field to given value.

### HasId

`func (o *GenDES3KeyResponse) HasId() bool`

HasId returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


