package audit

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	err = m.applyAuditChecks(audit)
	if err != nil {
		return err
	}

	//asign the db, db client and the collection
	m.db = db
	m.dbClient = client

	m.audit = audit

	return nil
}

func (m *database) applyAuditChecks(audit *collectionWrapper) error {
	log.Println("apply audit checks.....")

	//add indexes
	err := audit.AddIndex(bson.D{primitive.E{Key: "user_identifier", Value: 1}}, false)
	if err != nil {
		return err
	}
	err = audit.AddIndex(bson.D{primitive.E{Key: "used_group", Value: 1}}, false)
	if err != nil {
		return err
	}
	err = audit.AddIndex(bson.D{primitive.E{Key: "entity", Value: 1}}, false)
	if err != nil {
		return err
	}
	err = audit.AddIndex(bson.D{primitive.E{Key: "entity_id", Value: 1}}, false)
	if err != nil {
		return err
	}
	err = audit.AddIndex(bson.D{primitive.E{Key: "operation", Value: 1}}, false)
	if err != nil {
		return err
	}
	err = audit.AddIndex(bson.D{primitive.E{Key: "created_at", Value: 1}}, false)
	if err != nil {
		return err
	}

	log.Println("audit checks passed")
	return nil
}
