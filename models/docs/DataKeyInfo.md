# DataKeyInfo

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Handle** | **int32** | HSM key handle | 
**Id** | **int32** | Key identifier | 
**SizeBits** | **int32** | Size of the key in bits | 

## Methods

### NewDataKeyInfo

`func NewDataKeyInfo(handle int32, id int32, sizeBits int32, ) *DataKeyInfo`

NewDataKeyInfo instantiates a new DataKeyInfo object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDataKeyInfoWithDefaults

`func NewDataKeyInfoWithDefaults() *DataKeyInfo`

NewDataKeyInfoWithDefaults instantiates a new DataKeyInfo object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetHandle

`func (o *DataKeyInfo) GetHandle() int32`

GetHandle returns the Handle field if non-nil, zero value otherwise.

### GetHandleOk

`func (o *DataKeyInfo) GetHandleOk() (*int32, bool)`

GetHandleOk returns a tuple with the Handle field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHandle

`func (o *DataKeyInfo) SetHandle(v int32)`

SetHandle sets Handle field to given value.


### GetId

`func (o *DataKeyInfo) GetId() int32`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *DataKeyInfo) GetIdOk() (*int32, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *DataKeyInfo) SetId(v int32)`

SetId sets Id field to given value.


### GetSizeBits

`func (o *DataKeyInfo) GetSizeBits() int32`

GetSizeBits returns the SizeBits field if non-nil, zero value otherwise.

### GetSizeBitsOk

`func (o *DataKeyInfo) GetSizeBitsOk() (*int32, bool)`

GetSizeBitsOk returns a tuple with the SizeBits field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSizeBits

`func (o *DataKeyInfo) SetSizeBits(v int32)`

SetSizeBits sets SizeBits field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


