xpackage main // Имя текущего пакета
import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

// Функция main как точка входа

type RabbitCFG struct {
	uri          string
	exchange     string
	exchangeType string
	queue        string
	bindingKey   string
	consumerTag  string
}
type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	tag     string
	done    chan error
}

/*	uri          = flag.String("uri", "amqp://guest:guest@localhost:5672/", "AMQP URI")
	exchange     = flag.String("exchange", "test-exchange", "Durable, non-auto-deleted AMQP exchange name")
	exchangeType = flag.String("exchange-type", "direct", "Exchange type - direct|fanout|topic|x-custom")
	queue        = flag.String("queue", "test-queue", "Ephemeral AMQP queue name")
	bindingKey   = flag.String("key", "test-key", "AMQP binding key")
	consumerTag  = flag.String("consumer-tag", "simple-consumer", "AMQP consumer tag (should not be blank)")
	lifetime     = flag.Duration("lifetime", 5*time.Second, "lifetime of process before shutdown (0s=infinite)")
*/

func main() {
	q := RabbitCFG{"amqp://guest:guest@localhost:5672/", "test-exchange2", "direct", "test-queue", "test-key", "simple-consumer"}
	//	a := 1
	fmt.Println("Hello!", q)

	c := &Consumer{
		conn:    nil,
		channel: nil,
		tag:     q.consumerTag,
		done:    make(chan error),
	}

	var err error

	log.Printf("dialing %q", q.uri)
	c.conn, err = amqp.Dial(q.uri)
	if err != nil {
		fmt.Errorf("Dial: %s", err)
	}
	log.Printf("c.con=", c.conn)

	log.Printf("got Connection, getting Channel")
	c.channel, err = c.conn.Channel()
	if err != nil {
		fmt.Errorf("Channel: %s", err)
	}
	log.Printf("c.channel=", c.channel)
	if err = c.channel.ExchangeDeclare(
		q.exchange,     // name of the exchange
		q.exchangeType, // type
		true,           // durable
		false,          // delete when complete
		false,          // internal
		false,          // noWait
		nil,            // arguments
	); err != nil {
		fmt.Errorf("Exchange Declare: %s", err)
	}

	log.Printf("declared Exchange, declaring Queue %q", q.queue)
	queue, err := c.channel.QueueDeclare(
		q.queue, // name of the queue
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // noWait
		nil,     // arguments
	)
	if err != nil {
		fmt.Errorf("Queue Declare: %s", err)
	}

	log.Printf("queue=", queue)

	if err = c.channel.QueueBind(
		q.queue,      // name of the queue
		q.bindingKey, // bindingKey
		q.exchange,   // sourceExchange
		false,        // noWait
		nil,          // arguments
	); err != nil {
		fmt.Errorf("Queue Bind: %s", err)
	}

	if err = c.channel.Publish(
		q.exchange,   // publish to an exchange
		q.bindingKey, // routing to 0 or more queues
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            []byte("loko"),
			DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,              // 0-9
			// a bunch of application/implementation-specific fields
		},
	); err != nil {
		fmt.Errorf("Exchange Publish: %s", err)
	}

}
