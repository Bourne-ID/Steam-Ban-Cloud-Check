package steamapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"appengine"
	"appengine/urlfetch"
)

//SteamAccount object containing an array of all returned Steam Accounts
type SteamAccount struct {
	Players []SteamAccountDetails
}

//SteamAccountDetails gives the details of the SteamAccount
type SteamAccountDetails struct {
	SteamID          string
	CommunityBanned  bool
	VACBanned        bool
	NumberOfVACBans  int
	DaysSinceLastBan int
	NumberOfGameBans int
	EconomyBan       string
	LastUpdated      time.Time
	Updated          bool
}

func init() {
	http.HandleFunc("/", root)
}

func root(w http.ResponseWriter, r *http.Request) {
	steamIdsString := r.FormValue("steamids")
	steamIDArray := strings.Split(steamIdsString, ",")

	//do we have steam keys?
	if len(steamIDArray) == 0 {
		http.Error(w, "No Steam IDs present", http.StatusBadRequest)
		return
	}

	//has the api key been set up correctly?
	key, err := readAPIKey()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//Google specific
	c := appengine.NewContext(r)

	//check if we already have some of there entries in store and is recent
	//uncachedIDs, obtainedRecords, err := callMemcache(&c, steamIDArray)
	foundAccounts, missingIDs, err1 := callMemcache(&c, steamIDArray)
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusInternalServerError)
		return
	}

	//Steam API only allows 100 steamIds to be sent, group them up into a Map
	groupedSteamIDArray := groupSteamIDs(missingIDs)

	results, err := makeSteamAPICall(&c, groupedSteamIDArray, &key)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	results = append(results, foundAccounts...)

	//save entries if they have been updated
	if err := SaveAllToStore(c, results); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//going to assume there's no issue with return size - hope for gzip over the wire...
	marshalled, err := json.Marshal(results)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, string(marshalled))

}

func callMemcache(c *appengine.Context, idList []string) ([]SteamAccountDetails, []string, error) {
	return RetrieveMultiFromStore(c, idList)
}

func groupSteamIDs(idList []string) map[int][]string {
	i := 0
	j := 0

	groupList := make(map[int][]string)
	for i < len(idList) {
		var arr []string
		if i+100 > len(idList) {
			arr = idList[i:len(idList)]
		} else {
			arr = idList[i : i+100]
		}
		groupList[j] = arr

		i += 100
		j++
	}
	return groupList

}

func makeSteamAPICall(c *appengine.Context, groupedSteamIDs map[int][]string, key *string) ([]SteamAccountDetails, error) {
	mainEndpoint := "https://api.steampowered.com/ISteamUser/GetPlayerBans/v1/?key=" + *key + "&steamids="

	var steamAccounts []SteamAccountDetails
	for _, v := range groupedSteamIDs {
		endpoint := mainEndpoint + strings.Join(v, ",")

		client := urlfetch.Client(*c)
		resp, err := client.Get(endpoint)

		if err != nil {
			return nil, err
		}

		var m SteamAccount
		result, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(result, &m); err != nil {
			return nil, err
		}
		steamAccounts = append(steamAccounts, m.Players...)
	}

	return steamAccounts, nil
}

func readAPIKey() (string, error) {
	f, err := ioutil.ReadFile("key.txt")
	if err != nil {
		return "", errors.New("API Key file not found")
	}

	return string(f), nil
}
