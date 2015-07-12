package steamapi

import (
	"encoding/json"
	"time"

	"github.com/deckarep/golang-set"

	"appengine"
	"appengine/memcache"
)

//RetrieveFromStore retrieves the account from Store, or nil if no entry (assuming error is nil as well)
func RetrieveFromStore(c *appengine.Context, id string) (*SteamAccountDetails, error) {
	var acc SteamAccountDetails
	item, err := memcache.Get(*c, id)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(item.Value, acc); err != nil {
		return nil, err
	}
	return &acc, nil

}

//RetrieveMultiFromStore attempts to retrieve all elements from memcache
func RetrieveMultiFromStore(c *appengine.Context, ids []string) ([]SteamAccountDetails, []string, error) {
	allIDSet := mapset.NewSet()
	for _, id := range ids {
		allIDSet.Add(id)
	}

	foundIDs := mapset.NewSet()
	var foundAccounts []SteamAccountDetails

	items, err := memcache.GetMulti(*c, ids)
	if err != nil {
		return nil, nil, err
	}

	for key, item := range items {
		var acc SteamAccountDetails
		if err := json.Unmarshal(item.Value, &acc); err != nil {
			return nil, nil, err
		}
		foundIDs.Add(key)
		foundAccounts = append(foundAccounts, acc)
	}

	var missingIDs []string
	missingIDSet := allIDSet.Difference(foundIDs)
	for v := range missingIDSet.Iter() {
		missingIDs = append(missingIDs, v.(string))
	}

	return foundAccounts, missingIDs, nil

}

//SaveAllToStore itterates the array and stores them into the datastore
func SaveAllToStore(c appengine.Context, accounts []SteamAccountDetails) error {
	for _, account := range accounts {
		if err := SaveToStore(c, &account); err != nil {
			return err
		}
	}
	return nil
}

//SaveToStore Saves account details to store
func SaveToStore(c appengine.Context, account *SteamAccountDetails) error {
	account.LastUpdated = time.Now()
	account.Updated = false

	marshalledData, err := json.Marshal(account)
	if err != nil {
		return err
	}

	item := &memcache.Item{
		Key:        account.SteamID,
		Value:      marshalledData,
		Expiration: time.Duration(5) * time.Hour,
	}

	return memcache.Set(c, item)
}
