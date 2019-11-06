package coordinator

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/ericxiao417/sensor-farm/dto"
	"github.com/ericxiao417/sensor-farm/qutils"
	"github.com/streadway/amqp"
)

const url = "amqp://guest:guest@localhost:5672"

type QueueListener struct {
	conn    *amqp.Connection
	ch      *amqp.Channel
	sources map[string]<-chan amqp.Delivery
	ea      *EventAggregator
}

func NewQueueListener(ea *EventAggregator) *QueueListener {
	ql := QueueListener{
		sources: make(map[string]<-chan amqp.Delivery),
		ea:      ea,
	}

	ql.conn, ql.ch = qutils.GetChannel(url)

	return &ql
}

func (ql *QueueListener) DiscoverSensors() {
	ql.ch.ExchangeDeclare(
		qutils.SensorDiscoveryExchange, //name string,
		"fanout",                       //kind string,
		false,                          //durable bool,
		false,                          //autoDelete bool,
		false,                          //internal bool,
		false,                          //noWait bool,
		nil)                            //args amqp.Table)

	ql.ch.Publish(
		qutils.SensorDiscoveryExchange, //exchange string,
		"",                             //key string,
		false,                          //mandatory bool,
		false,                          //immediate bool,
		amqp.Publishing{})              //msg amqp.Publishing)
}

func (ql *QueueListener) ListenForNewSource() {
	q := qutils.GetQueue("", ql.ch, true)
	ql.ch.QueueBind(
		q.Name,       //name string,
		"",           //key string,
		"amq.fanout", //exchange string,
		false,        //noWait bool,
		nil)          //args amqp.Table)

	msgs, _ := ql.ch.Consume(
		q.Name, //queue string,
		"",     //consumer string,
		true,   //autoAck bool,
		false,  //exclusive bool,
		false,  //noLocal bool,
		false,  //noWait bool,
		nil)    //args amqp.Table)

	ql.DiscoverSensors()

	fmt.Println("listening for new sources")
	for msg := range msgs {
		if ql.sources[string(msg.Body)] == nil {
			fmt.Println("new source discovered")
			ql.ea.PublishEvent("DataSourceDiscovered", string(msg.Body))
			sourceChan, _ := ql.ch.Consume(
				string(msg.Body), //queue string,
				"",               //consumer string,
				true,             //autoAck bool,
				false,            //exclusive bool,
				false,            //noLocal bool,
				false,            //noWait bool,
				nil)              //args amqp.Table)

			ql.sources[string(msg.Body)] = sourceChan

			go ql.AddListener(sourceChan)
		}
	}
}

func (ql *QueueListener) AddListener(msgs <-chan amqp.Delivery) {
	for msg := range msgs {
		r := bytes.NewReader(msg.Body)
		d := gob.NewDecoder(r)
		sd := new(dto.SensorMessage)
		d.Decode(sd)

		fmt.Printf("Received message: %v\n", sd)

		ed := EventData{
			Name:      sd.Name,
			Timestamp: sd.Timestamp,
			Value:     sd.Value,
		}

		ql.ea.PublishEvent("MessageReceived_"+msg.RoutingKey, ed)
	}
}
