# DeleteKeyResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Message** | **string** | Confirmation message | 
**KeyLabel** | **string** | Label of the deleted key | 

## Methods

### NewDeleteKeyResponse

`func NewDeleteKeyResponse(message string, keyLabel string, ) *DeleteKeyResponse`

NewDeleteKeyResponse instantiates a new DeleteKeyResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDeleteKeyResponseWithDefaults

`func NewDeleteKeyResponseWithDefaults() *DeleteKeyResponse`

NewDeleteKeyResponseWithDefaults instantiates a new DeleteKeyResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetMessage

`func (o *DeleteKeyResponse) GetMessage() string`

GetMessage returns the Message field if non-nil, zero value otherwise.

### GetMessageOk

`func (o *DeleteKeyResponse) GetMessageOk() (*string, bool)`

GetMessageOk returns a tuple with the Message field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMessage

`func (o *DeleteKeyResponse) SetMessage(v string)`

SetMessage sets Message field to given value.


### GetKeyLabel

`func (o *DeleteKeyResponse) GetKeyLabel() string`

GetKeyLabel returns the KeyLabel field if non-nil, zero value otherwise.

### GetKeyLabelOk

`func (o *DeleteKeyResponse) GetKeyLabelOk() (*string, bool)`

GetKeyLabelOk returns a tuple with the KeyLabel field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKeyLabel

`func (o *DeleteKeyResponse) SetKeyLabel(v string)`

SetKeyLabel sets KeyLabel field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


