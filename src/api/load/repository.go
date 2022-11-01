package load

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	httperors "github.com/myrachanto/erroring"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	Loadrepository LoadRepoInterface = &loadrepository{}
	Repo                             = &loadrepository{}
	ctx                              = context.TODO()
	SyncTime                         = 30
	SyncTimeFrame                    = "Seconds"
)

type LoadRepoInterface interface {
	Synca() (bool, int)
	StartSychronization()
	GetAll() ([]*Synca, httperors.HttpErr)
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

	_, code := r.RecordSynca("started")
	res, itemscount := r.Synca()
	if !res {
		r.Ago = r.Ago + SyncTime
		resp := fmt.Sprintf("%d Items Were Sychronized successfully done at %d  %s ago", itemscount, r.Ago, SyncTimeFrame)
		r.UpdateSynca(code, "completed", resp, itemscount)
		fmt.Println(resp)
		// emailing.Emails.Emailing(res)
	} else {
		r.Ago = r.Ago + SyncTime
		resp := fmt.Sprintf("Sychronization failed at %d  %s ago \n", r.Ago, SyncTimeFrame)
		r.UpdateSynca(code, "failed", resp, itemscount)
		fmt.Println(resp)
		// emailing.Emails.Emailing(res)
	}
	_ = time.AfterFunc(time.Second*time.Duration(SyncTime), r.StartSychronization)
}
func (r *loadrepository) Synca() (bool, int) {
	lastsync, counts := r.LastSynchronization()
	counter := 0
	itemscount := 0
	res := true
	if counts == 1 {
		// do full synchronization
		for _, product := range r.DataFromDBA(false, time.Now()) {
			exist := r.CheckIfExistDBB(false, product)
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
func (r *loadrepository) LastSynchronization() (*Synca, int64) {
	GormDB, err1 := CentralRepo.Getconnected()
	if err1 != nil {
		log.Fatal("Got an error trying to connect to sync db")
		return nil, 0
	}
	defer CentralRepo.DbClose(GormDB)
	synca := &Synca{}
	var counter int64
	res := GormDB.Where("database_a = ? AND database_b = ?  AND ending = ? ", r.DatabaseA, r.Databaseb, "").Last(&synca)
	if res.Error != nil {
		return nil, 0
	}
	GormDB.Where("database_a = ? AND database_b = ?  AND ending = ? ", r.DatabaseA, r.Databaseb, "").Count(&counter)
	return synca, counter
}

func (r *loadrepository) DataFromDBA(status bool, dated ...time.Time) []*Product {
	if r.DatabaseA == "" {
		products, err := Fetcher(r.DatabaseAUrl)
		if err != nil {
			log.Println(err)
		}
		return products
	} else {
		return r.DataFromDBAs(status, dated...)
	}
}
func (r *loadrepository) DataFromDBAs(status bool, dated ...time.Time) []*Product {
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
			"created_at": bson.M{"$gte": dated[0]},
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
func (r *loadrepository) RecordSynca(status string) (bool, string) {
	GormDB, err1 := CentralRepo.Getconnected()
	if err1 != nil {
		log.Fatal("Got an error trying to connect to sync db")
		return false, ""
	}
	defer CentralRepo.DbClose(GormDB)
	name := r.GeneCode()
	dated := time.Now()
	synca := &Synca{Name: name, DatabaseA: r.DatabaseA, DatabaseB: r.Databaseb, Status: status, Dated: dated, Start: time.Now()}
	res := GormDB.Create(&synca)
	if res.Error != nil {
		return false, ""
	}
	return true, name
}
func (r *loadrepository) UpdateSynca(code string, status string, message string, itemscount int) bool {
	GormDB, err1 := CentralRepo.Getconnected()
	if err1 != nil {
		log.Fatal("Got an error trying to connect to sync db")
		return false
	}
	defer CentralRepo.DbClose(GormDB)
	synca := &Synca{}
	t := time.Now()
	dato := fmt.Sprintln(t.Format("2006-01-02 15:04:05"))
	// fmt.Println(dato)
	res := GormDB.Model(&synca).Where("name = ?", code).Updates(Synca{Ending: dato, Status: status, Message: message, Items: itemscount})
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

func (r *loadrepository) GetAll() ([]*Synca, httperors.HttpErr) {
	results := []*Synca{}
	GormDB, err1 := CentralRepo.Getconnected()
	if err1 != nil {
		log.Fatal("Got an error trying to connect to sync db")
		return nil, httperors.NewBadRequestError("Got an error trying to connect to sync db")
	}
	defer CentralRepo.DbClose(GormDB)
	GormDB.Order("dated desc").Find(&results)

	return results, nil
}
