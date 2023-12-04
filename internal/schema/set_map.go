package utils

import (
	"sort"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type setFn = func(d *schema.ResourceData, key string) error

// setMap sets the map of values to ResourceData, checking and returning the
// errors. Typically d.Set is not error checked. This helper makes checking
// those errors less tedious. Because this works with a map, the order of the
// errors would not be predictable, to avoid this the errors will be sorted.
func SetMap(d *schema.ResourceData, m map[string]interface{}) error {
	errs := &multierror.Error{}
	for key, v := range m {
		var err error
		if f, ok := v.(setFn); ok {
			err = f(d, key)
		} else {
			if key == "router" {
				d.Set("gateway", v)
			}
			err = d.Set(key, v)
		}

		if err != nil {
			errs = multierror.Append(errs, err)
		}
	}
	sort.Sort(errs)

	return errs.ErrorOrNil()
}
