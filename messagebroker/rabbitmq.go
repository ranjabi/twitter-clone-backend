package messagebroker

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMq struct {
	ctx  context.Context
	conn *amqp.Connection
	ch   *amqp.Channel
	q    amqp.Queue
}

func NewRabbitMq(ctx context.Context, conn *amqp.Connection, ch *amqp.Channel) RabbitMq {
	return RabbitMq{ctx: ctx, conn: conn, ch: ch}
}

func (r *RabbitMq) DeclareFeedQeueu() error {
	var err error
	r.q, err = r.ch.QueueDeclare(
		"feed", // name
		false,  // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *RabbitMq) SendFeed(tweetId int, userId int) error {
	body := fmt.Sprintf("%d,%d", tweetId, userId)
	err := r.ch.PublishWithContext(r.ctx,
		"",       // exchange
		r.q.Name, // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		return err
	}
	return nil
}

// TODO continue with receive
