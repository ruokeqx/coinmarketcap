package main

// func errprint(msg string, err error) {
// 	if err != nil {
// 		log.Fatalf("%s: %s", msg, err)
// 	}
// }

// func main() {
// 	conn, err := amqp.Dial("amqp://guest:guest@192.168.0.175:5672/")
// 	errprint("Failed to connect to RabbitMQ", err)
// 	defer conn.Close()

// 	ch, err := conn.Channel()
// 	errprint("Failed to open a channel", err)
// 	defer ch.Close()

// 	q, err := ch.QueueDeclare(
// 		"dspider", // name
// 		false,     // durable
// 		true,      // delete when unused
// 		false,     // exclusive
// 		false,     // no-wait
// 		nil,       // arguments
// 	)
// 	errprint("Failed to declare a queue", err)

// 	// 流量控制
// 	ch.Qos(32, 0, false)

// 	msgs, err := ch.Consume(
// 		q.Name, // queue
// 		"",     // consumer
// 		false,  // auto-ack
// 		false,  // exclusive
// 		false,  // no-local
// 		false,  // no-wait
// 		nil,    // args
// 	)
// 	errprint("Failed to register a consumer", err)

// 	// 创建数据库连接
// 	db := sqlInit()

// 	// 允许三个并发
// 	s := semaphore.NewWeighted(3)
// 	for {
// 		for msg := range msgs {
// 			go func(msg amqp.Delivery) {
// 				// 获取信号量
// 				s.Acquire(context.Background(), 1)
// 				defer s.Release(1)
// 				// 消息解析
// 				splitList := strings.Split(string(msg.Body), " ")
// 				coin_name := splitList[0]
// 				choice := splitList[1]
// 				tS, _ := strconv.ParseInt(splitList[2], 10, 64)
// 				flag, _ := strconv.ParseBool(splitList[3])
// 				// fmt.Println(coin_name, choice, tS, flag)

// 				// 爬虫
// 				tc := Coin{}
// 				id := 0
// 				db.AutoMigrate(&Coin{})
// 				if tc.Id == 0 {
// 					url := fmt.Sprintf(market_url, coin_name)
// 					id = GetId(Download(url))
// 				}
// 				InsertCoin(db, coin_name, id)

// 				if flag || !db.HasTable("chart-"+coin_name) {
// 					ParserChartData(db, coin_name, chart_url, id, choice)
// 				}

// 				if flag || !db.HasTable("history-"+coin_name) {
// 					GetHistoryData(db, coin_name, historical_url, id, tS)
// 				}

// 				// 手动确认 保证消息不丢失
// 				msg.Ack(false)
// 			}(msg)
// 		}
// 	}
// }
