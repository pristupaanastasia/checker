package model

import (
	"checker/config"
	"checker/logger"
	"checker/message"
	"checker/process"
	"checker/stat"
	"checker/telegram"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Model struct {
	conf    *config.Config
	checker *process.Checker
	bot     *telegram.Bot
	stat    *stat.Stat
	log     logger.Log
}

func NewModel(conf *config.Config, checker *process.Checker, bot *telegram.Bot, stat *stat.Stat, log logger.Log) *Model {
	log.Info("Start model handler")
	return &Model{
		conf:    conf,
		checker: checker,
		bot:     bot,
		stat:    stat,
		log:     log,
	}
}
func (m *Model) Status(w http.ResponseWriter, r *http.Request) {
	proxy := map[string]string{"status": m.checker.GetStatus()}

	m.GetInfo(w, r, proxy)
}

func (m *Model) ProxyListUrl(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		proxy := map[string]string{"proxy": m.conf.Url}

		m.GetInfo(w, r, proxy)
	case http.MethodPost:
		proxy := map[string]string{}
		if err := m.PostInfo(w, r, &proxy); err != nil {
			m.log.Error(err)
			return
		}

		m.conf.Url = proxy["proxy"]
		w.WriteHeader(http.StatusOK)
	}
}

func (m *Model) ProxyList(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		list := map[string][]message.ProxyResult{"proxy-list": m.checker.GetListProxy()}

		m.GetInfo(w, r, list)
	case http.MethodPost:
		list := map[string][]message.Proxy{}
		if err := m.PostInfo(w, r, &list); err != nil {
			m.log.Error(err)
			return
		}

		m.checker.UpdateListProxy(list["proxy-list"])
		w.WriteHeader(http.StatusOK)
	}
}

func (m *Model) ThreadCount(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		thread := map[string]int{"thread": m.conf.Goroutine}

		m.GetInfo(w, r, &thread)
	case http.MethodPost:
		thread := map[string]int{}
		if err := m.PostInfo(w, r, thread); err != nil {
			m.log.Error(err)
			return
		}

		m.conf.Goroutine = thread["thread"]
		w.WriteHeader(http.StatusOK)
	}
}

func (m *Model) TelegramId(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		id := map[string]int64{"id-telegram": m.conf.IdTelegram}

		m.GetInfo(w, r, id)
	case http.MethodPost:
		id := map[string]int64{}
		if err := m.PostInfo(w, r, &id); err != nil {
			m.log.Error(err)
			return
		}

		m.conf.IdTelegram = id["id-telegram"]
		w.WriteHeader(http.StatusOK)
	}
}
func (m *Model) TelegramToken(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		token := map[string]string{"token-tg": m.conf.TokenTg}

		m.GetInfo(w, r, token)
	case http.MethodPost:
		token := map[string]int64{}
		if err := m.PostInfo(w, r, &token); err != nil {
			m.log.Error(err)
			return
		}

		m.conf.IdTelegram = token["token-tg"]
		if err := m.bot.ReConnect(); err != nil {
			m.log.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
func (m *Model) Logs(w http.ResponseWriter, r *http.Request) {

	file, err := ioutil.ReadFile("info.txt")
	if err != nil {
		m.log.Error(err)
	} else {
		w.Write(file)
		w.WriteHeader(http.StatusOK)
	}
}

func (m *Model) Stats(w http.ResponseWriter, r *http.Request) {
	statJson := m.stat.GetStat()
	m.log.Info("STAT", statJson)
	m.GetInfo(w, r, statJson)
}
func (m *Model) StatsClear(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		m.stat.Clear()
	}
}

func (m *Model) GetInfo(w http.ResponseWriter, r *http.Request, v interface{}) {
	resp, err := json.Marshal(v)
	if err != nil {
		m.log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
func (m *Model) PostInfo(w http.ResponseWriter, r *http.Request, pointer interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		m.log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	if err = json.Unmarshal(body, pointer); err != nil {
		m.log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return err
	}
	return err
}
