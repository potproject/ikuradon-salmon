package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/potproject/ikuradon-salmon/dataaccess"
	"github.com/potproject/ikuradon-salmon/endpoints"
	"github.com/potproject/ikuradon-salmon/setting"
)

func main() {
	godotenv.Load()
	setting.SetSetting()
	list := flag.Bool("list", false, "show Database")
	flag.Parse()

	if setting.S.UseRedis {
		dataaccess.SetRedis()
	} else {
		dataaccess.SetLevel("data/level")
	}

	if *list {
		p, err := dataaccess.DA.ListAll()
		if err != nil {
			os.Exit(1)
		}
		for _, v := range p {
			fmt.Printf("%#v\n", v)
		}
		return
	}

	router := mux.NewRouter()
	// エンドポイント
	router.HandleFunc("/", endpoints.Healthcheck).Methods("GET")
	router.HandleFunc("/health-check", endpoints.Healthcheck).Methods("GET")
	router.HandleFunc("/api/v1/id", endpoints.GetID).Methods("GET")
	router.HandleFunc("/api/v2/subscribe", endpoints.PostSubscribe).Methods("POST")
	router.HandleFunc("/api/v2/unsubscribe", endpoints.PostUnsubscribe).Methods("POST")
	router.HandleFunc("/api/v1/webpush/{subscribeID}", endpoints.PostWebPush).Methods("POST")
	domain := fmt.Sprintf("%s:%d", setting.S.APIHost, setting.S.APIPort)
	realURL := "http://" + domain + "/"
	if setting.S.BaseURL != "" {
		realURL = setting.S.BaseURL
	}
	fmt.Println("Running to " + realURL)
	log.Fatal(http.ListenAndServe(domain, router))
}
