# DeleteKeyRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**KeyLabel** | **string** | Label of the key to delete | 
**Id** | **int32** | Key identifier (optional) | 

## Methods

### NewDeleteKeyRequest

`func NewDeleteKeyRequest(keyLabel string, id int32, ) *DeleteKeyRequest`

NewDeleteKeyRequest instantiates a new DeleteKeyRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDeleteKeyRequestWithDefaults

`func NewDeleteKeyRequestWithDefaults() *DeleteKeyRequest`

NewDeleteKeyRequestWithDefaults instantiates a new DeleteKeyRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetKeyLabel

`func (o *DeleteKeyRequest) GetKeyLabel() string`

GetKeyLabel returns the KeyLabel field if non-nil, zero value otherwise.

### GetKeyLabelOk

`func (o *DeleteKeyRequest) GetKeyLabelOk() (*string, bool)`

GetKeyLabelOk returns a tuple with the KeyLabel field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKeyLabel

`func (o *DeleteKeyRequest) SetKeyLabel(v string)`

SetKeyLabel sets KeyLabel field to given value.


### GetId

`func (o *DeleteKeyRequest) GetId() int32`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *DeleteKeyRequest) GetIdOk() (*int32, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *DeleteKeyRequest) SetId(v int32)`

SetId sets Id field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


