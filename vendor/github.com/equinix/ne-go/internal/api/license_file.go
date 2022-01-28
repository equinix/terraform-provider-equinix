package api

//LicenseFileUploadResponse describes response to license file
//upload request
type LicenseFileUploadResponse struct {
	FileID *string `json:"fileId,omitempty"`
}
