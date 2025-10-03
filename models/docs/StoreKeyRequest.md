# StoreKeyRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**KeyLabel** | **string** | Etiqueta para la clave almacenada | 
**Id** | **string** | Identificador único | 
**KeyValue** | **string** | Valor de la clave en Base64 | 

## Methods

### NewStoreKeyRequest

`func NewStoreKeyRequest(keyLabel string, id string, keyValue string, ) *StoreKeyRequest`

NewStoreKeyRequest instantiates a new StoreKeyRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewStoreKeyRequestWithDefaults

`func NewStoreKeyRequestWithDefaults() *StoreKeyRequest`

NewStoreKeyRequestWithDefaults instantiates a new StoreKeyRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetKeyLabel

`func (o *StoreKeyRequest) GetKeyLabel() string`

GetKeyLabel returns the KeyLabel field if non-nil, zero value otherwise.

### GetKeyLabelOk

`func (o *StoreKeyRequest) GetKeyLabelOk() (*string, bool)`

GetKeyLabelOk returns a tuple with the KeyLabel field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKeyLabel

`func (o *StoreKeyRequest) SetKeyLabel(v string)`

SetKeyLabel sets KeyLabel field to given value.


### GetId

`func (o *StoreKeyRequest) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *StoreKeyRequest) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *StoreKeyRequest) SetId(v string)`

SetId sets Id field to given value.


### GetKeyValue

`func (o *StoreKeyRequest) GetKeyValue() string`

GetKeyValue returns the KeyValue field if non-nil, zero value otherwise.

### GetKeyValueOk

`func (o *StoreKeyRequest) GetKeyValueOk() (*string, bool)`

GetKeyValueOk returns a tuple with the KeyValue field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKeyValue

`func (o *StoreKeyRequest) SetKeyValue(v string)`

SetKeyValue sets KeyValue field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


