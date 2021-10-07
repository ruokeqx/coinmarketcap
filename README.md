This is `coinmarketcap` distributed spider and webserver Golang edition.


## Usage:

You need a rabbitmq as the broker. Web server itself run as a publisher and spider run infinitly as consumer.

```bash
# To run spider
go run dspider.go spider.go sqlutil.go
```

```bash
# To run server
go run main.go router.go sqlutil.go spider.go
```
