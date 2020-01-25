package bsonquery_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	bq "github.com/samtech09/bsonquery"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testuser struct {
	ID    int    `bson:"_id"`
	Name  string `bson:"name"`
	Age   int    `bson:"age"`
	Score int    `bson:"score"`
}
type testuser2 struct {
	ID         int    `bson:"_id"`
	Name       string `bson:"name"`
	Age        int    `bson:"age"`
	Score      int    `bson:"score"`
	ExtraField string `bson:"extra"`
}

var client *mongo.Client
var c *mongo.Collection

func (t *testuser) ToInterface(list []testuser) []interface{} {
	var islice []interface{} = make([]interface{}, len(list))
	for i, d := range list {
		islice[i] = d
	}
	return islice
}
func (t *testuser2) ToInterface(list []testuser2) []interface{} {
	var islice []interface{} = make([]interface{}, len(list))
	for i, d := range list {
		islice[i] = d
	}
	return islice
}

func TestInit(t *testing.T) {
	var err error
	dbname := "testdb"
	uri := fmt.Sprintf("mongodb://%s:%d/%s", "192.168.60.206", 27017, dbname)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		t.Errorf("Error connecting mongodb: %s", err.Error())
		t.FailNow()
	}

	c = client.Database(dbname).Collection("testusers")
	c.Drop(context.Background()) // drop exisitng documents if any

	// insert documents
	usr := []testuser{}
	usr = append(usr, testuser{1, "test1", 44, 9})
	usr = append(usr, testuser{2, "test2", 32, 7})
	usr = append(usr, testuser{3, "test3Reg", 29, 9})
	usr = append(usr, testuser{4, "test3Big", 42, 5})
	usr = append(usr, testuser{5, "test3Bit", 36, 6})
	u := testuser{}
	islice := u.ToInterface(usr)
	for i, d := range usr {
		islice[i] = d
	}
	_, err = c.InsertMany(context.Background(), islice)
	if err != nil {
		t.Errorf("Error inserting documents: %s", err.Error())
		t.FailNow()
	}

	// add testuser2 with extra field
	usr2 := []testuser2{}
	usr2 = append(usr2, testuser2{6, "test3Extra", 39, 6, "abc"})
	usr2 = append(usr2, testuser2{7, "test2Ex1", 32, 7, "123"})
	u2 := testuser2{}
	islice2 := u2.ToInterface(usr2)
	for i, d := range usr2 {
		islice2[i] = d
	}
	_, err = c.InsertMany(context.Background(), islice2)
	if err != nil {
		t.Errorf("Error inserting documents: %s", err.Error())
		t.FailNow()
	}
}

func TestFilterSimple(t *testing.T) {
	// test filter
	//   should return single record i.e. testuser{4, "test2", 37, 7}
	filter := bq.Builder().
		And(bq.C().EQ("name", "test2"), bq.C().GT("age", 29)).
		Build()
	exp := 1
	cur, err := c.Find(context.Background(), filter)
	if err != nil {
		t.Errorf("Error finding docs with filter: %s", err.Error())
		t.FailNow()
	}
	cnt := countDocs(cur)
	if cnt != exp {
		t.Errorf("TestFilterSimple failed. Expected %d,  Got: %d", exp, cnt)
	}
}

func TestFilterOr(t *testing.T) {
	// or query (age > 40 or score > 7)
	// documents satisfying any condition will be selected
	//   satisfied records are
	//		testuser{1, "test1", 44, 9}
	//		testuser{3, "test3Reg", 29, 9}
	//		testuser{4, "test3Big", 42, 5}
	/*
		{ "$or": [ {"age": {"$gt": 30}}, {"score": {"$gt": 5}} ] }
	*/
	exp := 3
	filter := bq.Builder().
		Or(bq.C().GT("age", 40), bq.C().GT("score", 7)).
		Build()
	cur, err := c.Find(context.Background(), filter)
	if err != nil {
		t.Errorf("Error finding docs with filter: %s", err.Error())
		t.FailNow()
	}
	cnt := countDocs(cur)
	if cnt != exp {
		t.Errorf("TestFilterOr failed. Expected %d,  Got: %d", exp, cnt)
	}
}

func TestFilterRegex(t *testing.T) {
	// check Regex, name like %3Bi%
	exp := 2
	filter := bq.Builder().
		And(bq.C().Regex("name", "3Bi", false)).
		Build()
	cur, err := c.Find(context.Background(), filter)
	if err != nil {
		t.Errorf("Error finding docs with filter: %s", err.Error())
		t.FailNow()
	}
	cnt := countDocs(cur)
	if cnt != exp {
		t.Errorf("TestFilterRegex failed. Expected %d,  Got: %d", exp, cnt)
	}
}

func TestFilterExist(t *testing.T) {
	// find documents having field 'extra', field must exist, field-value could be null
	exp := 2
	filter := bq.Builder().
		And(bq.C().Exist("extra", true)).
		Build()
	cur, err := c.Find(context.Background(), filter)
	if err != nil {
		t.Errorf("Error finding docs with filter: %s", err.Error())
		t.FailNow()
	}
	cnt := countDocs(cur)
	if cnt != exp {
		t.Errorf("TestFilterExist failed. Expected %d,  Got: %d", exp, cnt)
	}
}

func TestCleanup(t *testing.T) {
	client.Disconnect(context.Background())
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
