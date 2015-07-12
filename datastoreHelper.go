package steamapi

import (
	"strconv"
	"time"

	"appengine"
	"appengine/datastore"
)

//RetrieveFromStore retrieves the account from Store, or nil if no entry (assuming error is nil as well)
func RetrieveFromStore(c *appengine.Context, id string) (*SteamAccountDetails, error) {
	var acc SteamAccountDetails
	key, keyErr := steamAccountDetailsKey(c, &id)
	if keyErr != nil {
		return nil, keyErr
	}
	err := datastore.Get(*c, key, acc)

	if err == datastore.ErrNoSuchEntity {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return &acc, nil
	}
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
	key, keyErr := steamAccountDetailsKey(&c, &account.SteamID)
	if keyErr != nil {
		return keyErr
	}
	_, err := datastore.Put(c, key, account)
	return err
}

func steamAccountDetailsKey(c *appengine.Context, id *string) (*datastore.Key, error) {
	intID, err := strconv.ParseInt(*id, 0, 64)
	if err != nil {
		return nil, err
	}
	return datastore.NewKey(*c, "SteamAccountDetails", "", intID, nil), nil
}
