package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"lp3/internal/data"
	"lp3/internal/mail"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type config struct {
	db   string
	port int
	env  string
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
}

type application struct {
	cfg    config
	logger *log.Logger
	model  *data.Model
	mailer *mail.Mailer
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 8000, "Port number")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development | production)")

	flag.StringVar(&cfg.db, "db", "postgres://postgres:freeroam@localhost/tg?sslmode=disable", "Database Connection String")

	flag.StringVar(&cfg.smtp.host, "smtp-host", "sandbox.smtp.mailtrap.io", "SMTP Host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 2525, "SMTP Port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", "ed7f23147046f3", "SMTP Username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", "14c3b7d8c6c7e1", "SMTP Password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "aamerasim45@gmail.com", "SMTP Sender")

	db, err := OpenDB(cfg.db)
	if err != nil {
		fmt.Printf("Could not initialize database: %v\n", err)
		return
	}
	defer db.Close()

	app := &application{}
	app.cfg = cfg
	app.logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	app.model = data.NewModel(db)
	app.mailer = mail.NewMailer(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender)

	app.logger.Println("Application succesfully Initialized")
	router := app.route()

	srv := http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  time.Minute,
	}

	// // ================================== !! =======================
	// app.logger.Println("sending mail via net/smtp")
	// subject := "Simple HTML Email"
	// body := "Here is a simple plain text"
	// sender := "gapi@mywebsite.com"
	// recipient := []string{
	// 	"a4mer.2a45@gmail.com",
	// }
	// message := fmt.Sprintf("From: %s\r\n", sender)
	// message += fmt.Sprintf("To: %s\r\n", recipient)
	// message += fmt.Sprintf("Subject: %s\r\n", subject)
	// message += fmt.Sprintf("\r\n%s\r\n", body)
	// auth := smtp.PlainAuth("", cfg.smtp.username, cfg.smtp.password, cfg.smtp.host)
	// addr := fmt.Sprintf("%s:%d", cfg.smtp.host, cfg.smtp.port)
	// if err = smtp.SendMail(addr, auth, sender, recipient, []byte(message)); err != nil {
	// 	app.logger.Printf("Unable to send mail: %v\n", err)
	// }
	// // ================================== !! =======================
	// app.logger.Println("Sent Mail")

	err = srv.ListenAndServe()
	if err != nil {
		app.logger.Fatalf("Unable to run server: %v\n", err)
	}
}

func OpenDB(db_dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", db_dsn)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
