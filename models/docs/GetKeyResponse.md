# GetKeyResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**KeyInfo** | [**DataKeyInfo**](DataKeyInfo.md) |  | 

## Methods

### NewGetKeyResponse

`func NewGetKeyResponse(keyInfo DataKeyInfo, ) *GetKeyResponse`

NewGetKeyResponse instantiates a new GetKeyResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewGetKeyResponseWithDefaults

`func NewGetKeyResponseWithDefaults() *GetKeyResponse`

NewGetKeyResponseWithDefaults instantiates a new GetKeyResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetKeyInfo

`func (o *GetKeyResponse) GetKeyInfo() DataKeyInfo`

GetKeyInfo returns the KeyInfo field if non-nil, zero value otherwise.

### GetKeyInfoOk

`func (o *GetKeyResponse) GetKeyInfoOk() (*DataKeyInfo, bool)`

GetKeyInfoOk returns a tuple with the KeyInfo field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKeyInfo

`func (o *GetKeyResponse) SetKeyInfo(v DataKeyInfo)`

SetKeyInfo sets KeyInfo field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


