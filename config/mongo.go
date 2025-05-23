package config

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

func InitMongo() *mongo.Client {
	// Cargar variables del .env
	if err := godotenv.Load(); err != nil {
		log.Println("No se pudo cargar el archivo .env")
	}

	mongoURI := os.Getenv("MONGO_URI")

	if mongoURI == "" {
		log.Fatal("Defina los ENV")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Error al conectar con MongoDB: %v", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("❌ Ping a MongoDB falló: %v", err)
	}

	log.Println("Conexión a MongoDB establecida")

	MongoClient = client
	return client
}
func GetDB() *mongo.Database {
	dbName := os.Getenv("MONGO_DB_NAME")
	if dbName == "" {
		log.Fatal("❌ MONGO_DB_NAME no está definido en el entorno")
	}
	return MongoClient.Database(dbName)
}
