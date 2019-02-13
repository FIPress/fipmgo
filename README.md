# fipmgo

[fipmgo](https://fipress.org/project/fipmgo) is a wrapper of [mgo](). It provides some convenient api to `mgo`.

**Usage**

1. Init connection
```
dbConn := "mongodb://localhost/test"
fipmgo.InitMgoWithConn(dbConn)
```

2. Get accessor by collection name
```
var userAccessor = &MgoAccessor{"user"}
```

3. Use the accessor to do CRUD on the connection
```
userAccessor.Insert(user)
userAccessor.Update(user)
userAccessor.UpdatePartial(user)
userAccessor.Delete(user)
userAccessor.DeleteAll(query)

//query
userAccessor.Get(id,result)
userAccessor.FindOne(query,selector,result,sort-field...)
userAccessor.FindAll(query,selector,result)
userAccessor.FindPage(query,selector,result,pageSize,pageIndex,sort-field...)

//pipeline
userAccessor.PipeOne(pipeline, result)
userAccessor.PipeAll(pipeline, result)
```

For detailed usage, please visit the [ptoject page](https://fipress.org/project/fipmgo).