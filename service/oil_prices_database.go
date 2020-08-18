package service

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	// the time format in database
	timeLayout = "20060102"
)

// OilPriceDB - class for oil price database
type OilPriceDB struct {
	client   *mongo.Client
	database *mongo.Database
}

// OilPriceData - format of oil price data
type OilPriceData struct {
	ID			primitive.ObjectID		`json:"id" bson:"_id, omitempty"`
	Name 		string 					`json:"name" bson:"name"`
	Prices		[]float64				`json:"oilPrices" bson:"oilPrices"`
}

// Init initializes OilPriceDB
func (db *OilPriceDB) Init(dburl string) error {
	// 连接到数据库
	client, err := mongo.NewClient(options.Client().ApplyURI(dburl))
	if err != nil {
		return err
	}
	db.client = client
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = db.client.Connect(ctx)
	if err != nil {
		return err
	}

	// check connection
	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	}
	fmt.Println("Successfully connected to server")

	// check for updates
	db.database = db.client.Database("oilPricesDB")
	updateFlag, err := db.needsUpdate(ctx, db.database.Collection("time"))
	if err != nil {
		return err
	}
	fmt.Println("Needs update? ", updateFlag)
	if updateFlag {
		return db.update()
	}
	return nil
}

// Close - disconnects from database
func (db *OilPriceDB) Close() {
	db.client.Disconnect(context.TODO())
}

// Read - reads oil price data
func (db *OilPriceDB) Read(province string, oilType string) ([]float64, error) {
	collection := db.database.Collection("prices")
	result := OilPriceData{}
	filter := bson.D{{"name", province}}
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return make([]float64, 0), err
	}
	idx, ok := oilIndex[oilType]
	if ok {
		return result.Prices[idx : idx+1], nil
	} else {
		return result.Prices, nil
	}
}

// write - writes oil price data into database
func (db *OilPriceDB) write(database *mongo.Database, prices map[string][]float64) error {
	collection := database.Collection("prices")
	for province, oilPrices := range prices {
		result := OilPriceData{}
		filter := bson.D{{"name", province}}
		err := collection.FindOne(context.TODO(), filter).Decode(&result)
		if err != nil {

			// if the province is new, add a new record
			if err == mongo.ErrNoDocuments {
				_, err := collection.InsertOne(context.TODO(), &OilPriceData{
					ID:		primitive.NewObjectID(),
					Name:	province,
					Prices:	oilPrices,
				})
				if err != nil {
					return err
				}
			} else {
				return err
			}
		} else {
			update := bson.D{
				{"$set", bson.D{
					{"oilPrices", bson.A{oilPrices}},
				}},
			}
			_, err := collection.UpdateOne(context.TODO(), filter, update)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// updates oil price data
func (db *OilPriceDB) update() error {
	data, err := crawlOilPrices()
	if err != nil {
		return err
	}
	return db.write(db.database, data)
}

// checks if the data were last updated within past 24 hours
func (db *OilPriceDB) needsUpdate(ctx context.Context, collection *mongo.Collection) (bool, error) {
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	result := struct {
		Time string
	}{}
	filter := bson.D{{"name", "time"}}
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		
		// if time is not found, add it
		if err == mongo.ErrNoDocuments {
			_, err := collection.InsertOne(context.TODO(), bson.D{{"name", "time"}, {"time", time.Now().Format(timeLayout)}})
			if err != nil {
				return false, err
			}
			return true, nil
		} else {
			return false, err
		}
	}

	last, err := time.Parse(timeLayout, result.Time)
	if err != nil {
		return false, err
	}
	now := time.Now()
	diff := now.Sub(last).Hours()
	if diff >= 24 {
		update := bson.D{
			{"$set", bson.D{
				{"time", now.Format(timeLayout)},
			}},
		}
		_, err := collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}
