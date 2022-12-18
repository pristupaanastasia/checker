# Checker proxy url

## API

| Url             | Method   | Value                                               | Description                                                                       |
|-----------------|----------|-----------------------------------------------------|-----------------------------------------------------------------------------------|
| /status         | GET      | {"status":string}                                   | return status process: "enabled"/"stopped"/"error"                                |
| /thread_count   | GET/POST | {"thread":int}                                      | return and change thread count                                                    |
| /proxy_list_url | GET/POST | {"proxy":string}                                    | return and change url with info about proxy                                       |
| /proxy_list     | GET/POST | {"proxy-list":[]{"url":string,"result":string}}     | return and change list of proxy                                                   |
| /telegram_id    | GET/POST | {"id-telegram":int}                                 | return and change telegram id chat                                                |
| /telegram_token | GET/POST | {"token-tg":string}                                 | return and change token of telegram bot                                           |
| /logs           | GET      | file with logs                                      | return file with logs with date                                                   |
| /stats          | GET      | {"success":float,"speed":int,"average-speed":float} | return stat about count of success check proxy, speed and average-speed all check |
| /stats_clear    | POST     |                                                     | clear stat                                                                        |

## Config

program handle 2 way to set env: from config.json and ENV

| Name json    | Name ENV       | Value                                | Description                   |
|--------------|----------------|--------------------------------------|-------------------------------|
| start-parser | START_PARSER   | "on"/"off"                           | start or stop parse proxy url |
| goroutine    | GOROUTINE      | int                                  | count of goroutine            |
| url          | URL            | string                               | url with proxy list           |
| id-telegram  | ID_TELEGRAM    | int                                  | id chat telegram              |
| token-tg     | TOKEN_TELEGRAM | string                               | token of telegram bot         |
| port         | PORT           | string                               | port of server                |
|              | LOGLEVEL       | "DEBUG", "ALL","ERROR","WARN","INFO" | log level                     |


## Start app

1.  docker-compose build
    docker-compose up
2. go run main.go

