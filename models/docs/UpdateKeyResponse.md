# UpdateKeyResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Message** | **string** | Confirmation message | 
**Handle** | **int32** | New HSM key handle | 
**KeyLabel** | **string** | Label of the updated key | 
**CipherKey** | Pointer to **string** | Encrypted key value (optional) | [optional] 

## Methods

### NewUpdateKeyResponse

`func NewUpdateKeyResponse(message string, handle int32, keyLabel string, ) *UpdateKeyResponse`

NewUpdateKeyResponse instantiates a new UpdateKeyResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUpdateKeyResponseWithDefaults

`func NewUpdateKeyResponseWithDefaults() *UpdateKeyResponse`

NewUpdateKeyResponseWithDefaults instantiates a new UpdateKeyResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetMessage

`func (o *UpdateKeyResponse) GetMessage() string`

GetMessage returns the Message field if non-nil, zero value otherwise.

### GetMessageOk

`func (o *UpdateKeyResponse) GetMessageOk() (*string, bool)`

GetMessageOk returns a tuple with the Message field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMessage

`func (o *UpdateKeyResponse) SetMessage(v string)`

SetMessage sets Message field to given value.


### GetHandle

`func (o *UpdateKeyResponse) GetHandle() int32`

GetHandle returns the Handle field if non-nil, zero value otherwise.

### GetHandleOk

`func (o *UpdateKeyResponse) GetHandleOk() (*int32, bool)`

GetHandleOk returns a tuple with the Handle field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHandle

`func (o *UpdateKeyResponse) SetHandle(v int32)`

SetHandle sets Handle field to given value.


### GetKeyLabel

`func (o *UpdateKeyResponse) GetKeyLabel() string`

GetKeyLabel returns the KeyLabel field if non-nil, zero value otherwise.

### GetKeyLabelOk

`func (o *UpdateKeyResponse) GetKeyLabelOk() (*string, bool)`

GetKeyLabelOk returns a tuple with the KeyLabel field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKeyLabel

`func (o *UpdateKeyResponse) SetKeyLabel(v string)`

SetKeyLabel sets KeyLabel field to given value.


### GetCipherKey

`func (o *UpdateKeyResponse) GetCipherKey() string`

GetCipherKey returns the CipherKey field if non-nil, zero value otherwise.

### GetCipherKeyOk

`func (o *UpdateKeyResponse) GetCipherKeyOk() (*string, bool)`

GetCipherKeyOk returns a tuple with the CipherKey field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCipherKey

`func (o *UpdateKeyResponse) SetCipherKey(v string)`

SetCipherKey sets CipherKey field to given value.

### HasCipherKey

`func (o *UpdateKeyResponse) HasCipherKey() bool`

HasCipherKey returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


