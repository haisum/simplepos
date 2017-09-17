package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/haisum/simplepos/db"
	"github.com/haisum/simplepos/request/items"
	"github.com/spf13/viper"
	"net/http"
	"os"
)

func init() {
	viper.SetDefault("HTTP_PORT", 8443)
	viper.SetDefault("HTTP_HOST", "")
	viper.SetDefault("HTTP_CERTFILE", "keys/cert.pem")
	viper.SetDefault("HTTP_KEYFILE", "keys/server.key")
	viper.SetDefault("DB_DRIVER", "sqlite3")
	viper.SetDefault("DSN", "simplepos.sqlite3")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.WithError(err).Warn("Couldn't read config file.")
	}
	viper.SetEnvPrefix("POS")
	viper.AutomaticEnv()
}

func main() {
	fmt.Println("Welcome to simple POS")

	log.Infof("Connecting database with driver %s and dsn %s.", viper.GetString("DB_DRIVER"), viper.GetString("DSN"))
	conn, err := db.Connect(viper.GetString("DB_DRIVER"), viper.GetString("DSN"))
	if err != nil {
		fmt.Printf("%+v", err)
		os.Exit(1)
	}
	defer conn.Db.Close()
	_, err = db.CheckAndCreateDB()
	if err != nil {
		fmt.Printf("%+v", err)
		os.Exit(1)
	}

	r := mux.NewRouter()
	//Item handlers
	r.Handle("/items", handlers.MethodHandler{"GET": items.List, "DELETE": items.Delete, "POST": items.Add, "PUT": items.Update})

	static := http.FileServer(http.Dir("./static/"))
	r.PathPrefix("/").Handler(static)

	log.Infof("Listening for requests on port %d", viper.GetInt("HTTP_PORT"))
	log.Fatal(http.ListenAndServeTLS(fmt.Sprintf("%s:%d", viper.GetString("HTTP_HOST"), viper.GetInt("HTTP_PORT")), viper.GetString("HTTP_CERTFILE"), viper.GetString("HTTP_KEYFILE"), handlers.CombinedLoggingHandler(os.Stdout, r)))
}
