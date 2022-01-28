package rest

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/go-resty/resty/v2"
)

//PagingConfig is used to describe pagination aspects like naming of query parameters
type PagingConfig struct {
	//TotalCountFieldName is name of a field in a response struct that indicates total number of elements in collection
	TotalCountFieldName string
	//ContentFieldName is a name of a field that holds slice of elements
	ContentFieldName string
	//SizeParamName is a name of page size query parameter
	SizeParamName string
	//PageParamName is a name of page number query parameter
	PageParamName string
	//FirstPageNumber is a first page number returned by the server (typically 0 or 1)
	FirstPageNumber int
	//AdditionalParams is a map of additional query parameters that will be set in GET request
	AdditionalParams map[string]string
}

//DefaultPagingConfig returns PagingConfig with default values
func DefaultPagingConfig() *PagingConfig {
	return &PagingConfig{
		TotalCountFieldName: "TotalCount",
		ContentFieldName:    "Content",
		SizeParamName:       "size",
		PageParamName:       "page",
		FirstPageNumber:     1,
		AdditionalParams:    make(map[string]string),
	}
}

//SetTotalCountFieldName sets totalCount field name
func (c *PagingConfig) SetTotalCountFieldName(v string) *PagingConfig {
	c.TotalCountFieldName = v
	return c
}

//SetContentFieldName sets totalCount field name
func (c *PagingConfig) SetContentFieldName(v string) *PagingConfig {
	c.ContentFieldName = v
	return c
}

//SetSizeParamName sets size query parameter name
func (c *PagingConfig) SetSizeParamName(v string) *PagingConfig {
	c.SizeParamName = v
	return c
}

//SetPageParamName sets page query parameter name
func (c *PagingConfig) SetPageParamName(v string) *PagingConfig {
	c.PageParamName = v
	return c
}

//SetFirstPageNumber sets number of a fist page as returned by the server
func (c *PagingConfig) SetFirstPageNumber(v int) *PagingConfig {
	c.FirstPageNumber = v
	return c
}

//SetAdditionalParams sets additional query parameters that will be used in a request
func (c *PagingConfig) SetAdditionalParams(v map[string]string) *PagingConfig {
	c.AdditionalParams = v
	return c
}

//OffsetPaginationConfig is used to describe pagination aspects
type OffsetPaginationConfig struct {
	PaginationFieldName string
	OffsetFieldName     string
	LimitFieldName      string
	TotalFieldName      string
	DataFieldName       string
	AdditionalParams    map[string]string
}

//DefaultOffsetPagingConfig returns OffsetPaginationConfig with default values
func DefaultOffsetPagingConfig() *OffsetPaginationConfig {
	return &OffsetPaginationConfig{
		PaginationFieldName: "Pagination",
		DataFieldName:       "Data",
		TotalFieldName:      "Total",
		OffsetFieldName:     "offset",
		LimitFieldName:      "limit",
		AdditionalParams:    make(map[string]string),
	}
}

//SetPaginationFieldName sets pagination field name
func (c *OffsetPaginationConfig) SetPaginationFieldName(name string) *OffsetPaginationConfig {
	c.PaginationFieldName = name
	return c
}

//SetOffsetFieldName sets offset field name
func (c *OffsetPaginationConfig) SetOffsetFieldName(name string) *OffsetPaginationConfig {
	c.OffsetFieldName = name
	return c
}

//SetLimitFieldName sets limit field name
func (c *OffsetPaginationConfig) SetLimitFieldName(name string) *OffsetPaginationConfig {
	c.LimitFieldName = name
	return c
}

//SetTotalFieldName sets total field name
func (c *OffsetPaginationConfig) SetTotalFieldName(name string) *OffsetPaginationConfig {
	c.TotalFieldName = name
	return c
}

//SetDataFieldName sets data field name
func (c *OffsetPaginationConfig) SetDataFieldName(name string) *OffsetPaginationConfig {
	c.DataFieldName = name
	return c
}

//SetAdditionalParams sets additional query parameters that will be used in a request
func (c *OffsetPaginationConfig) SetAdditionalParams(v map[string]string) *OffsetPaginationConfig {
	c.AdditionalParams = v
	return c
}

//GetPaginated uses HTTP GET requests to retrieve list of all objects from paginated responses.
//Requests are executed against given path, pagination aspects are controlled with PagingConfig
func (c Client) GetPaginated(path string, result interface{}, conf *PagingConfig) ([]interface{}, error) {
	if reflect.ValueOf(result).Kind() != reflect.Ptr {
		return nil, fmt.Errorf("operation failed, provided result is not a ptr")
	}
	req := c.R().SetResult(result).
		SetQueryParams(conf.AdditionalParams).
		SetQueryParam(conf.SizeParamName, strconv.Itoa(c.PageSize))
	if err := c.Execute(req, resty.MethodGet, path); err != nil {
		return nil, err
	}
	totalValue, err := getFieldValueFromStruct(result, conf.TotalCountFieldName, reflect.Int)
	if err != nil {
		return nil, err
	}
	totalCount := totalValue.Interface().(int)
	contentValue, err := getFieldValueFromStruct(result, conf.ContentFieldName, reflect.Slice)
	if err != nil {
		return nil, err
	}
	content := make([]interface{}, 0, totalCount)
	content = appendSliceValue(content, contentValue)
	recordsFetched := c.PageSize
	isLast := false
	if recordsFetched >= totalCount {
		isLast = true
	}
	for pageNum := conf.FirstPageNumber + 1; !isLast; pageNum++ {
		resValue := reflect.ValueOf(result)
		if resValue.Kind() == reflect.Ptr {
			resValue = resValue.Elem()
		}
		nextResult := reflect.New(resValue.Type()).Interface()
		req := c.R().SetResult(nextResult).
			SetQueryParams(conf.AdditionalParams).
			SetQueryParam(conf.SizeParamName, strconv.Itoa(c.PageSize)).
			SetQueryParam(conf.PageParamName, strconv.Itoa(pageNum))
		if err := c.Execute(req, resty.MethodGet, path); err != nil {
			return nil, err
		}
		resContent, err := getFieldValueFromStruct(nextResult, conf.ContentFieldName, reflect.Slice)
		if err != nil {
			return nil, err
		}
		content = appendSliceValue(content, resContent)
		recordsFetched += c.PageSize
		if recordsFetched >= totalCount {
			isLast = true
		}
	}
	return content, nil
}

//GetOffsetPaginated uses HTTP GET requests to retrieve list of all objects from
//paginated responses that use offset & limit attributes in a separate pagination object.
//Requests are executed against given path, pagination aspects are controlled with PagingConfig
func (c Client) GetOffsetPaginated(path string, result interface{}, conf *OffsetPaginationConfig) ([]interface{}, error) {
	if reflect.ValueOf(result).Kind() != reflect.Ptr {
		return nil, fmt.Errorf("operation failed, provided result is not a ptr")
	}
	req := c.R().SetResult(result).
		SetQueryParams(conf.AdditionalParams).
		SetQueryParam(conf.LimitFieldName, strconv.Itoa(c.PageSize))
	if err := c.Execute(req, resty.MethodGet, path); err != nil {
		return nil, err
	}
	paginationData, err := getFieldValueFromStruct(result, conf.PaginationFieldName, reflect.Struct)
	if err != nil {
		return nil, err
	}
	totalValue, err := getFieldValueFromStruct(paginationData.Interface(), conf.TotalFieldName, reflect.Int)
	if err != nil {
		return nil, err
	}
	dataValue, err := getFieldValueFromStruct(result, conf.DataFieldName, reflect.Slice)
	if err != nil {
		return nil, err
	}
	totalCount := totalValue.Interface().(int)
	data := make([]interface{}, 0, totalCount)
	data = appendSliceValue(data, dataValue)
	for offset := c.PageSize; offset < totalCount; {
		resValue := reflect.ValueOf(result)
		if resValue.Kind() == reflect.Ptr {
			resValue = resValue.Elem()
		}
		nextResult := reflect.New(resValue.Type()).Interface()
		req := c.R().SetResult(nextResult).
			SetQueryParams(conf.AdditionalParams).
			SetQueryParam(conf.LimitFieldName, strconv.Itoa(c.PageSize)).
			SetQueryParam(conf.OffsetFieldName, strconv.Itoa(offset))
		if err := c.Execute(req, resty.MethodGet, path); err != nil {
			return nil, err
		}
		responseData, err := getFieldValueFromStruct(nextResult, conf.DataFieldName, reflect.Slice)
		if err != nil {
			return nil, err
		}
		data = appendSliceValue(data, responseData)
		offset += c.PageSize
	}
	return data, nil
}

//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Unexported package methods
//_______________________________________________________________________

func getFieldValueFromStruct(target interface{}, fieldName string, fieldKind reflect.Kind) (*reflect.Value, error) {
	resultVal := reflect.ValueOf(target)
	if resultVal.Kind() == reflect.Ptr {
		resultVal = resultVal.Elem()
	}
	if resultVal.Kind() != reflect.Struct {
		return nil, fmt.Errorf("provided target is %s and not a struct", resultVal.Kind())
	}
	val := resultVal.FieldByName(fieldName)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != fieldKind {
		return nil, fmt.Errorf("kind of %s field in target struct is %s and not %s", fieldName, val.Kind(), fieldKind)
	}
	return &val, nil
}

func appendSliceValue(target []interface{}, source *reflect.Value) []interface{} {
	transformed := make([]interface{}, source.Len())
	for i := 0; i < source.Len(); i++ {
		transformed[i] = source.Index(i).Interface()
	}
	return append(target, transformed...)
}
