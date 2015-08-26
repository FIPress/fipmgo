package fipmgo

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)


type ObjectId bson.ObjectId
type M map[string]interface{}

func NewObjectId() ObjectId {
	return ObjectId(bson.NewObjectId())
}

func ObjectIdFromString(id string) ObjectId {
	return ObjectId(bson.ObjectIdHex(id))
}

func (o ObjectId) String() string {
	return bson.ObjectId(o).Hex()
}

func (o ObjectId) Time() time.Time {
	return bson.ObjectId(o).Time()
}
