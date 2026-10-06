package main

import (
	"context"
	"log"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-pkg/database/seeder"
	"github.com/MamangRust/monolith-point-of-sale-pkg/dotenv"
	"github.com/MamangRust/monolith-point-of-sale-pkg/hash"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
)

func main() {
	if err := dotenv.Viper(); err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	l, err := logger.NewLogger("seeder", nil)
	if err != nil {
		log.Fatalf("Error creating logger: %v", err)
	}

	mongoClient, err := salemongo.NewMongoClient(l)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}
	mongoDB := salemongo.NewDatabase(mongoClient)

	ctx := context.Background()

	s := seeder.NewSeeder(seeder.Deps{
		Db:     mongoDB,
		Ctx:    ctx,
		Logger: l,
		Hash:   hash.NewHashingPassword(),
	})

	if err := s.Run(); err != nil {
		log.Fatalf("Seeding failed: %v", err)
	}

	l.Info("Seeding completed successfully.")
}
