/*
Equinix Fabric API

Equinix Fabric is an advanced software-defined interconnection solution that enables you to directly, securely and dynamically connect to distributed infrastructure and digital ecosystems on platform Equinix via a single port, Customers can use Fabric to connect to: </br> 1. Cloud Service Providers - Clouds, network and other service providers.  </br> 2. Enterprises - Other Equinix customers, vendors and partners.  </br> 3. Myself - Another customer instance deployed at Equinix. </br>

API version: 4.2
Contact: api-support@equinix.com
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package v4

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
	"reflect"
)


// StatisticsApiService StatisticsApi service
type StatisticsApiService service

type ApiGetStatsByPortUUIDRequest struct {
	ctx context.Context
	ApiService *StatisticsApiService
	portId string
	startDateTime *time.Time
	endDateTime *time.Time
}

// startDateTime
func (r ApiGetStatsByPortUUIDRequest) StartDateTime(startDateTime time.Time) ApiGetStatsByPortUUIDRequest {
	r.startDateTime = &startDateTime
	return r
}

// endDateTime
func (r ApiGetStatsByPortUUIDRequest) EndDateTime(endDateTime time.Time) ApiGetStatsByPortUUIDRequest {
	r.endDateTime = &endDateTime
	return r
}

func (r ApiGetStatsByPortUUIDRequest) Execute() (*Statistics, *http.Response, error) {
	return r.ApiService.GetStatsByPortUUIDExecute(r)
}

/*
GetStatsByPortUUID Get Stats by uuid

This API provides service-level traffic metrics so that you can view access and gather key information required to manage service subscription sizing and capacity.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param portId Port UUID
 @return ApiGetStatsByPortUUIDRequest
*/
func (a *StatisticsApiService) GetStatsByPortUUID(ctx context.Context, portId string) ApiGetStatsByPortUUIDRequest {
	return ApiGetStatsByPortUUIDRequest{
		ApiService: a,
		ctx: ctx,
		portId: portId,
	}
}

// Execute executes the request
//  @return Statistics
func (a *StatisticsApiService) GetStatsByPortUUIDExecute(r ApiGetStatsByPortUUIDRequest) (*Statistics, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *Statistics
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "StatisticsApiService.GetStatsByPortUUID")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/fabric/v4/ports/{portId}/stats"
	localVarPath = strings.Replace(localVarPath, "{"+"portId"+"}", url.PathEscape(parameterToString(r.portId, "")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.startDateTime == nil {
		return localVarReturnValue, nil, reportError("startDateTime is required and must be specified")
	}
	if r.endDateTime == nil {
		return localVarReturnValue, nil, reportError("endDateTime is required and must be specified")
	}

	localVarQueryParams.Add("startDateTime", parameterToString(*r.startDateTime, ""))
	localVarQueryParams.Add("endDateTime", parameterToString(*r.endDateTime, ""))
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = ioutil.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		if localVarHTTPResponse.StatusCode == 401 {
			var v []Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 403 {
			var v []Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 500 {
			var v []Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.model = v
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type ApiGetStatsByVcUUIDRequest struct {
	ctx context.Context
	ApiService *StatisticsApiService
	connectionId string
	startDateTime *time.Time
	endDateTime *time.Time
	viewPoint *ViewPoint
}

// startDateTime
func (r ApiGetStatsByVcUUIDRequest) StartDateTime(startDateTime time.Time) ApiGetStatsByVcUUIDRequest {
	r.startDateTime = &startDateTime
	return r
}

// endDateTime
func (r ApiGetStatsByVcUUIDRequest) EndDateTime(endDateTime time.Time) ApiGetStatsByVcUUIDRequest {
	r.endDateTime = &endDateTime
	return r
}

// viewPoint
func (r ApiGetStatsByVcUUIDRequest) ViewPoint(viewPoint ViewPoint) ApiGetStatsByVcUUIDRequest {
	r.viewPoint = &viewPoint
	return r
}

func (r ApiGetStatsByVcUUIDRequest) Execute() (*Statistics, *http.Response, error) {
	return r.ApiService.GetStatsByVcUUIDExecute(r)
}

/*
GetStatsByVcUUID Get Stats by uuid

This API provides service-level metrics so that you can view access and gather key information required to manage service subscription sizing and capacity

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param connectionId Connection UUID
 @return ApiGetStatsByVcUUIDRequest
*/
func (a *StatisticsApiService) GetStatsByVcUUID(ctx context.Context, connectionId string) ApiGetStatsByVcUUIDRequest {
	return ApiGetStatsByVcUUIDRequest{
		ApiService: a,
		ctx: ctx,
		connectionId: connectionId,
	}
}

// Execute executes the request
//  @return Statistics
func (a *StatisticsApiService) GetStatsByVcUUIDExecute(r ApiGetStatsByVcUUIDRequest) (*Statistics, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *Statistics
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "StatisticsApiService.GetStatsByVcUUID")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/fabric/v4/connections/{connectionId}/stats"
	localVarPath = strings.Replace(localVarPath, "{"+"connectionId"+"}", url.PathEscape(parameterToString(r.connectionId, "")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.startDateTime == nil {
		return localVarReturnValue, nil, reportError("startDateTime is required and must be specified")
	}
	if r.endDateTime == nil {
		return localVarReturnValue, nil, reportError("endDateTime is required and must be specified")
	}
	if r.viewPoint == nil {
		return localVarReturnValue, nil, reportError("viewPoint is required and must be specified")
	}

	localVarQueryParams.Add("startDateTime", parameterToString(*r.startDateTime, ""))
	localVarQueryParams.Add("endDateTime", parameterToString(*r.endDateTime, ""))
	localVarQueryParams.Add("viewPoint", parameterToString(*r.viewPoint, ""))
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = ioutil.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		if localVarHTTPResponse.StatusCode == 401 {
			var v []Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 403 {
			var v []Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.model = v
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type ApiGetTopUtilizedPortStatisticsRequest struct {
	ctx context.Context
	ApiService *StatisticsApiService
	metros *[]string
	sort *Sort
	top *int32
	duration *Duration
	direction *QueryDirection
	metricInterval *MetricInterval
}

// Two-letter prefix indicating the metropolitan area in which a specified Equinix asset is located.
func (r ApiGetTopUtilizedPortStatisticsRequest) Metros(metros []string) ApiGetTopUtilizedPortStatisticsRequest {
	r.metros = &metros
	return r
}

// Key or set of keys that organizes the search payload by property (such as createdDate or metroCode) or by direction. Ascending (ASC) is the default value. The \&quot;‒\&quot; prefix indicates descending (DESC) order.
func (r ApiGetTopUtilizedPortStatisticsRequest) Sort(sort Sort) ApiGetTopUtilizedPortStatisticsRequest {
	r.sort = &sort
	return r
}

// Filter returning only the specified number of most heavily trafficked ports. The standard value is [1...10], and the default is 5.
func (r ApiGetTopUtilizedPortStatisticsRequest) Top(top int32) ApiGetTopUtilizedPortStatisticsRequest {
	r.top = &top
	return r
}

// duration
func (r ApiGetTopUtilizedPortStatisticsRequest) Duration(duration Duration) ApiGetTopUtilizedPortStatisticsRequest {
	r.duration = &duration
	return r
}

// Direction of traffic from the requester&#39;s viewpoint. The default is outbound.
func (r ApiGetTopUtilizedPortStatisticsRequest) Direction(direction QueryDirection) ApiGetTopUtilizedPortStatisticsRequest {
	r.direction = &direction
	return r
}

// metricInterval
func (r ApiGetTopUtilizedPortStatisticsRequest) MetricInterval(metricInterval MetricInterval) ApiGetTopUtilizedPortStatisticsRequest {
	r.metricInterval = &metricInterval
	return r
}

func (r ApiGetTopUtilizedPortStatisticsRequest) Execute() (*TopUtilizedStatistics, *http.Response, error) {
	return r.ApiService.GetTopUtilizedPortStatisticsExecute(r)
}

/*
GetTopUtilizedPortStatistics Top Port Statistics

This API provides top utilized service-level traffic metrics so that you can view access and gather key information required to manage service subscription sizing and capacity.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @return ApiGetTopUtilizedPortStatisticsRequest
*/
func (a *StatisticsApiService) GetTopUtilizedPortStatistics(ctx context.Context) ApiGetTopUtilizedPortStatisticsRequest {
	return ApiGetTopUtilizedPortStatisticsRequest{
		ApiService: a,
		ctx: ctx,
	}
}

// Execute executes the request
//  @return TopUtilizedStatistics
func (a *StatisticsApiService) GetTopUtilizedPortStatisticsExecute(r ApiGetTopUtilizedPortStatisticsRequest) (*TopUtilizedStatistics, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *TopUtilizedStatistics
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "StatisticsApiService.GetTopUtilizedPortStatistics")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/fabric/v4/ports/stats"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.metros == nil {
		return localVarReturnValue, nil, reportError("metros is required and must be specified")
	}

	if r.sort != nil {
		localVarQueryParams.Add("sort", parameterToString(*r.sort, ""))
	}
	if r.top != nil {
		localVarQueryParams.Add("top", parameterToString(*r.top, ""))
	}
	if r.duration != nil {
		localVarQueryParams.Add("duration", parameterToString(*r.duration, ""))
	}
	if r.direction != nil {
		localVarQueryParams.Add("direction", parameterToString(*r.direction, ""))
	}
	if r.metricInterval != nil {
		localVarQueryParams.Add("metricInterval", parameterToString(*r.metricInterval, ""))
	}
	{
		t := *r.metros
		if reflect.TypeOf(t).Kind() == reflect.Slice {
			s := reflect.ValueOf(t)
			for i := 0; i < s.Len(); i++ {
				localVarQueryParams.Add("metros", parameterToString(s.Index(i), "multi"))
			}
		} else {
			localVarQueryParams.Add("metros", parameterToString(t, "multi"))
		}
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = ioutil.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		if localVarHTTPResponse.StatusCode == 401 {
			var v []Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.model = v
			return localVarReturnValue, localVarHTTPResponse, newErr
		}
		if localVarHTTPResponse.StatusCode == 403 {
			var v []Error
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.model = v
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}
