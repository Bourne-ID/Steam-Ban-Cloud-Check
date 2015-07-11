package steamapi

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"appengine"
	"appengine/urlfetch"
)

func init() {
	http.HandleFunc("/", root)
	//http.HandleFunc("/sign", sign)
}

func root(w http.ResponseWriter, r *http.Request) {
	steamIdsString := r.URL.Query().Get("steamids")
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

	c := appengine.NewContext(r)

	results, err := makeSteamAPICall(&c, &steamIDArray, &key)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(results))

}

func makeSteamAPICall(c *appengine.Context, id *[]string, key *string) ([]byte, error) {
	endpoint := "https://api.steampowered.com/ISteamUser/GetPlayerBans/v1/?key=" + *key + "&steamids="
	endpoint += strings.Join(*id, ",")

	client := urlfetch.Client(*c)
	resp, err := client.Get(endpoint)

	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(resp.Body)
}

func readAPIKey() (string, error) {
	f, err := ioutil.ReadFile("key.txt")
	if err != nil {
		return "", errors.New("API Key file not found")
	}

	return string(f), nil
}
