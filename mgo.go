package fipmgo

import (
	"github.com/fipress/fiplog"
	"gopkg.in/mgo.v2"
)

const (
	defaultPageSize = 20
)

var (
	idAggregation = []M{
		{"$group": M{"_id": 1, "ids": M{"$push": "$_id"}}}}
	//`[{$group:{_id:null,ids: {$push: "$_id"}}}]`
)

//todo: connection pool, performance?
var mgoConn *MgoConn

func InitMgo(url, db string) {
	mgoConn = NewMgoConn(url, db)
}

func InitMgoWithConn(conn string) {
	mgoConn = NewMgoConnWithUrl(conn)
}

func InitMgoWithAuth(url, db, user, pwd string) {
	mgoConn = NewMgoConnWithAuth(url, db, user, pwd)
}

func GetMgoConn() (conn *MgoConn) {
	if mgoConn == nil {
		//todo: ensure, try the initial method 3 times
		InitMgo("localhost", "test")
	}
	if mgoConn == nil {
		panic("MongoDB connection failed!")
	}
	return mgoConn
}

func CloseMgoConn() {
	if mgoConn != nil {
		mgoConn.Close()
	}
}

type MgoConn struct {
	*mgo.Database
}

func NewMgoConnWithUrl(connStr string) *MgoConn {
	conn, err := mgo.Dial(connStr)
	if err != nil {
		fiplog.GetLogger().Error("Connect to database failed, conn:", connStr, ",error:", err)
		return nil
	}
	return &MgoConn{conn.DB("")}
}

func NewMgoConn(url, db string) *MgoConn {
	//m := &MgoManager{url,db}
	conn, err := mgo.Dial(url)
	if err != nil {
		fiplog.GetLogger().Error("Connect to database failed, url:", url, ",error:", err)
		return nil
	}
	database := conn.DB(db)
	return &MgoConn{database}
}

func NewMgoConnWithAuth(url, db, user, pwd string) *MgoConn {
	//m := &MgoManager{url,db}
	conn, err := mgo.Dial(url)
	if err != nil {
		fiplog.GetLogger().Error("Connect to database failed, url:", url, ",error:", err)
	}
	database := conn.DB(db)
	err = database.Login(user, pwd)
	if err != nil {
		fiplog.GetLogger().Error("Database auth error:", err)
	}
	return &MgoConn{database}
}

func (m *MgoConn) Close() {
	m.Database.Session.Close()
}

type MgoEntity interface {
	//	Coll() string
	GetId() interface{}
	//    SetId(id interface {})
}

//type EntityConstructor func() MgoEntity

/*
type Accessor interface {
	Get(id interface {}) interface {}
	Insert(entities ...interface {}) bool
}
*/

type MgoAccessor struct {
	CollName string
	//constructor interface {}
}

/*func (m *MgoAccessor) CreateEntity() {
	panic("Abstract method, to be implemented.")
}*/

func (m *MgoAccessor) coll() *mgo.Collection {
	return GetMgoConn().C(m.CollName)
}

// Get entity by id. Should pass the pointer of an empty entity in.
// Usage:
//    	e := &Entity{}
//		ok := testAccessor.Get(id, e)
func (m *MgoAccessor) Get(id interface{}, entity interface{}) bool {
	err := m.coll().FindId(id).One(entity)
	if err != nil {
		fiplog.GetLogger().Error("Get entity failed. id:", id, ",coll:", m.CollName, "error:", err)
		return false
	}
	return true
}

func (m *MgoAccessor) GetById(id interface{}, selector interface{}, entity interface{}) bool {
	err := m.coll().FindId(id).Select(selector).One(entity)
	if err != nil {
		fiplog.GetLogger().Error("Get partial entity failed. id:", id, ",coll:", m.CollName, "error:", err)
		return false
	}
	return true
}

func (m *MgoAccessor) GetPartial(id interface{}, entity interface{}, fields ...string) bool {
	if len(fields) == 0 {
		fiplog.GetLogger().Error("Get partial entity, no field specified")
		return false
	}

	selector := M{}
	for _, field := range fields {
		selector[field] = 1
	}

	return m.GetById(id, selector, entity)
}

func (m *MgoAccessor) FindOne(query interface{}, selector interface{}, result interface{}, sortFields ...string) bool {
	q := m.coll().Find(query).Select(selector)

	if len(sortFields) != 0 {
		q = q.Sort(sortFields...)
	}

	err := q.One(result)

	if err != nil {
		fiplog.GetLogger().Error("Find one failed. query:", query, ",coll:", m.CollName, "error:", err)
		return false
	}
	return true
}

// Get a list of entities by the selector. Should pass the pointer of an empty list in
// Usage:
//		var list []Entity
//		ok := testAccessor.FindAll(id, &list)
func (m *MgoAccessor) FindAll(query interface{}, selector interface{}, result interface{}) {
	err := m.coll().Find(query).Select(selector).All(result)
	fiplog.GetLogger().Debug("query:", query, ",result", result)

	if err != nil {
		fiplog.GetLogger().Error("Find all failed. query:", query, ",coll:", m.CollName, "error:", err)
	}
}

//pageIndex starts with 1
func (m *MgoAccessor) FindPage(query, selector, result interface{}, pageSize, pageIndex int, sortFields ...string) (pageCount int, ok bool) {
	ok = true
	if pageSize == 0 {
		pageSize = defaultPageSize
	}

	if pageIndex < 1 {
		pageIndex = 1
	}

	q := m.coll().Find(query).Select(selector)
	total, err := q.Count()
	if err != nil {
		fiplog.GetLogger().Error("Find page failed. query:", query, ",coll:", m.CollName, "error:", err)
		ok = false
		return
	}
	skip := pageSize * (pageIndex - 1)
	if total <= skip {
		return
	}

	fiplog.GetLogger().Debug("find page. total:", total, ",query:", query)
	pageCount = total / pageSize

	/*if pageCount == 0 || pageCount < pageIndex {
		fiplog.GetLogger().Error("Find page failed, pageCound:",pageCount,", pageIndex:",pageIndex)
		return
	}*/
	if len(sortFields) != 0 {
		q = q.Sort(sortFields...)
	}

	err = q.Skip(skip).Limit(pageSize).All(result)
	if err != nil {
		fiplog.GetLogger().Error("Find all failed. query:", query, ",coll:", m.CollName, "error:", err)
		ok = false
		return
	}
	fiplog.GetLogger().Debug("find page. result:", result, ",skip:", skip, ",limit:", pageSize)
	return
}

func (m *MgoAccessor) Find(query interface{}) *mgo.Query {
	return m.coll().Find(query)
}

func (m *MgoAccessor) Insert(entity ...interface{}) bool {
	err := m.coll().Insert(entity...)
	if err != nil {
		fiplog.GetLogger().Error("Insert entity failed. error:", ",coll:", m.CollName, err)
		return false
	}
	//entity.SetId(id)
	return true
}

/*
type ChangeInfo struct {
	*mgo.ChangeInfo
}*/
func (m *MgoAccessor) Upsert(entity MgoEntity) bool {
	_, err := m.coll().UpsertId(entity.GetId(), entity)
	if err != nil {
		fiplog.GetLogger().Error("Upsert entity failed. id:", entity.GetId(), ",coll:", m.CollName, "error:", err)
		return false
	}
	return true
}

func (m *MgoAccessor) UpsertPartial(id interface{}, update interface{}) bool {
	_, err := m.coll().UpsertId(id, update)
	if err != nil {
		fiplog.GetLogger().Error("Upsert entity failed. id:", id, "error:", err)
		return false
	}
	return true
}

func (m *MgoAccessor) Update(entity MgoEntity) bool {
	err := m.coll().UpdateId(entity.GetId(), entity)
	if err != nil {
		fiplog.GetLogger().Error("Update entity failed. id:", entity.GetId(), ",coll:", m.CollName, "error:", err)
		return false
	}
	return true
}

func (m *MgoAccessor) UpdatePartialById(id, update interface{}) bool {
	fiplog.GetLogger().Debug("id:", id, "update:", update)
	err := m.coll().UpdateId(id, update)
	if err != nil {
		fiplog.GetLogger().Error("Update partial by id failed. id:", id, ",coll:", m.CollName, "error:", err)
		return false
	}
	return true
}

func (m *MgoAccessor) UpdatePartial(selector interface{}, update interface{}) bool {
	//fiplog.GetLogger().Debug("update:",selector,update)
	err := m.coll().Update(selector, update)
	if err != nil {
		fiplog.GetLogger().Error("Partial update failed.", selector, update, ",coll:", m.CollName, "error:", err)
		return false
	}

	return true
}

func (m *MgoAccessor) InsertToArrayField(id interface{}, field string, val interface{}) bool {
	/*bytes, err := json.Marshal(val)
	if err != nil {
		fiplog.GetLogger().Error("insert to array field failed.field:", field, "error:", err)
		return false
	}*/
	//fiplog.GetLogger().Debug("json:",string(bytes))
	err := m.coll().UpdateId(id, M{"$push": M{field: val}})
	if err != nil {
		fiplog.GetLogger().Error("insert to array field. field:", field, ",coll:", m.CollName, "error:", err)
		return false
	}
	return true
}

/*
func (m *MgoAccessor) DeleteFromArrayField(id interface{}, field string, val interface{} )

func (m *MgoAccessor) UpdateInArrayField(id interface{}, field string, selector interface{}, update interface{} )  {
	err := m.coll().Update(id, M{"$push",M{field,string(bytes)}})
	if err != nil {
		fiplog.GetLogger().Error("insert to array field. field:", field, "error:", err)
		return false
	}
	return true
}
*/
func (m *MgoAccessor) Delete(id interface{}) bool {
	err := m.coll().RemoveId(id)
	if err != nil {
		fiplog.GetLogger().Error("Delete failed. Id:", id, ",coll:", m.CollName, ",error:", err)
		return false
	}
	return true
}

func (m *MgoAccessor) DeleteAll(query interface{}) int {
	info, err := m.coll().RemoveAll(query)
	if err != nil {
		fiplog.GetLogger().Error("Delete all failed. query:", query, ",coll:", m.CollName, ",error:", err)
		return 0
	}
	return info.Removed
}

func (m *MgoAccessor) Exists(query interface{}) (exists bool, err error) {
	count, err := m.coll().Find(query).Count()
	if err != nil {
		fiplog.GetLogger().Error("Check if exists failed. query:", query, ",coll:", m.CollName, ",error:", err)
		return
	}
	exists = count != 0
	return
}

func (m *MgoAccessor) IdExists(id interface{}) (exists bool, err error) {
	count, err := m.coll().FindId(id).Count()
	if err != nil {
		fiplog.GetLogger().Error("Check if id exists failed. id:", id, ",coll:", m.CollName, ",error:", err)
		return
	}
	exists = count != 0
	return
}

func (m *MgoAccessor) PipeOne(pipeline, result interface{}) bool {
	pipe := m.coll().Pipe(pipeline)
	err := pipe.One(result)
	if err != nil {
		fiplog.GetLogger().Error("pipe one failed. pipeline:", pipeline, ",coll:", m.CollName, ",error:", err)
		return false
	} else {
		return true
	}
}

func (m *MgoAccessor) PipeAll(pipeline, result interface{}) bool {
	pipe := m.coll().Pipe(pipeline)
	err := pipe.All(result)
	if err != nil {
		fiplog.GetLogger().Error("pipe all failed. pipeline:", pipeline, ",coll:", m.CollName, ",error:", err)
		return false
	} else {
		return true
	}
}

type a struct {
	Id  int
	Ids interface{}
}

func (m *MgoAccessor) GetAllIds() interface{} {
	var result []a
	//result.Ids = list
	err := m.coll().Pipe(idAggregation).All(&result)
	if err != nil {
		fiplog.GetLogger().Error("Get all id failed. coll ", m.CollName, ",error:", err)
		return nil
	}
	return result[0].Ids
}
