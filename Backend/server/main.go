package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	rabbitHost := os.Getenv("RABBITMQ_HOST")
	redisHost := os.Getenv("REDIS_HOST")
	if rabbitHost == "" {
		rabbitHost = "localhost"
	}
	if redisHost == "" {
		redisHost = "localhost"
	}
	startRabbitMQ(rabbitHost)
	intializeRouter()
	stopRabbitMQ()
}

/***************
	REST API
****************/
// TODO: CRUD matches and events
// TODO: Do not use global variables but a database (REDIS)

func publishAllEvents(c *gin.Context) { // Test
	// Publish all events to RabbitMQ
	for _, evento := range eventos { // TODO: Get events from database not from global variable
		publishEvent(evento)
		time.Sleep(1 * time.Second)
	}
	c.JSON(http.StatusOK, gin.H{"message": "All events published"})
}

func createEvent(c *gin.Context) {
	var evento Event
	// Bind JSON to Event struct
	if err := c.BindJSON(&evento); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Save event in database
	evento.ID = len(eventos) + 1
	eventos = append(eventos, evento) // TODO: Remove this line when using a database
	
	// Publish event to RabbitMQ
	publishEvent(evento)
	c.JSON(http.StatusCreated, evento)
}

func getMatches(c *gin.Context) {
	c.JSON(http.StatusOK, matches)
}

func intializeRouter() {
	// Create a new router
	router := gin.Default()

	// Define the routes
	router.GET("/events", publishAllEvents)
	router.POST("/events", createEvent)

	router.GET("/matches", getMatches)

	// Run the server
	router.Run(":8080")
}


/***************
	  REDIS
****************/
// TODO: Implement DB connection
// TODO: Implement CRUD operations

/***************
    RABBITMQ
****************/

var rabbitmq RabbitMQ

func startRabbitMQ(rabbitHost string) {
	rabbitmq.conn, rabbitmq.err = amqp.Dial("amqp://guest:guest@" + rabbitHost + ":5672/")
	failOnError(rabbitmq.err, "Failed to connect to RabbitMQ")

	rabbitmq.ch, rabbitmq.err = rabbitmq.conn.Channel()
	failOnError(rabbitmq.err, "Failed to open a channel")

	rabbitmq.err = rabbitmq.ch.ExchangeDeclare(
		"match_events", // name
		"topic",        // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)
	failOnError(rabbitmq.err, "Failed to declare an exchange")

	rabbitmq.ctx, rabbitmq.cancel = context.WithTimeout(context.Background(), 5*time.Second)
}

func stopRabbitMQ() {
	rabbitmq.conn.Close()
	rabbitmq.ch.Close()
	rabbitmq.cancel()
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func publishEvent(event Event) {
	topic := "match." + strconv.Itoa(event.MatchID) + ".event." + event.Type
	body, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to marshal event: %v", err)
	}

	err = rabbitmq.ch.PublishWithContext(rabbitmq.ctx,
		"match_events", // exchange
		topic,          // routing key
		false,          // mandatory
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", body)
}

/***************
	 MODELS
****************/

type Event struct {
	ID      int    `json:"id"`
	MatchID int    `json:"matchId"`
	Team    string `json:"team"`
	Player  string `json:"player"`
	Type    string `json:"type"`
	Minute  int    `json:"minute"`
}

type Match struct {
	ID           int    `json:"id"`
	HomeTeam     string `json:"homeTeam"`
	HomeTeamAbbr string `json:"homeTeamAbbr"`
	HomeImg      string `json:"homeImg"`
	AwayTeam     string `json:"awayTeam"`
	AwayTeamAbbr string `json:"awayTeamAbbr"`
	AwayImg      string `json:"awayImg"`
	Date         string `json:"date"`
	Time         string `json:"time"`
}

type RabbitMQ struct {
	conn   *amqp.Connection
	ch     *amqp.Channel
	ctx    context.Context
	cancel context.CancelFunc
	err    error
}

// TODO: Delete this and use a database (REDIS)
var matches = []Match{
	{ID: 1, HomeTeam: "Real Madrid", HomeTeamAbbr: "RMA", HomeImg: "/team_logos/Real Madrid.png", AwayTeam: "FC Barcelona", AwayTeamAbbr: "BAR", AwayImg: "/team_logos/FC Barcelona.png", Date: "2023-10-01", Time: "20:00"},
	{ID: 2, HomeTeam: "Atlético de Madrid", HomeTeamAbbr: "ATM", HomeImg: "/team_logos/Atlético de Madrid.png", AwayTeam: "Sevilla FC", AwayTeamAbbr: "SEV", AwayImg: "/team_logos/Sevilla FC.png", Date: "2023-10-02", Time: "20:00"},
	{ID: 3, HomeTeam: "Valencia CF", HomeTeamAbbr: "VAL", HomeImg: "/team_logos/Valencia CF.png", AwayTeam: "Villarreal CF", AwayTeamAbbr: "VIL", AwayImg: "/team_logos/Villarreal CF.png", Date: "2023-10-03", Time: "22:00"},
	{ID: 4, HomeTeam: "Real Sociedad", HomeTeamAbbr: "RSO", HomeImg: "/team_logos/Real Sociedad.png", AwayTeam: "Athletic Bilbao", AwayTeamAbbr: "ATH", AwayImg: "/team_logos/Athletic Bilbao.png", Date: "2023-10-04", Time: "18:00"},
	{ID: 5, HomeTeam: "Real Betis Balompié", HomeTeamAbbr: "BET", HomeImg: "/team_logos/Real Betis Balompié.png", AwayTeam: "Deportivo Alavés", AwayTeamAbbr: "ALA", AwayImg: "/team_logos/Deportivo Alavés.png", Date: "2023-10-05", Time: "20:00"},
	{ID: 6, HomeTeam: "Celta de Vigo", HomeTeamAbbr: "CEL", HomeImg: "/team_logos/Celta de Vigo.png", AwayTeam: "RCD Espanyol Barcelona", AwayTeamAbbr: "ESP", AwayImg: "/team_logos/RCD Espanyol Barcelona.png", Date: "2023-10-06", Time: "22:00"},
}

var eventos = []Event{
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
	{ID: 22, MatchID: 2, Team: "Atlético de Madrid", Player: "Suárez", Type: "penalty", Minute: 90},
	{ID: 10, MatchID: 1, Team: "Real Madrid", Player: "Bellingham", Type: "goal", Minute: 85},
	{ID: 23, MatchID: 2, Team: "Atlético de Madrid", Player: "Oblak", Type: "start", Minute: 0},
	{ID: 11, MatchID: 1, Team: "FC Barcelona", Player: "Ter Stegen", Type: "start", Minute: 0},
	{ID: 24, MatchID: 2, Team: "Sevilla FC", Player: "Bono", Type: "half-time", Minute: 45},
	{ID: 12, MatchID: 1, Team: "Real Madrid", Player: "Courtois", Type: "half-time", Minute: 45},
	{ID: 25, MatchID: 2, Team: "Sevilla FC", Player: "Rakitic", Type: "end", Minute: 90},
	{ID: 13, MatchID: 1, Team: "Real Madrid", Player: "Modric", Type: "end", Minute: 90},
}
