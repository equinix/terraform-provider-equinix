package api

//ConfirmL2ConnectionRequest patch l2 connections request
type ConfirmL2ConnectionRequest struct {
	AccessKey *string `json:"accessKey,omitempty"`
	SecretKey *string `json:"secretKey,omitempty"`
}

//ConfirmL2ConnectionResponse patch l2 connection response
type ConfirmL2ConnectionResponse struct {
	Message             *string `json:"message,omitempty"`
	PrimaryConnectionID *string `json:"primaryConnectionId,omitempty"`
}
