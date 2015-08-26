package fipmgo

import (
	"testing"
	"gopkg.in/mgo.v2/bson"
	"fiplog"
	"fmt"
//	"gopkg.in/mgo.v2/mgo"
)

type Test struct {
	Id bson.ObjectId `bson:"_id,omitempty"`
	Name string
	Age int
}

func (test *Test) GetId() interface {} {
	return test.Id
}

func (test *Test) SetId(id interface {} ) {
	fiplog.GetLogger().Info("id:",id)
	fmt.Println("id:",id)
	switch rid := id.(type) {
	case bson.ObjectId:
		test.Id = rid
	default:

		panic("type not match")
	}
}

func CreateTest() *Test {
	return &Test{}
}


/*func (test *Test) Create() interface {} {
	return &
}*/


var testAccessor = &MgoAccessor{"test"}



func TestInsert(t *testing.T) {
	InitMgoWithAuth("localhost","test","dba","Dcn103@Mongo")
	tt := &Test{}
	tt.Id = bson.NewObjectId()
	tt.Name = "Abby"
	tt.Age = 33

	//testAccessor.Insert(tt)

	t2 := &Test{}
	t2.Id = bson.NewObjectId()
	t2.Name = "Tony"
	t2.Age = 35

	testAccessor.Insert(tt,t2)
	/*err := GetMgoConn().C("test").Insert(tt)
	if err != nil {
		t.Error(err)
	}*/

	t.Log("return id", tt.Id.Hex())
	//var v = &Test{}
	/*_, v := testAccessor.Get(tt.Id).(*Test)

	t.Log(v.Id)
	t.Log(v.Name)
	t.Log(v.Age)*/
}

func TestUpdate(t *testing.T) {
	InitMgoWithAuth("localhost","test","dba","Dcn103@Mongo")
	tt := &Test{}
	tt.Id = bson.ObjectIdHex("552681af421aa93b68000001")
	//tt.Name = "Abby"
	tt.Age = 34
	testAccessor.Update(tt)
}

func TestUpdatePartial(t *testing.T) {
	InitMgoWithAuth("localhost","test","dba","Dcn103@Mongo")
	//tt := &Test{}
	//tt.Id = bson.ObjectIdHex("552681af421aa93b68000001")
	//tt.Name = "Abby"
	//tt.Age = 34
	testAccessor.UpdatePartialById(bson.ObjectIdHex("552681af421aa93b68000001"),bson.M{"$inc":bson.M{"age":7}})
}

func TestDelete(t *testing.T) {
	InitMgoWithAuth("localhost","test","dba","Dcn103@Mongo")
	testAccessor.Delete(bson.ObjectIdHex("552681af421aa93b68000002"))
}

func TestGet(t *testing.T) {
	InitMgoWithAuth("localhost","test","dba","Dcn103@Mongo")
	t1 := &Test{}
	ok := testAccessor.Get(bson.ObjectIdHex("552681af421aa93b68000001"),t1)
	if ok {
		t.Log("name:", t1.Name)
		t.Log("age:", t1.Age)
	} else {
		t.Log("not found")
	}
}

func TestFindAll(t *testing.T) {
	InitMgoWithAuth("localhost","test","dba","Dcn103@Mongo")
	var ret []*Test
	testAccessor.FindAll(bson.M{"age":3},&ret)
	for _,r := range ret {
			t.Log("id:",r.Id)
			t.Log("name:",r.Name)
			t.Log("age:",r.Age)
	}

}

func TestSingleUrl(t *testing.T) {
	t.Log("start")
	InitMgoWithConn("mongodb://dev:Dcn103Mongo@localhost/test")
	t.Log("inited")
	t2 := &Test{}
	t2.Id = bson.NewObjectId()
	t2.Name = "Jimmy"
	t2.Age = 37

	testAccessor.Insert(t2)

	t3 := new(Test)
	testAccessor.Get(t2.Id,t3)
	t.Log("t3.Name",t3.Name)
}
