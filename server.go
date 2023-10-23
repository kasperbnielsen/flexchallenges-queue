package main

import (
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func connectToRabbit(conn *amqp.Connection) {

	ch, err := conn.Channel()
	failOnError(err, "Failed to open channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     //arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	failOnError(err, "Couldnt register consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			fmt.Printf("received a message: %s", d.Body)
			var matchArray []string = []string{}
			err := json.Unmarshal(d.Body, &matchArray)
			failOnError(err, "Couldnt unmarshal matches")
			for _, item := range matchArray {
				fmt.Println(item)
				getMatch(item)
			}

			t := time.Duration(13)
			time.Sleep(t * time.Second)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

type Result struct {
	Puuid        string
	RevisionDate string
}

type User struct {
	ID           primitive.ObjectID `bson:"_id"`
	Puuid        string             `bson:"puuid"`
	RevisionDate float64            `bson:"revisionDate"`
	//accountID     string             `bson:"accountId"`
	//name          string             `bson:"name"`
	//profileIconID float64            `bson:"profileIconId"`
	//region        string             `bson:"region"`
	//summonerID    string             `bson:"summonerId"`
	//username      string             `bson:"username"`
	//iD2           string             `bson:"id"`
	//summonerLevel float64            `bson:"summonerLevel"`
}

func getAllActiveUsers() []Result {

	URI := "mongodb+srv://gokasper:a4e43ce97@production.edtakaz.mongodb.net/?retryWrites=true&w=majority"

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(URI))

	failOnError(err, "Couldnt connect to mongodb")

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database("production").Collection("puuid")

	cursor, err := coll.Find(context.TODO(), bson.D{})

	failOnError(err, "Couldnt find entries")

	var result []User = []User{}

	err = cursor.All(context.TODO(), &result)
	failOnError(err, "Cant find documents")
	results := []Result{}

	for _, user := range result {
		output, err := json.Marshal(user.RevisionDate)
		failOnError(err, "error here")
		output2, err := json.Marshal(user.Puuid)
		failOnError(err, "error puuid not found")
		temp := Result{string(output2[1:]), string(output)}
		results = append(results, temp)
	}

	return results

}

var matches map[string]bool = make(map[string]bool)

func queueAll(conn *amqp.Connection, users []Result) {

	for i := 0; i < len(users); i++ {
		if i%95 == 0 && i != 0 {
			fmt.Println("Ratelimit! Waiting 120 seconds")
			time.Sleep(120 * time.Second)
		} else if i%15 == 0 {
			fmt.Println("Waiting 1 second")
			time.Sleep(1 * time.Second)
		}
		matchIds := getMatchIds(users[i].Puuid[:len(users[i].Puuid)-1], "1")
		var tempMatches []string = []string{}
		err := json.Unmarshal(matchIds, &tempMatches)
		for j := 0; j < len(tempMatches); j++ {
			failOnError(err, "Couldnt convert matches from bytes")
			matches[tempMatches[j]] = true
		}
	}

	var temp []string = []string{}
	count := 0
	for k := range matches {
		if count < 10 {
			temp = append(temp, k)
			count++
		} else {
			produceToQueue(conn, temp)
			temp = temp[:0]
			count = 0
		}
	}
}

func produceToQueue(conn *amqp.Connection, job []string) {
	ch, err := conn.Channel()
	failOnError(err, "Failed to open channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     //arguments
	)
	failOnError(err, "Failed to declare a queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	data, err := json.Marshal(job)
	failOnError(err, "couldnt convert to json")

	err = ch.PublishWithContext(ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(data),
		})
	failOnError(err, "Couldnt publish matches")
	fmt.Println(data)
}

func checkUser() {

}

func main() {

	//ch := make(chan bool)

	conn, err := amqp.Dial("amqps://ryiuuexx:bJGW34j_OGN4hVmDeae-1lg3Vl9oBJaA@stingray.rmq.cloudamqp.com/ryiuuexx")
	failOnError(err, "Couldnt connect")

	//matchids := getMatchIds(users[1].Puuid[:len(users[1].Puuid)-1], "1")

	//var tempvar []string

	//err = json.Unmarshal(matchids, &tempvar)

	//failOnError(err, "couldnt unmarshal")

	//matches["EUW231123"] = true
	//matches["EUW439043"] = true
	//users := getAllActiveUsers()

	//fmt.Println(users)

	//go queueAll(conn, users)

	//fmt.Println("Producing to queue ...")

	//<-ch

	//fmt.Println("finished")

	connectToRabbit(conn)

	/*for i := 1; i < len(users); i++ {
		res := getMatchIds(users[i].Puuid[:len(users[i].Puuid)-1], "1")
	}*/

}
