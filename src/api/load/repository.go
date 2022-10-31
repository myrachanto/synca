package load

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	Loadrepository LoadRepoInterface = &loadrepository{}
	Repo                             = &loadrepository{}
	ctx                              = context.TODO()
	SyncTime                         = 10
	SyncTimeFrame                    = "Seconds"
)

type LoadRepoInterface interface {
	Synca() (bool, int)
	StartSychronization()
}
type loadrepository struct {
	Ago            int
	DatabaseA      string
	DatabaseAUrl   string
	Databaseb      string
	DatabasebURL   string
	CollectionName string
}

func NewloadRepo(dba, dbaurl, dbb, dbburl, table string) LoadRepoInterface {
	return &loadrepository{
		DatabaseA:      dba,
		DatabaseAUrl:   dbaurl,
		Databaseb:      dbb,
		DatabasebURL:   dbburl,
		CollectionName: table,
	}
}
func (r *loadrepository) StartSychronization() {
	// fmt.Println("Items initialized!", r.DatabaseA, r.DatabaseAUrl, r.Databaseb, r.DatabasebURL, r.CollectionName)
	res, itemscount := r.Synca()
	if !res {
		r.Ago = r.Ago + SyncTime
		resp := fmt.Sprintf("%d Items Were Sychronized successfully done at %d  %s ago \n", itemscount, r.Ago, SyncTimeFrame)
		r.RecordSynca(false, resp, itemscount)
		fmt.Println(resp)
		// emailing.Emails.Emailing(res)
	} else {
		r.Ago = r.Ago + SyncTime
		resp := fmt.Sprintf("Sychronization failed at %d  %s ago \n", r.Ago, SyncTimeFrame)
		r.RecordSynca(true, resp, itemscount)
		fmt.Println(resp)
		// emailing.Emails.Emailing(res)
	}
	_ = time.AfterFunc(time.Second*time.Duration(SyncTime), r.StartSychronization)
}
func (r *loadrepository) Synca() (bool, int) {
	lastsync, state := r.LastSynchronization()
	// log.Println("Synca ---------------------------- hello")
	counter := 0
	itemscount := 0
	res := true
	if !state {
		// do full synchronization
		// log.Println("Synca ---------------------------- hello step 1")
		for _, product := range r.DataFromDBA(false, time.Now()) {
			// log.Println("Synca ---------------------------- hello step 1a", product)
			exist := r.CheckIfExistDBB(false, product)
			// log.Println("Synca ---------------------------- hello step 1b", exist)
			if !exist {
			checka:
				for !r.InsertDataDBB(product) {
					r.InsertDataDBB(product)
					counter++
					itemscount++
					if counter >= 5 {
						res = false
						break checka
					}
				}
			}
		}
	} else {
		// log.Println("Synca ---------------------------- hello=================")
		dataFromA := r.DataFromDBA(true, lastsync.Dated)
	asdfs:
		for _, product := range dataFromA {
			exist := r.CheckIfExistDBB(false, product, lastsync.Dated)
			if !exist {
				for !r.InsertDataDBB(product) {
					r.InsertDataDBB(product)
					counter++
					itemscount++
					if counter >= 5 {
						res = false
						break asdfs
					}
				}
			}
		}
	}
	return res, itemscount
}
func (r *loadrepository) LastSynchronization() (*Synca, bool) {
	GormDB, err1 := CentralRepo.Getconnected()
	if err1 != nil {
		log.Fatal("Got an error trying to connect to sync db")
		return nil, false
	}
	defer CentralRepo.DbClose(GormDB)
	synca := &Synca{}
	res := GormDB.Where("database_a = ? AND database_b = ?", r.DatabaseA, r.Databaseb).Last(&synca)
	if res.Error != nil {
		return nil, false
	}
	return synca, true
}
func (r *loadrepository) DataFromDBA(status bool, dated ...time.Time) []*Product {
	// fmt.Println("dba -----------------step 1", r.DatabaseAUrl, r.DatabaseA)
	conn, err := Mongodb(r.DatabaseAUrl, r.DatabaseA)
	if err != nil {
		log.Panicln("could not connect to database A")
	}
	collection := conn.Collection(r.CollectionName)
	defer DbClose(conn.Client())
	// fmt.Println("dba -----------------step 1a")
	results := []*Product{}
	// fmt.Println("dba -----------------step 1b")
	if status {
		// fmt.Println("===================== status")
		filter := bson.M{
			"dated": bson.M{"$gte": dated[0]},
		}
		cursor, err := collection.Find(ctx, filter)
		if err != nil {
			return nil
		}
		if err = cursor.All(ctx, &results); err != nil {
			return nil
		}
		return results
	} else {
		// fmt.Println("===================== else")
		filter := bson.M{}
		cursor, err := collection.Find(ctx, filter)
		if err != nil {
			return nil
		}
		if err = cursor.All(ctx, &results); err != nil {
			return nil
		}
		return results

	}
}
func (r *loadrepository) CheckIfExistDBB(status bool, product *Product, dated ...time.Time) bool {
	// fmt.Println("----------------exist cool 1")
	conn, err := Mongodb(r.DatabasebURL, r.Databaseb)
	if err != nil {
		log.Panicln("could not connect to database B")
	}
	collection := conn.Collection(r.CollectionName)
	defer DbClose(conn.Client())
	result := &Product{}
	if status {
		filter := bson.M{
			"dated": bson.M{"$gte": dated[0]},
			"url":   product.Name,
		}
		// fmt.Println("----------------exist cool 2")
		err1 := collection.FindOne(ctx, filter).Decode(&result)
		return err1 == nil
	}
	// fmt.Println("----------------exist cool 3")
	filter := bson.M{"code": product.Code}
	err1 := collection.FindOne(ctx, filter).Decode(&result)
	return err1 == nil
}
func (r *loadrepository) InsertDataDBB(product *Product) bool {
	fmt.Println("++++++++++++++++++++++++ insert db step 1")
	conn, err := Mongodb(r.DatabasebURL, r.Databaseb)
	if err != nil {
		log.Panicln("could not Insert to database B")
	}
	fmt.Println("++++++++++++++++++++++++", product)
	collection := conn.Collection(r.CollectionName)
	defer DbClose(conn.Client())
	product.ID = primitive.NilObjectID
	res, err1 := collection.InsertOne(ctx, product)
	fmt.Println("++++++++++++++++++++++++", res.InsertedID)
	return err1 == nil
}
func (r *loadrepository) RecordSynca(status bool, message string, itemscount int) bool {
	GormDB, err1 := CentralRepo.Getconnected()
	if err1 != nil {
		log.Fatal("Got an error trying to connect to sync db")
		return false
	}
	defer CentralRepo.DbClose(GormDB)
	name := r.GeneCode()
	synca := &Synca{Name: name, DatabaseA: r.DatabaseA, DatabaseB: r.Databaseb, Status: status, Message: message, Items: itemscount, Dated: time.Now()}
	res := GormDB.Create(&synca)
	return res.Error == nil
}

func (r *loadrepository) GeneCode() string {
	GormDB, err1 := CentralRepo.Getconnected()
	if err1 != nil {
		log.Fatal("Got an error trying to connect to sync db")
		return ""
	}
	defer CentralRepo.DbClose(GormDB)
	synca := &Synca{}
	special := Stamper()
	err := GormDB.Last(&synca)
	if err.Error != nil {
		var c1 uint = 1
		code := "SyncRec" + strconv.FormatUint(uint64(c1), 10) + "-" + special
		return code
	}
	c1 := synca.ID + 1
	code := "SyncRec" + strconv.FormatUint(uint64(c1), 10) + "-" + special

	return code

}

func Stamper() string {
	timestamp := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	special := timestamp[0:7]
	return special
}
