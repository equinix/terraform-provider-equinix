package datalist

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// This is the configuration for a "data list" resource. It represents the schema and operations
// needed to create the data list resource.
type ResourceConfig struct {
	// The schema for a single instance of the resource.
	RecordSchema map[string]*schema.Schema

	// The name of the attribute in the resource through which to expose results.
	ResultAttributeName string

	// The description of the attribute in the resource through which to expose results.
	ResultAttributeDescription string

	// Given a record returned from the GetRecords function, flatten the record to a
	// map acceptable to the Set method on schema.ResourceData.
	FlattenRecord func(record, meta interface{}, extra map[string]interface{}) (map[string]interface{}, error)

	// Return all of the records on which the data list resource should operate.
	// The `meta` argument is the same meta argument passed into the resource's Read
	// function.
	GetRecords func(meta interface{}, extra map[string]interface{}) ([]interface{}, error)

	// Extra parameters to expose on the datasource alongside `filter` and `sort`.
	ExtraQuerySchema map[string]*schema.Schema
}

// Returns a new "data list" resource given the specified configuration. This
// is a resource with `filter` and `sort` attributes that can select a subset
// of records from a list of records for a particular type of resource.
func NewResource(config *ResourceConfig) *schema.Resource {
	err := validateResourceConfig(config)
	if err != nil {
		// Panic if the resource config is invalid since this will prevent the resource
		// from operating.
		log.Panicf("datalist.NewResource: invalid resource configuration: %v", err)
	}

	recordSchema := map[string]*schema.Schema{}
	for attributeName, attributeSchema := range config.RecordSchema {
		newAttributeSchema := &schema.Schema{}
		*newAttributeSchema = *attributeSchema
		newAttributeSchema.Computed = true
		newAttributeSchema.Required = false
		newAttributeSchema.Optional = false
		recordSchema[attributeName] = newAttributeSchema
	}

	filterAttributes := computeFilterAttributes(recordSchema)
	sortAttributes := computeSortAttributes(recordSchema)

	datasourceSchema := map[string]*schema.Schema{
		"filter": filterSchema(filterAttributes),
		"sort":   sortSchema(sortAttributes),
		config.ResultAttributeName: {
			Type:        schema.TypeList,
			Computed:    true,
			Description: config.ResultAttributeDescription,
			Elem: &schema.Resource{
				Schema: recordSchema,
			},
		},
	}

	for attr, value := range config.ExtraQuerySchema {
		datasourceSchema[attr] = value
	}

	return &schema.Resource{
		ReadContext: dataListResourceRead(config),
		Schema:      datasourceSchema,
	}
}

func dataListResourceRead(config *ResourceConfig) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		extra := map[string]interface{}{}
		for attr := range config.ExtraQuerySchema {
			extra[attr] = d.Get(attr)
		}

		records, err := config.GetRecords(meta, extra)
		if err != nil {
			return diag.Errorf("Unable to load records: %s", err)
		}

		flattenedRecords := make([]map[string]interface{}, len(records))
		for i, record := range records {
			flattenedRecord, err := config.FlattenRecord(record, meta, extra)
			if err != nil {
				return diag.FromErr(err)
			}
			flattenedRecords[i] = flattenedRecord
		}

		if v, ok := d.GetOk("filter"); ok {
			filters, err := expandFilters(config.RecordSchema, v.(*schema.Set).List())
			if err != nil {
				return diag.FromErr(err)
			}
			flattenedRecords = applyFilters(config.RecordSchema, flattenedRecords, filters)
		}

		if v, ok := d.GetOk("sort"); ok {
			sorts := expandSorts(v.([]interface{}))
			flattenedRecords = applySorts(config.RecordSchema, flattenedRecords, sorts)
		}

		d.SetId(resource.UniqueId())

		if err := d.Set(config.ResultAttributeName, flattenedRecords); err != nil {
			return diag.Errorf("unable to set `%s` attribute: %s", config.ResultAttributeName, err)
		}

		return nil
	}
}

// Compute the set of filter attributes for the resource.
func computeFilterAttributes(recordSchema map[string]*schema.Schema) []string {
	var filterAttributes []string

	for attr, schemaForAttr := range recordSchema {
		if schemaForAttr.Type != schema.TypeMap {
			filterAttributes = append(filterAttributes, attr)
		}
	}

	return filterAttributes
}

// Compute the set of sort attributes for the source.
func computeSortAttributes(recordSchema map[string]*schema.Schema) []string {
	var sortAttributes []string

	for attr, schemaForAttr := range recordSchema {
		supported := false
		switch schemaForAttr.Type {
		case schema.TypeString, schema.TypeBool, schema.TypeInt, schema.TypeFloat:
			supported = true
		}

		if supported {
			sortAttributes = append(sortAttributes, attr)
		}
	}

	return sortAttributes
}

// Validate a ResourceConfig to ensure it conforms to this package's assumptions.
func validateResourceConfig(config *ResourceConfig) error {
	// Ensure that ResultAttributeName exists.
	if config.ResultAttributeName == "" {
		return fmt.Errorf("ResultAttributeName must be specified")
	}

	return nil
}
