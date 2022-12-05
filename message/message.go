package message

type Proxy struct {
	GoogleError     string `json:"google_error"`
	GoogleStatus    int    `json:"google_status"`
	GoogleTotalTime int    `json:"google_total_time"`
	Ip              string `json:"ip"`
	Port            string `json:"port"`
	YahooError      string `json:"yahoo_error"`
	YahooStatus     int    `json:"yahoo_status"`
	YahooTotalTime  int    `json:"yahoo_total_time"`
	YandexError     string `json:"yandex_error"`
	YandexStatus    int    `json:"yandex_status"`
	YandexTotalTime int    `json:"yandex_total_time"`
}

type ProxyList struct {
	Date    string  `json:"date"`
	Proxies []Proxy `json:"proxies"`
}
type ResultJson struct {
	Result string `json:"result"`
}

type ProxyResult struct {
	Url    string `json:"url"`
	Result string `json:"result"`
}
type Stat struct {
	Success      float64 `json:"success"`
	Speed        int64   `json:"speed"`
	AverageSpeed float64 `json:"average-speed"`
}

const SUCCESS = "success"
const FAILED = "failed"
