package fipmgo

type M map[string]interface{}

/*
import (
	"gopkg.in/mgo.v2/bson"
)

type ObjectId bson.ObjectId



func NewObjectId() ObjectId {
	return ObjectId(bson.NewObjectId())
}

func ObjectIdFromString(id string) ObjectId {
	return ObjectId(bson.ObjectIdHex(id))
}

func (id ObjectId) IsEmpty() bool {
	return string(id) == ""
}

// MarshalJSON turns a bson.ObjectId into a json.Marshaller.
func (id ObjectId) MarshalJSON() ([]byte, error) {
	return bson.ObjectId(id).MarshalJSON()
}

// UnmarshalJSON turns *bson.ObjectId into a json.Unmarshaller.
func (id *ObjectId) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		*id = NewObjectId()
		return nil
	}
	/*
	*id = bson.ObjectId.UnmarshalJSON(data)

	if len(data) != 26 || data[0] != '"' || data[25] != '"' {
		return errors.New(fmt.Sprintf("Invalid ObjectId in JSON: %s", string(data)))
	}
	var buf [12]byte
	_, err := hex.Decode(buf[:], data[1:25])
	if err != nil {
		return errors.New(fmt.Sprintf("Invalid ObjectId in JSON: %s (%s)", string(data), err))
	}
	*id = ObjectId(string(buf[:]))

	return nil
}
*/
