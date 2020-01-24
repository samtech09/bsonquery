package bsonquery

import (
	"context"
	"testing"

	"github.com/samtech09/dbtools/mango"
	"go.mongodb.org/mongo-driver/mongo"
)

type testuser struct {
	ID    int    `bson:"_id"`
	Name  string `bson:"name"`
	Age   int    `bson:"age"`
	Score int    `bson:"score"`
}

func (t testuser) GetID() interface{} {
	return t.ID
}
func (t *testuser) ToInterface(list []testuser) []interface{} {
	var islice []interface{} = make([]interface{}, len(list))
	for i, d := range list {
		islice[i] = d
	}
	return islice
}

func TestMongo(t *testing.T) {
	cfg := mango.MongoConfig{}
	cfg.Host = "192.168.60.206"
	cfg.Port = 27017
	cfg.DbName = "testdb"

	ses := mango.InitSession(cfg)
	defer ses.Cleanup()
	c := ses.GetColl("testusers")

	c.Drop(context.Background())

	// insert documents
	usr := []testuser{}
	usr = append(usr, testuser{1, "test1", 44, 9})
	usr = append(usr, testuser{2, "test2", 32, 7})
	usr = append(usr, testuser{3, "test3", 29, 9})
	usr = append(usr, testuser{4, "test3", 42, 5})
	usr = append(usr, testuser{5, "test3", 36, 6})
	u := testuser{}
	islice := u.ToInterface(usr)
	for i, d := range usr {
		islice[i] = d
	}
	err := ses.InsertBulk(c, islice...)
	if err != nil {
		t.Errorf("Error inserting documents: %s", err.Error())
		t.FailNow()
	}

	// test filter
	//   should return single record i.e. testuser{4, "test2", 37, 7}
	filter := Builder().
		And(C().EQ("name", "test2"), C().GT("age", 29)).
		Build()
	exp := 1
	cur, err := c.Find(context.Background(), filter)
	if err != nil {
		t.Errorf("Error finding docs with filter: %s", err.Error())
		t.FailNow()
	}
	cnt := countDocs(cur)
	if cnt != exp {
		t.Errorf("Filter-1 failed. Expected %d,  Got: %d", exp, cnt)
	}

	// and with or query (score > 5) and (age > 30 or age < 45)
	//   satisfied records are
	//		testuser{1, "test1", 44, 9}
	//		testuser{2, "test2", 32, 7}
	/*
		{ "$or": [ {"age": {"$gt": 30}}, {"score": {"$gt": 5}} ] }
	*/
	exp = 2
	filter = Builder().
		Or(C().GT("age", 40), C().GT("score", 7)).
		Build()
	cur, err = c.Find(context.Background(), filter)
	if err != nil {
		t.Errorf("Error finding docs with filter: %s", err.Error())
		t.FailNow()
	}
	cnt = countDocs(cur)
	if cnt != exp {
		t.Errorf("Filter-2 failed. Expected %d,  Got: %d", exp, cnt)
	}
}

func countDocs(cur *mongo.Cursor) int {
	count := 0
	for cur.Next(context.TODO()) {
		count++
	}
	if err := cur.Err(); err != nil {
		panic("Error interating curser: " + err.Error())
	}
	return count
}
