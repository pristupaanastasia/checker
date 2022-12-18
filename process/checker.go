package process

import (
	"checker/config"
	"checker/logger"
	"checker/message"
	"checker/stat"
	"checker/telegram"
	"context"
	"encoding/json"
	"golang.org/x/sync/semaphore"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
)

type Checker struct {
	s         *semaphore.Weighted
	mu        sync.Mutex
	proxyList map[string]bool

	conf   *config.Config
	bot    *telegram.Bot
	stat   *stat.Stat
	log    logger.Log
	status string
}

const (
	ENABLED = "enabled"
	STOPPED = "stopped"
	ERROR   = "error"
)
const (
	ON  = "on"
	OFF = "off"
)

func NewProcess(conf *config.Config, bot *telegram.Bot, stat *stat.Stat, log logger.Log) *Checker {
	log.Info("Get process config")
	check := &Checker{
		stat:      stat,
		s:         semaphore.NewWeighted(int64(conf.Goroutine)),
		conf:      conf,
		proxyList: make(map[string]bool),
		bot:       bot,
		log:       log,
		status:    ENABLED,
		mu:        sync.Mutex{},
	}
	return check
}
func (ch *Checker) Process(ctx context.Context) {
	go ch.conf.Update(ctx)

	ch.log.Info("Start process")
	ch.Start(ctx)
}

func (ch *Checker) Start(ctx context.Context) {

	var ipList chan message.Proxy
	var ip chan string
	var deleteip chan string

	deleteip = make(chan string)
	ip = make(chan string, ch.conf.Goroutine)
	ipList = make(chan message.Proxy, ch.conf.Goroutine*2)
	go ch.SendUrl(ctx, ipList)
	for {
		select {
		case <-ctx.Done():
			return
		case proxy := <-ip:
			ch.AddProxyUrl(proxy)
		case proxy := <-deleteip:
			ch.DeleteProxyUrl(proxy)
		default:
			if ok := ch.s.TryAcquire(1); ok && ch.conf.StartParser == ON {
				ch.log.Info("Start goroutine")
				go ch.GetProxy(ctx, ip, deleteip, ipList)

			}
		}
	}
}
func (ch *Checker) SendUrl(ctx context.Context, ip chan message.Proxy) {
	var listProxy []message.Proxy
	for {
		res, err := http.Get(ch.conf.Url)
		if err != nil {
			ch.log.Error(err)
			ch.status = ERROR
			continue
		}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			ch.log.Error(err)
			ch.status = ERROR
			continue
		}
		ch.log.Info("GET info")
		var list message.ProxyList
		err = json.Unmarshal(body, &list)
		if err != nil {
			ch.log.Error(err)
			ch.status = ERROR
			continue
		}
		listProxy = list.Proxies

		for _, proxy := range listProxy {
			select {
			case <-ctx.Done():
				return
			default:
				ip <- proxy

			}
		}
		listProxy = []message.Proxy{}
	}

}

func (ch *Checker) GetProxy(ctx context.Context, ip chan string, deleteip chan string, list chan message.Proxy) {

	ch.status = ENABLED
	for {
		select {
		case <-ctx.Done():
			return
		case proxy := <-list:
			existProxy := ch.ProxyExist(proxy.Ip)
			if result := ch.CheckProxy(proxy.Ip); result.Result == message.SUCCESS {
				if !existProxy {
					go ch.bot.Send(message.ProxyResult{Url: proxy.Ip, Result: message.SUCCESS})
					ch.log.Info(message.ProxyResult{Url: proxy.Ip, Result: message.SUCCESS})
					if err := ch.WriteSuccess(message.ProxyResult{Url: proxy.Ip, Result: message.SUCCESS}); err != nil {
						ch.log.Error(err)
					}
				}
				ch.stat.Success()
				ip <- proxy.Ip

			} else if existProxy && result.Result == message.FAILED {
				deleteip <- proxy.Ip
			}
			ch.stat.Speed()
			if ch.conf.StartParser == OFF {
				ch.status = STOPPED
				return
			}
		}
	}

	//ch.mu.Lock()
	//for ip, _ := range ch.proxyList {
	//	if result := ch.CheckProxy(ip); result.Result == message.SUCCESS {
	//		ch.proxyList[ip] = true
	//	}
	//}
	//ch.mu.Unlock()

}

func (ch *Checker) CheckProxy(url string) message.ResultJson {
	_, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return message.ResultJson{Result: message.FAILED}
	}

	return message.ResultJson{Result: message.SUCCESS}
}

func (ch *Checker) GetListProxy() []message.ProxyResult {

	list := make([]message.ProxyResult, 0)
	ch.mu.Lock()
	for ip, _ := range ch.proxyList {
		list = append(list, message.ProxyResult{Url: ip, Result: message.SUCCESS})
	}
	ch.mu.Unlock()
	return list
}

func (ch *Checker) AddProxyUrl(ip string) {
	ch.mu.Lock()
	ch.proxyList[ip] = true
	ch.mu.Unlock()
}
func (ch *Checker) DeleteProxyUrl(ip string) {
	ch.mu.Lock()
	delete(ch.proxyList, ip)
	ch.mu.Unlock()
}
func (ch *Checker) ProxyExist(ip string) bool {
	ch.mu.Lock()
	_, ok := ch.proxyList[ip]
	ch.mu.Unlock()
	return ok
}
func (ch *Checker) UpdateListProxy(list []message.ProxyResult) {
	ch.mu.Lock()
	for _, ip := range list {
		ch.proxyList[ip.Url] = true
	}
	ch.mu.Unlock()
}

func (ch *Checker) GetStatus() string {
	return ch.status
}

func (ch *Checker) WriteSuccess(v interface{}) error {
	var err error
	resp, err := json.Marshal(v)
	if err != nil {
		return err
	}

	resp = append(resp, []byte("\n")...)
	var filename = "success.txt"
	var f *os.File

	if checkFileIsExist(filename) { // Если файл существует
		f, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0666) // Открыть файл
	} else {
		f, err = os.Create(filename) // Создать файл
	}
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(resp) // Запись в файл (массив байтов)
	if err != nil {
		return err
	}
	err = f.Sync()
	return err
}
func checkFileIsExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}
