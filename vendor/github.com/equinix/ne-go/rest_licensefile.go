package ne

import (
	"io"
	"net/http"

	"github.com/equinix/ne-go/internal/api"
)

//UploadLicenseFile performs multipart upload of a license file from a given reader interface
//along with provided data. Uploaded file identifier is returned on success.
func (c RestClient) UploadLicenseFile(metroCode, deviceTypeCode, deviceManagementMode, licenseMode, fileName string, reader io.Reader) (*string, error) {
	path := "/ne/v1/devices/licenseFiles"
	respBody := api.LicenseFileUploadResponse{}
	req := c.R().
		SetFileReader("file", fileName, reader).
		SetFormData(map[string]string{
			"metroCode":            metroCode,
			"deviceTypeCode":       deviceTypeCode,
			"licenseType":          licenseMode,
			"deviceManagementType": deviceManagementMode,
		}).
		SetResult(&respBody)
	if err := c.Execute(req, http.MethodPost, path); err != nil {
		return nil, err
	}
	return respBody.FileID, nil
}
