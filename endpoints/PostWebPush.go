package endpoints

import (
	"bytes"
	"fmt"
	"net/http"
)

func PostWebPush(w http.ResponseWriter, r *http.Request) {
	fmt.Println("RemoteAddr:", r.RemoteAddr)
	fmt.Println("Path      :", r.URL.Path)
	fmt.Println("Query     :", r.URL.RawQuery)
	// Header all
	for name, values := range r.Header {
		// Loop over all values for the name.
		for _, value := range values {
			fmt.Println("Header    :", name, value)
		}
	}
	// Body Get
	bufbody := new(bytes.Buffer)
	bufbody.ReadFrom(r.Body)
	body := bufbody.String()
	fmt.Println("Body     :", body)
}
