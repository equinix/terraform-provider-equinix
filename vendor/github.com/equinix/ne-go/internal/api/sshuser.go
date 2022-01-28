package api

//SSHUser describes network edge SSH user
type SSHUser struct {
	UUID        *string  `json:"uuid,omitempty"`
	Username    *string  `json:"username,omitempty"`
	Password    *string  `json:"password,omitempty"`
	DeviceUUIDs []string `json:"deviceUUIDs,omitempty"`
}

//SSHUserRequest describes network edge SSH user creation request
type SSHUserRequest struct {
	Username   *string `json:"username,omitempty"`
	Password   *string `json:"password,omitempty"`
	DeviceUUID *string `json:"deviceUuid,omitempty"`
}

//SSHUserUpdateRequest describes network edge SSH user update request
type SSHUserUpdateRequest struct {
	Password *string `json:"password,omitempty"`
}

//SSHUsersResponse describes response for a get ssh user list request
type SSHUsersResponse struct {
	Pagination Pagination `json:"pagination,omitempty"`
	Data       []SSHUser  `json:"data,omitempty"`
}
