package model

import (
	"encoding/json"
	"github.com/pristupaanastasia/checker/config"
	"github.com/pristupaanastasia/checker/logger"
	"github.com/pristupaanastasia/checker/message"
	"github.com/pristupaanastasia/checker/process"
	"github.com/pristupaanastasia/checker/telegram"
	"io/ioutil"
	"net/http"
)

type Model struct {
	conf    *config.Config
	checker *process.Checker
	bot     *telegram.Bot
	log     logger.Log
}

func NewModel(conf *config.Config, checker *process.Checker, bot *telegram.Bot, log logger.Log) *Model {
	return &Model{
		conf:    conf,
		checker: checker,
		bot:     bot,
		log:     log,
	}
}
func (m *Model) Status(w http.ResponseWriter, r *http.Request) {

}

func (m *Model) ProxyListUrl(w http.ResponseWriter, r *http.Request) {

}

func (m *Model) ProxyList(w http.ResponseWriter, r *http.Request) {

}

func (m *Model) ThreadCount(w http.ResponseWriter, r *http.Request) {

}
func (m *Model) CheckUrl(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		url := message.Proxy{Ip: m.conf.Url}

		resp, err := json.Marshal(url)
		if err != nil {
			m.log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	case http.MethodPost:
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			m.log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		mes := message.Proxy{}
		if err = json.Unmarshal(body, &mes); err != nil {
			m.log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		m.conf.Url = mes.Ip
		w.WriteHeader(http.StatusOK)
	}
}
func (m *Model) TelegramId(w http.ResponseWriter, r *http.Request) {

}
func (m *Model) TelegramToken(w http.ResponseWriter, r *http.Request) {

}
func (m *Model) Logs(w http.ResponseWriter, r *http.Request) {

}
func (m *Model) Stats(w http.ResponseWriter, r *http.Request) {

}
func (m *Model) StatsClear(w http.ResponseWriter, r *http.Request) {

}
