package process

import (
	"context"
	"encoding/json"
	"github.com/pristupaanastasia/checker/config"
	"github.com/pristupaanastasia/checker/logger"
	"github.com/pristupaanastasia/checker/message"
	"github.com/pristupaanastasia/checker/telegram"
	"golang.org/x/sync/semaphore"
	"io/ioutil"
	"net/http"
)

type Checker struct {
	s         *semaphore.Weighted
	proxyList map[string]bool
	conf      *config.Config
	bot       *telegram.Bot
	log       logger.Log
}

func NewProcess(conf *config.Config, bot *telegram.Bot, log logger.Log) *Checker {
	check := &Checker{
		s:         semaphore.NewWeighted(int64(conf.Goroutine)),
		conf:      conf,
		proxyList: make(map[string]bool),
		bot:       bot,
		log:       log,
	}
	return check
}
func (ch *Checker) Process(ctx context.Context) {
	go ch.conf.Update(ctx)

	ch.Start(ctx)
}

func (ch *Checker) Start(ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			return
		default:
			if ok := ch.s.TryAcquire(1); ok {
				go ch.Parse(ctx)
			}
		}
	}
}

func (ch *Checker) Parse(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			ch.GetProxy()
		}
	}
}
func (ch *Checker) GetProxy() error {
	res, err := http.Get(ch.conf.Url)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	list := message.ProxyList{}
	err = json.Unmarshal(body, &list)
	if err != nil {
		return err
	}

	for _, proxy := range list.Proxies {
		if result := ch.CheckProxy(proxy.Ip); result.Result == message.SUCCESS {
			ch.proxyList[proxy.Ip] = true
		}
	}

	for ip, _ := range ch.proxyList {
		if result := ch.CheckProxy(ip); result.Result == message.SUCCESS {
			ch.proxyList[ip] = true
		}
	}
	return err
}

func (ch *Checker) CheckProxy(url string) message.ResultJson {
	_, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return message.ResultJson{Result: message.FAILED}
	}
	ch.bot.Send(message.TelegramResult{Url: url, Result: message.SUCCESS})
	return message.ResultJson{Result: message.SUCCESS}
}
