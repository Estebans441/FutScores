package main

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Event struct {
	ID      int    `json:"id"`
	MatchID int    `json:"matchId"`
	Team    string `json:"team"`
	Player  string `json:"player"`
	Type    string `json:"type"`
	Minute  int    `json:"minute"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	eventos := []Event{
		{ID: 1, MatchID: 1, Team: "FC Barcelona", Player: "Gavi", Type: "goal", Minute: 15},
		{ID: 14, MatchID: 2, Team: "Atlético de Madrid", Player: "Suárez", Type: "goal", Minute: 10},
		{ID: 2, MatchID: 1, Team: "Real Madrid", Player: "Bellingham", Type: "goal", Minute: 30},
		{ID: 15, MatchID: 2, Team: "Sevilla FC", Player: "En-Nesyri", Type: "goal", Minute: 25},
		{ID: 3, MatchID: 1, Team: "FC Barcelona", Player: "Lewandowski", Type: "penalty", Minute: 35},
		{ID: 16, MatchID: 2, Team: "Atlético de Madrid", Player: "Koke", Type: "yellow card", Minute: 35},
		{ID: 4, MatchID: 1, Team: "Real Madrid", Player: "Vinicius Jr.", Type: "red card", Minute: 40},
		{ID: 17, MatchID: 2, Team: "Sevilla FC", Player: "Rakitic", Type: "substitution", Minute: 50},
		{ID: 5, MatchID: 1, Team: "FC Barcelona", Player: "Gundogan", Type: "yellow card", Minute: 50},
		{ID: 18, MatchID: 2, Team: "Atlético de Madrid", Player: "Hermoso", Type: "offside", Minute: 60},
		{ID: 6, MatchID: 1, Team: "Real Madrid", Player: "Joselu", Type: "substitution", Minute: 60},
		{ID: 19, MatchID: 2, Team: "Sevilla FC", Player: "Ocampos", Type: "corner kick", Minute: 70},
		{ID: 7, MatchID: 1, Team: "FC Barcelona", Player: "Araujo", Type: "offside", Minute: 65},
		{ID: 20, MatchID: 2, Team: "Atlético de Madrid", Player: "Llorente", Type: "goal", Minute: 80},
		{ID: 8, MatchID: 1, Team: "Real Madrid", Player: "Modric", Type: "corner kick", Minute: 75},
		{ID: 21, MatchID: 2, Team: "Sevilla FC", Player: "Navas", Type: "free kick", Minute: 85},
		{ID: 9, MatchID: 1, Team: "FC Barcelona", Player: "Ferran Torres", Type: "free kick", Minute: 80},
		{ID: 22, MatchID: 2, Team: "Atlético de Madrid", Player: "Suárez", Type: "goal", Minute: 90},
		{ID: 10, MatchID: 1, Team: "Real Madrid", Player: "Bellingham", Type: "goal", Minute: 85},
		{ID: 23, MatchID: 2, Team: "Atlético de Madrid", Player: "Oblak", Type: "start", Minute: 0},
		{ID: 11, MatchID: 1, Team: "FC Barcelona", Player: "Ter Stegen", Type: "start", Minute: 0},
		{ID: 24, MatchID: 2, Team: "Sevilla FC", Player: "Bono", Type: "half-time", Minute: 45},
		{ID: 12, MatchID: 1, Team: "Real Madrid", Player: "Courtois", Type: "half-time", Minute: 45},
		{ID: 25, MatchID: 2, Team: "Sevilla FC", Player: "Rakitic", Type: "end", Minute: 90},
		{ID: 13, MatchID: 1, Team: "Real Madrid", Player: "Modric", Type: "end", Minute: 90},
	}

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"match_events", // name
		"direct",       // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, evento := range eventos {
		body, err := json.Marshal(evento)
		if err != nil {
			log.Printf("Failed to marshal event: %v", err)
			continue
		}
		publishMessage(ctx, ch, body, strconv.Itoa(evento.MatchID))
		time.Sleep(1 * time.Second)
	}

}

func publishMessage(ctx context.Context, ch *amqp.Channel, body []byte, matchId string) {
	err := ch.PublishWithContext(ctx,
		"match_events", // exchange
		matchId,        // routing key
		false,          // mandatory
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s", body)
}
