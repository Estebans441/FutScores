package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client
var ctx = context.Background()

func main() {
	rabbitHost := os.Getenv("RABBITMQ_HOST")
	redisHost := os.Getenv("REDIS_HOST")
	if rabbitHost == "" {
		rabbitHost = "localhost"
	}
	if redisHost == "" {
		redisHost = "localhost"
	}

	// Connect to Redis
	connectToRedis(redisHost)

	// Initialize Redis data with sample matches
	initializeSampleMatches()

	// Start RabbitMQ
	startRabbitMQ(rabbitHost)
	initializeRouter()
	stopRabbitMQ()
}

/***************
	ROUTER FUNC - API
****************/

func initializeRouter() {
	router := gin.Default()

	router.Use(cors.Default())

	// Routes for CRUD operations on matches
	router.POST("/matches", createMatch)
	router.GET("/matches/:id", getMatch)
	router.GET("/matches", getAllMatches)
	router.PUT("/matches/:id", updateMatch)
	router.DELETE("/matches/:id", deleteMatch)
	router.GET("/matches/:id/events", getAllEvents)

	// Routes for CRUD operations on events
	router.POST("/events", createEvent)
	router.GET("/events/:id", getEvent)
	router.DELETE("/events/:id", deleteEvent)
	router.PUT("/events/:id", updateEvent)

	router.Run(":8080") // Start server on port 8080
}

// Function to initialize sample matches in Redis
func initializeSampleMatches() {
	// Check if any match already exists in Redis
	keys, err := redisClient.Keys(ctx, "match:*").Result()
	if err != nil {
		log.Printf("Unable to check if matches exist in Redis: %v", err)
		return
	}

	// If there are already matches, skip initialization
	if len(keys) > 0 {
		log.Println("Matches already exist in Redis, skipping initialization")
		return
	}

	// Initialize matches if Redis is empty
	var matches = []Match{
		{HomeTeam: "Real Madrid", HomeTeamAbbr: "RMA", HomeImg: "/team_logos/Real Madrid.png", AwayTeam: "FC Barcelona", AwayTeamAbbr: "BAR", AwayImg: "/team_logos/FC Barcelona.png", Date: "2023-10-01", Time: "20:00"},
		{HomeTeam: "Atlético de Madrid", HomeTeamAbbr: "ATM", HomeImg: "/team_logos/Atlético de Madrid.png", AwayTeam: "Sevilla FC", AwayTeamAbbr: "SEV", AwayImg: "/team_logos/Sevilla FC.png", Date: "2023-10-02", Time: "20:00"},
		{HomeTeam: "Valencia CF", HomeTeamAbbr: "VAL", HomeImg: "/team_logos/Valencia CF.png", AwayTeam: "Villarreal CF", AwayTeamAbbr: "VIL", AwayImg: "/team_logos/Villarreal CF.png", Date: "2023-10-03", Time: "22:00"},
		{HomeTeam: "Real Sociedad", HomeTeamAbbr: "RSO", HomeImg: "/team_logos/Real Sociedad.png", AwayTeam: "Athletic Bilbao", AwayTeamAbbr: "ATH", AwayImg: "/team_logos/Athletic Bilbao.png", Date: "2023-10-04", Time: "18:00"},
		{HomeTeam: "Real Betis Balompié", HomeTeamAbbr: "BET", HomeImg: "/team_logos/Real Betis Balompié.png", AwayTeam: "Deportivo Alavés", AwayTeamAbbr: "ALA", AwayImg: "/team_logos/Deportivo Alavés.png", Date: "2023-10-05", Time: "20:00"},
		{HomeTeam: "Celta de Vigo", HomeTeamAbbr: "CEL", HomeImg: "/team_logos/Celta de Vigo.png", AwayTeam: "RCD Espanyol Barcelona", AwayTeamAbbr: "ESP", AwayImg: "/team_logos/RCD Espanyol Barcelona.png", Date: "2023-10-06", Time: "22:00"},
	}

	for _, match := range matches {
		// Increment the match ID counter in Redis
		newID, err := redisClient.Incr(ctx, "match_id_counter").Result()
		if err != nil {
			log.Printf("Unable to increment match ID counter in Redis: %v", err)
			continue
		}

		match.ID = int(newID)
		matchID := strconv.Itoa(match.ID)

		// Save match to Redis as hash
		err = redisClient.HSet(ctx, "match:"+matchID, map[string]interface{}{
			"id":           match.ID,
			"homeTeam":     match.HomeTeam,
			"homeTeamAbbr": match.HomeTeamAbbr,
			"homeImg":      match.HomeImg,
			"awayTeam":     match.AwayTeam,
			"awayTeamAbbr": match.AwayTeamAbbr,
			"awayImg":      match.AwayImg,
			"date":         match.Date,
			"time":         match.Time,
		}).Err()

		if err != nil {
			log.Printf("Unable to save match in Redis: %v", err)
		}
	}
}

/***************
	REDIS
****************/

// Function to connect to Redis
func connectToRedis(redisHost string) {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisHost + ":6379", // Redis server address
		Password: "",                  // No password set
		DB:       0,                   // Use default DB
	})

	// Test connection
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Cannot connect to Redis: %v", err)
	}
	log.Println("Connected to Redis successfully")
}

/***************
	REST API - EVENTS
****************/

func createEvent(c *gin.Context) {
	var event Event
	if err := c.BindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Increment the event ID counter in Redis
	newID, err := redisClient.Incr(ctx, "event_id_counter").Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to generate new event ID: " + err.Error()})
		return
	}

	event.ID = int(newID)
	eventID := strconv.Itoa(event.ID) // Convert ID to string

	// Save event to Redis as hash
	err = redisClient.HSet(ctx, "event:"+eventID, map[string]interface{}{
		"id":      event.ID,
		"matchId": event.MatchID,
		"team":    event.Team,
		"player":  event.Player,
		"type":    event.Type,
		"minute":  event.Minute,
	}).Err()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save event in Redis: " + err.Error()})
		return
	}

	// Publish event to RabbitMQ
	publishEvent(event)

	c.JSON(http.StatusCreated, event)
}

func getEvent(c *gin.Context) {
	eventID := c.Param("id")
	// Retrieve event from Redis hash
	result, err := redisClient.HGetAll(ctx, "event:"+eventID).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve event: " + err.Error()})
		return
	}

	// Check if the event exists
	if len(result) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	// Deserialize to Event struct
	var event Event
	event.ID, _ = strconv.Atoi(result["id"])
	event.MatchID, _ = strconv.Atoi(result["matchId"])
	event.Team = result["team"]
	event.Player = result["player"]
	event.Type = result["type"]
	event.Minute, _ = strconv.Atoi(result["minute"])

	c.JSON(http.StatusOK, event)
}

// Handler to get all events of a match
func getAllEvents(c *gin.Context) {
	matchID := c.Param("id")

	// Retrieve all event keys from Redis
	keys, err := redisClient.Keys(ctx, "event:*").Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve event keys: " + err.Error()})
		return
	}

	var allEvents []Event
	// Iterate over each key and retrieve event data
	for _, key := range keys {
		result, err := redisClient.HGetAll(ctx, key).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve event data: " + err.Error()})
			return
		}

		var event Event
		event.ID, _ = strconv.Atoi(result["id"])
		event.MatchID, _ = strconv.Atoi(result["matchId"])
		event.Team = result["team"]
		event.Player = result["player"]
		event.Type = result["type"]
		event.Minute, _ = strconv.Atoi(result["minute"])

		matchIDInt, err := strconv.Atoi(matchID)
		if err != nil {
			log.Printf("Unable to convert matchID to int: %v", err)
			continue
		}
		if event.MatchID == matchIDInt {
			allEvents = append(allEvents, event)
		}
	}

	c.JSON(http.StatusOK, allEvents)
}

// Handler to delete an event by ID
func deleteEvent(c *gin.Context) {
	eventID := c.Param("id")

	// Delete event from Redis
	err := redisClient.Del(ctx, "event:"+eventID).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to delete event: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully"})
}

// Handler to update an existing event by ID
func updateEvent(c *gin.Context) {
	eventID := c.Param("id")
	var event Event
	if err := c.BindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify if the event exists
	exists, err := redisClient.Exists(ctx, "event:"+eventID).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to verify event existence: " + err.Error()})
		return
	} else if exists == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	// Update event in Redis as hash
	err = redisClient.HSet(ctx, "event:"+eventID, map[string]interface{}{
		"id":      event.ID,
		"matchId": event.MatchID,
		"team":    event.Team,
		"player":  event.Player,
		"type":    event.Type,
		"minute":  event.Minute,
	}).Err()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update event in Redis: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, event)
}

/***************
	REST API - MATCHES
****************/

// Handler to create a new match
func createMatch(c *gin.Context) {
	var match Match
	if err := c.BindJSON(&match); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Increment the match ID counter in Redis
	newID, err := redisClient.Incr(ctx, "match_id_counter").Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to generate new match ID: " + err.Error()})
		return
	}

	match.ID = int(newID)
	matchID := strconv.Itoa(match.ID) // Convert ID to string

	// Save match to Redis as hash
	err = redisClient.HSet(ctx, "match:"+matchID, map[string]interface{}{
		"id":           match.ID,
		"homeTeam":     match.HomeTeam,
		"homeTeamAbbr": match.HomeTeamAbbr,
		"homeImg":      match.HomeImg,
		"awayTeam":     match.AwayTeam,
		"awayTeamAbbr": match.AwayTeamAbbr,
		"awayImg":      match.AwayImg,
		"date":         match.Date,
		"time":         match.Time,
	}).Err()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save match in Redis: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, match)
}

// Handler to retrieve a match by ID
func getMatch(c *gin.Context) {
	matchID := c.Param("id")

	// Retrieve match from Redis hash
	result, err := redisClient.HGetAll(ctx, "match:"+matchID).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve match: " + err.Error()})
		return
	}

	// Check if the match exists
	if len(result) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	// Deserialize to Match struct
	var match Match
	match.ID, _ = strconv.Atoi(result["id"])
	match.HomeTeam = result["homeTeam"]
	match.HomeTeamAbbr = result["homeTeamAbbr"]
	match.HomeImg = result["homeImg"]
	match.AwayTeam = result["awayTeam"]
	match.AwayTeamAbbr = result["awayTeamAbbr"]
	match.AwayImg = result["awayImg"]
	match.Date = result["date"]
	match.Time = result["time"]

	c.JSON(http.StatusOK, match)
}

// Handler to get all matches
func getAllMatches(c *gin.Context) {
	// Retrieve all match keys from Redis
	keys, err := redisClient.Keys(ctx, "match:*").Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve match keys: " + err.Error()})
		return
	}

	var allMatches []Match
	// Iterate over each key and retrieve match data
	for _, key := range keys {
		result, err := redisClient.HGetAll(ctx, key).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve match data: " + err.Error()})
			return
		}

		var match Match
		match.ID, _ = strconv.Atoi(result["id"])
		match.HomeTeam = result["homeTeam"]
		match.HomeTeamAbbr = result["homeTeamAbbr"]
		match.HomeImg = result["homeImg"]
		match.AwayTeam = result["awayTeam"]
		match.AwayTeamAbbr = result["awayTeamAbbr"]
		match.AwayImg = result["awayImg"]
		match.Date = result["date"]
		match.Time = result["time"]

		allMatches = append(allMatches, match)
	}

	c.JSON(http.StatusOK, allMatches)
}

// Handler to delete a match by ID
func deleteMatch(c *gin.Context) {
	matchID := c.Param("id")

	// Delete match from Redis
	err := redisClient.Del(ctx, "match:"+matchID).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to delete match: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Match deleted successfully"})
}

// Handler to update an existing match by ID
func updateMatch(c *gin.Context) {
	matchID := c.Param("id")
	var match Match
	if err := c.BindJSON(&match); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify if the match exists
	exists, err := redisClient.Exists(ctx, "match:"+matchID).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to verify match existence: " + err.Error()})
		return
	} else if exists == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	// Update match in Redis as hash
	err = redisClient.HSet(ctx, "match:"+matchID, map[string]interface{}{
		"id":           match.ID,
		"homeTeam":     match.HomeTeam,
		"homeTeamAbbr": match.HomeTeamAbbr,
		"homeImg":      match.HomeImg,
		"awayTeam":     match.AwayTeam,
		"awayTeamAbbr": match.AwayTeamAbbr,
		"awayImg":      match.AwayImg,
		"date":         match.Date,
		"time":         match.Time,
	}).Err()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update match in Redis: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, match)
}

/***************
    RABBITMQ
****************/

// RabbitMQ configuration (unchanged)
var rabbitmq RabbitMQ

// Function to start RabbitMQ
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

// Function to stop RabbitMQ
func stopRabbitMQ() {
	rabbitmq.conn.Close()
	rabbitmq.ch.Close()
	rabbitmq.cancel()
}

// Function to handle errors
func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

// Function to publish events to RabbitMQ (unchanged)
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
