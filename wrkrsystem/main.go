package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"wrksystem/models"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}
	models.Setup(os.Getenv("DBHOST"), os.Getenv("DBUSERNAME"), os.Getenv("DBPASSWORD"), os.Getenv("DATABASE"), os.Getenv("DBPORT"))

	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	msgs, _ := ch.Consume("transactions", "", true, false, false, false, nil)
	for msg := range msgs {
		var tx models.Transaction
		json.Unmarshal(msg.Body, &tx)
		processTransaction(tx)
	}
}

func processTransaction(tx models.Transaction) {
	err := models.db.Transaction(func(txDB *gorm.DB) error {
		var sender, receiver models.Client
		if err := txDB.First(&sender, tx.SenderID).Error; err != nil {
			return err
		}
		if err := txDB.First(&receiver, tx.ReceiverID).Error; err != nil {
			return err
		}

		if sender.Balance < tx.Amount {
			tx.Status = "failed"
		} else {
			sender.Balance -= tx.Amount
			receiver.Balance += tx.Amount
			tx.Status = "completed"

			if err := txDB.Save(&sender).Error; err != nil {
				return err
			}
			if err := txDB.Save(&receiver).Error; err != nil {
				return err
			}
		}

		tx.CreatedAt = time.Now()
		return txDB.Create(&tx).Error
	})

	if err != nil {
		log.Printf("Transaction failed: %v", err)
	}
}
