# GetAllKeysResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**KeysByLabel** | [**map[string][]GetKeyResponseKeyInfo**](array.md) | Map of label to array of key information | 
**TotalKeys** | **int32** | Total number of keys found | 
**TotalLabels** | **int32** | Total number of unique labels | 

## Methods

### NewGetAllKeysResponse

`func NewGetAllKeysResponse(keysByLabel map[string][]GetKeyResponseKeyInfo, totalKeys int32, totalLabels int32, ) *GetAllKeysResponse`

NewGetAllKeysResponse instantiates a new GetAllKeysResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewGetAllKeysResponseWithDefaults

`func NewGetAllKeysResponseWithDefaults() *GetAllKeysResponse`

NewGetAllKeysResponseWithDefaults instantiates a new GetAllKeysResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetKeysByLabel

`func (o *GetAllKeysResponse) GetKeysByLabel() map[string][]GetKeyResponseKeyInfo`

GetKeysByLabel returns the KeysByLabel field if non-nil, zero value otherwise.

### GetKeysByLabelOk

`func (o *GetAllKeysResponse) GetKeysByLabelOk() (*map[string][]GetKeyResponseKeyInfo, bool)`

GetKeysByLabelOk returns a tuple with the KeysByLabel field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKeysByLabel

`func (o *GetAllKeysResponse) SetKeysByLabel(v map[string][]GetKeyResponseKeyInfo)`

SetKeysByLabel sets KeysByLabel field to given value.


### GetTotalKeys

`func (o *GetAllKeysResponse) GetTotalKeys() int32`

GetTotalKeys returns the TotalKeys field if non-nil, zero value otherwise.

### GetTotalKeysOk

`func (o *GetAllKeysResponse) GetTotalKeysOk() (*int32, bool)`

GetTotalKeysOk returns a tuple with the TotalKeys field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotalKeys

`func (o *GetAllKeysResponse) SetTotalKeys(v int32)`

SetTotalKeys sets TotalKeys field to given value.


### GetTotalLabels

`func (o *GetAllKeysResponse) GetTotalLabels() int32`

GetTotalLabels returns the TotalLabels field if non-nil, zero value otherwise.

### GetTotalLabelsOk

`func (o *GetAllKeysResponse) GetTotalLabelsOk() (*int32, bool)`

GetTotalLabelsOk returns a tuple with the TotalLabels field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotalLabels

`func (o *GetAllKeysResponse) SetTotalLabels(v int32)`

SetTotalLabels sets TotalLabels field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


