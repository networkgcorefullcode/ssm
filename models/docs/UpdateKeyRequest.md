# UpdateKeyRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**KeyLabel** | **string** | Label of the key to update | 
**Id** | **int32** | Key identifier | 
**KeyValue** | **string** | New key value in hexadecimal format | 
**KeyType** | **string** | Type of cryptographic key | 

## Methods

### NewUpdateKeyRequest

`func NewUpdateKeyRequest(keyLabel string, id int32, keyValue string, keyType string, ) *UpdateKeyRequest`

NewUpdateKeyRequest instantiates a new UpdateKeyRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUpdateKeyRequestWithDefaults

`func NewUpdateKeyRequestWithDefaults() *UpdateKeyRequest`

NewUpdateKeyRequestWithDefaults instantiates a new UpdateKeyRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetKeyLabel

`func (o *UpdateKeyRequest) GetKeyLabel() string`

GetKeyLabel returns the KeyLabel field if non-nil, zero value otherwise.

### GetKeyLabelOk

`func (o *UpdateKeyRequest) GetKeyLabelOk() (*string, bool)`

GetKeyLabelOk returns a tuple with the KeyLabel field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKeyLabel

`func (o *UpdateKeyRequest) SetKeyLabel(v string)`

SetKeyLabel sets KeyLabel field to given value.


### GetId

`func (o *UpdateKeyRequest) GetId() int32`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *UpdateKeyRequest) GetIdOk() (*int32, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *UpdateKeyRequest) SetId(v int32)`

SetId sets Id field to given value.


### GetKeyValue

`func (o *UpdateKeyRequest) GetKeyValue() string`

GetKeyValue returns the KeyValue field if non-nil, zero value otherwise.

### GetKeyValueOk

`func (o *UpdateKeyRequest) GetKeyValueOk() (*string, bool)`

GetKeyValueOk returns a tuple with the KeyValue field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKeyValue

`func (o *UpdateKeyRequest) SetKeyValue(v string)`

SetKeyValue sets KeyValue field to given value.


### GetKeyType

`func (o *UpdateKeyRequest) GetKeyType() string`

GetKeyType returns the KeyType field if non-nil, zero value otherwise.

### GetKeyTypeOk

`func (o *UpdateKeyRequest) GetKeyTypeOk() (*string, bool)`

GetKeyTypeOk returns a tuple with the KeyType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKeyType

`func (o *UpdateKeyRequest) SetKeyType(v string)`

SetKeyType sets KeyType field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


