package main

import (
	"context"
	"github.com/gorilla/mux"
	config "github.com/pristupaanastasia/checker/config"
	"github.com/pristupaanastasia/checker/logger"
	"github.com/pristupaanastasia/checker/model"
	"github.com/pristupaanastasia/checker/process"
	"github.com/pristupaanastasia/checker/telegram"
	"net/http"

	"os"
)

func main() {
	level := os.Getenv("LOGLEVEL")
	logfile := os.Getenv("LOGFILE")
	log, err := logger.New(level, logfile)
	if err != nil {
		log.Fatal(err)
	}
	defer log.Delete()
	r := mux.NewRouter()

	conf, err := config.NewConfig(log)
	if err != nil {
		log.Fatal(err)
	}
	bot, err := telegram.NewTg(conf, log)
	if err != nil {
		log.Fatal()
	}
	ctx := context.Background()
	checker := process.NewProcess(conf, bot, log)
	go checker.Process(ctx)
	handler := model.NewModel(conf, checker, bot, log)
	r.HandleFunc("/status", handler.Status)
	r.HandleFunc("/thread_count", handler.ThreadCount)
	r.HandleFunc("/proxy_list_url", handler.ProxyListUrl)
	r.HandleFunc("/proxy_list", handler.ProxyList)
	r.HandleFunc("/check_url", handler.CheckUrl)
	r.HandleFunc("/telegram_id", handler.TelegramId)
	r.HandleFunc("/telegram_token", handler.TelegramToken)
	r.HandleFunc("/logs", handler.Logs)
	r.HandleFunc("/stats", handler.Stats)
	r.HandleFunc("/stats_clear", handler.StatsClear).Methods("POST")

	log.Info("http listening on", conf.Port)
	if err = http.ListenAndServe(conf.Port, r); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
