package main

import (
	"CRUD-DATABASE/handler"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form"
	"github.com/jmoiron/sqlx"
	"github.com/justinas/nosurf"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

var sessionManager *scs.SessionManager

func main() {
	tableMigration := `
    CREATE TABLE IF NOT EXISTS students (
		id BIGSERIAL,
        name TEXT NOT NULL,
		email TEXT NOT NULL,
		roll INT NOT NULL,
		english INT NOT NULL,
		bangla INT NOT NULL,
		mathematics INT NOT NULL,
		grade TEXT,
		gpa FLOAT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMP DEFAULT NULL,

		PRIMARY KEY(id),
		UNIQUE(email)
    );`

	adminTableMigration := `
    CREATE TABLE IF NOT EXISTS admin (
		id BIGSERIAL,
        username TEXT NOT NULL,
		password TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

		PRIMARY KEY(id),
		UNIQUE(username)

    );`

	config := viper.NewWithOptions(
		viper.EnvKeyReplacer(
			strings.NewReplacer(".", "_"),
		),
	)
	config.SetConfigFile("env/config")
	config.SetConfigType("ini")
	config.AutomaticEnv()

	if err := config.ReadInConfig(); err != nil {
		log.Fatalf("error loading configuration: %v", err)
	}

	lt := config.GetDuration("session.lifetime")
	it := config.GetDuration("session.idletime")
	sessionManager = scs.New()
	sessionManager.Lifetime = lt * time.Hour
	sessionManager.IdleTimeout = it * time.Minute
	sessionManager.Cookie.Name = "web-session"
	sessionManager.Cookie.HttpOnly = true
	sessionManager.Cookie.Secure = true

	db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.GetString("database.host"),
		config.GetString("database.port"),
		config.GetString("database.user"),
		config.GetString("database.password"),
		config.GetString("database.dbname"),
		config.GetString("database.sslmode"),
	))
	if err != nil {
		log.Fatalln(err)
	}

	res := db.MustExec(tableMigration)
	row, err := res.RowsAffected()
	if err != nil {
		log.Fatalln(err)
	}

	if row < 0 {
		log.Fatalln("failed to run schema")
	}
	
	resAd := db.MustExec(adminTableMigration)
	rowAd, err := resAd.RowsAffected()
	if err != nil {
		log.Fatalln(err)
	}

	if rowAd < 0 {
		log.Fatalln("failed to run schema")
	}

	decoder := form.NewDecoder()
	chi := handler.NewHandler(sessionManager, decoder, db)
	p := config.GetInt("server.port")

	newChi := nosurf.New(chi)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", p), sessionManager.LoadAndSave(newChi)); err != nil {
		log.Fatalf("%#v", err)
	}
}
