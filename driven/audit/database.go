package audit

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type database struct {
	mongoDBAuth  string
	mongoDBName  string
	mongoTimeout time.Duration

	db       *mongo.Database
	dbClient *mongo.Client

	audit *collectionWrapper
}

func (m *database) start() error {
	log.Println("audit database -> start")

	//connect to the database
	clientOptions := options.Client().ApplyURI(m.mongoDBAuth)
	connectContext, cancel := context.WithTimeout(context.Background(), m.mongoTimeout)
	client, err := mongo.Connect(connectContext, clientOptions)
	cancel()
	if err != nil {
		return err
	}

	//ping the database
	pingContext, cancel := context.WithTimeout(context.Background(), m.mongoTimeout)
	err = client.Ping(pingContext, nil)
	cancel()
	if err != nil {
		return err
	}

	//apply checks
	db := client.Database(m.mongoDBName)
	audit := &collectionWrapper{database: m, coll: db.Collection("audit")}
	err = m.applyAuditChecks(configs)
	if err != nil {
		return err
	}

	//asign the db, db client and the collection
	m.db = db
	m.dbClient = client

	m.audit = audit

	return nil
}

func (m *database) applyAuditChecks(configs *collectionWrapper) error {
	log.Println("apply audit checks.....")

	//TODO

	log.Println("audit checks passed")
	return nil
}
