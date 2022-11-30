package message

type Proxy struct {
	Ip string `json:"ip"`
}

type ProxyList struct {
	Proxies []Proxy `json:"proxies"`
}
type ResultJson struct {
	Result string `json:"result"`
}

type TelegramResult struct {
	Url    string `json:"url"`
	Result string `json:"result"`
}

const SUCCESS = "success"
const FAILED = "failed"
