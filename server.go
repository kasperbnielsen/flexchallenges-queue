package main

import (
	"context"
	//amqp "github.com/rabbitmq/amqp091-go"
	"log"
	//"time"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

/*func connectToRabbit() {
	conn, err := amqp.Dial("amqps://ryiuuexx:bJGW34j_OGN4hVmDeae-1lg3Vl9oBJaA@stingray.rmq.cloudamqp.com/ryiuuexx")
	failOnError(err, "Couldnt connect")



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


	}*/

	type Result struct {
		Puuid string
		RevisionDate string
	}

	type User struct {
		ID primitive.ObjectID `bson:"_id"`
		Puuid string `bson:"puuid`
		accountId string
		name string
		profileIconId float64
		region string
		RevisionDate float64 `bson:"revisionDate`
		summonerId string
		summonerLevel float64
		username string
	}

func main() {

	

	URI := "mongodb+srv://gokasper:a4e43ce97@production.edtakaz.mongodb.net/?retryWrites=true&w=majority"

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(URI))

	failOnError(err, "Couldnt connect to mongodb");

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database("production").Collection("puuid");

	cursor, err := coll.Find(context.TODO(), bson.D{{"region", "EUW1"}})

	failOnError(err, "Couldnt find entries");

	var result []User

	err = cursor.All(context.TODO(), &result)
	failOnError(err, "Cant find documents")

	var results []Result

	results = make([]Result, 1)

	

	for _, user := range result {
		cursor.Decode(&user)
		output, err := json.Marshal(user.RevisionDate)
		failOnError(err, "error here")
		output2, err := json.Marshal(user.Puuid)
		failOnError(err, "error puuid not found")
		temp := Result{string(output2[1:]), string(output)}
		results = append(results, temp)
	}


	res := getMatchIds(results[1].Puuid[:len(results[1].Puuid)-1], fmt.Sprint(int(result[1].RevisionDate)))

	fmt.Println(res)

}
