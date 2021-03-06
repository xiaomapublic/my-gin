package mongodb

import (
	"gopkg.in/mgo.v2"
	"my-gin/libraries/mongodb"
	"time"
)

type MyGinData struct {
	Id                 string    `json:"id" bson:"_id"`
	Ad_id              string    `json:"ad_id" bson:"ad_id"`
	Campaign_id        string    `json:"campaign_id" bson:"campaign_id"`
	Product_id         int       `json:"product_id" bson:"product_id"`
	Advertiser_id      int       `json:"advertiser_id" bson:"advertiser_id"`
	Request_count      int       `json:"request_count" bson:"request_count"`
	Cpm_count          int       `json:"cpm_count" bson:"cpm_count"`
	Cpc_original_count int       `json:"cpc_original_count" bson:"cpc_original_count"`
	Division_id        int       `json:"division_id" bson:"division_id"`
	Status             int       `json:"status" bson:"status"`
	Created_at         time.Time `json:"created_at" bson:"created_at"`
	Updated_at         time.Time `json:"updated_at" bson:"updated_at"`
}

type MyGin struct {
}

func (*MyGin) Mongodb() *mgo.Collection {
	return mongodb.MongoSession["mygin"].C("mygin")
}
