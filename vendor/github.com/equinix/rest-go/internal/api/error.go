package api

//ErrorResponses describes error response built with
//multiple error responses
type ErrorResponses []ErrorResponse

//ErrorResponse describes error response with standardized
//application error description
type ErrorResponse struct {
	ErrorCode    string `json:"errorCode,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
	MoreInfo     string `json:"moreInfo,omitempty"`
	Property     string `json:"property,omitempty"`
}
