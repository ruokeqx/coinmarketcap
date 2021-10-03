This is `coinmarketcap` spider and webserver Golang edition.


## Usage:

If historical data is required for the first run, you need run spider for some time before running server.

```bash
# To run spider, you need to uncomment the main function in the spider.go
go run spider.go sqlutil.go
```

```bash
# To run server, you need to comment out the main function in the spider.go
go run main.go router.go sqlutil.go spider.go
```



