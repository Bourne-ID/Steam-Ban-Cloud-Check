package steamapi

import (
	"fmt"
	"net/http"
	"time"
)

func init() {
	http.HandleFunc("/", root)
	//http.HandleFunc("/sign", sign)
}

func root(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello, the time is now ", time.Now())
}
