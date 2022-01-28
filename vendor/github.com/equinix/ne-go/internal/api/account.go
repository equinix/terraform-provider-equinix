package api

//AccountResponse describes response to Network Edge
//billing account query
type AccountResponse struct {
	Accounts []Account `json:"accounts,omitempty"`
}

//Account describes Network Edge billing account
type Account struct {
	Name   *string `json:"accountName,omitempty"`
	Number *string `json:"accountNumber,omitempty"`
	UCMID  *string `json:"accountUcmId,omitempty"`
	Status *string `json:"accountStatus,omitempty"`
}
