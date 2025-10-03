# ProblemDetails

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Title** | Pointer to **string** | Título breve del problema | [optional] 
**Detail** | Pointer to **string** | Descripción detallada del problema | [optional] 
**Status** | Pointer to **int32** | Código de estado HTTP | [optional] 
**Error** | Pointer to **string** | Código de error interno | [optional] 
**Instance** | Pointer to **string** | URI que identifica la ocurrencia específica | [optional] 

## Methods

### NewProblemDetails

`func NewProblemDetails() *ProblemDetails`

NewProblemDetails instantiates a new ProblemDetails object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewProblemDetailsWithDefaults

`func NewProblemDetailsWithDefaults() *ProblemDetails`

NewProblemDetailsWithDefaults instantiates a new ProblemDetails object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetTitle

`func (o *ProblemDetails) GetTitle() string`

GetTitle returns the Title field if non-nil, zero value otherwise.

### GetTitleOk

`func (o *ProblemDetails) GetTitleOk() (*string, bool)`

GetTitleOk returns a tuple with the Title field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTitle

`func (o *ProblemDetails) SetTitle(v string)`

SetTitle sets Title field to given value.

### HasTitle

`func (o *ProblemDetails) HasTitle() bool`

HasTitle returns a boolean if a field has been set.

### GetDetail

`func (o *ProblemDetails) GetDetail() string`

GetDetail returns the Detail field if non-nil, zero value otherwise.

### GetDetailOk

`func (o *ProblemDetails) GetDetailOk() (*string, bool)`

GetDetailOk returns a tuple with the Detail field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDetail

`func (o *ProblemDetails) SetDetail(v string)`

SetDetail sets Detail field to given value.

### HasDetail

`func (o *ProblemDetails) HasDetail() bool`

HasDetail returns a boolean if a field has been set.

### GetStatus

`func (o *ProblemDetails) GetStatus() int32`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *ProblemDetails) GetStatusOk() (*int32, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *ProblemDetails) SetStatus(v int32)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *ProblemDetails) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetError

`func (o *ProblemDetails) GetError() string`

GetError returns the Error field if non-nil, zero value otherwise.

### GetErrorOk

`func (o *ProblemDetails) GetErrorOk() (*string, bool)`

GetErrorOk returns a tuple with the Error field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetError

`func (o *ProblemDetails) SetError(v string)`

SetError sets Error field to given value.

### HasError

`func (o *ProblemDetails) HasError() bool`

HasError returns a boolean if a field has been set.

### GetInstance

`func (o *ProblemDetails) GetInstance() string`

GetInstance returns the Instance field if non-nil, zero value otherwise.

### GetInstanceOk

`func (o *ProblemDetails) GetInstanceOk() (*string, bool)`

GetInstanceOk returns a tuple with the Instance field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInstance

`func (o *ProblemDetails) SetInstance(v string)`

SetInstance sets Instance field to given value.

### HasInstance

`func (o *ProblemDetails) HasInstance() bool`

HasInstance returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


