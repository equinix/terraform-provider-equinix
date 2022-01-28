package ne

import (
	"net/http"
	"net/url"

	"github.com/equinix/ne-go/internal/api"
)

//GetAccounts retrieves accounts and their details for a given metro code using Network Edge API
func (c RestClient) GetAccounts(metroCode string) ([]Account, error) {
	path := "/ne/v1/accounts/" + url.PathEscape(metroCode)
	respBody := api.AccountResponse{}
	req := c.R().SetResult(&respBody)

	if err := c.Execute(req, http.MethodGet, path); err != nil {
		return nil, err
	}
	return mapAccountsAPIToDomain(respBody.Accounts), nil
}

func mapAccountsAPIToDomain(apiAccounts []api.Account) []Account {
	transformed := make([]Account, len(apiAccounts))
	for i := range apiAccounts {
		transformed[i] = Account{
			Name:   apiAccounts[i].Name,
			Number: apiAccounts[i].Number,
			UCMID:  apiAccounts[i].UCMID,
			Status: apiAccounts[i].Status,
		}
	}
	return transformed
}
