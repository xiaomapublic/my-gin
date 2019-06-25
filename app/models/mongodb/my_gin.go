package mongodb

import (
	"gopkg.in/mgo.v2"
	"my-gin/app/libraries/mongodb"
	"time"
)

type MyGinData struct {
	Id               string         		`bson:"_id"`
	Name             string                 `bson:"name"`
	Location_id      int                    `bson:"location_id"`
	Product_id       int                    `bson:"product_id"`
	Advertiser_id    int                    `bson:"advertiser_id"`
	Creative_type_id int                    `bson:"creative_type_id"`
	Created_at       time.Time             	`bson:"created_at"`
}

type MyGin struct {

}

func (*MyGin) Mongodb() *mgo.Collection {
	return mongodb.MongoSession["default"].C("my_gin")
}
