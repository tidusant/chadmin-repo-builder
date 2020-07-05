package builder

import (
	"time"

	c3mcommon "github.com/tidusant/chadmin-common"
	"github.com/tidusant/chadmin-log"
	"github.com/tidusant/chadmin-repo-models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"os"
)

var (
	db *mgo.Database
)

func init() {
	log.Infof("init repo build")
	strErr := ""
	db, strErr = c3mcommon.ConnectDB("chbuild")
	if strErr != "" {
		log.Infof(strErr)
		os.Exit(1)
	}
}

//query and update https://stackoverflow.com/questions/11417784/mongodb-in-go-golang-with-mgo-how-do-i-update-a-record-find-out-if-update-wa
func CreateBuild(bs models.BuildScript) string {
	col := db.C("builds")
	//remove old build
	cond := bson.M{"status": 0, "shopid": bs.ShopID, "objectid": bs.ObjectID}
	//"objectid": buildscript.ObjectID, "collection": buildscript.Collection}
	if bs.ObjectID != "home" && bs.ObjectID != "script" {
		cond["collection"] = bs.Collection
	}
	//change := bson.M{"$set": bson.M{"status": -1}}
	_, err := col.RemoveAll(cond)

	bs.Status = 0
	bs.Created = time.Now().Unix()
	bs.Modified = time.Now().Unix()
	err = col.Insert(bs)
	c3mcommon.CheckError("insert build script", err)
	if bs.ObjectID != "home" {
		bs.ObjectID = "home"
		CreateBuild(bs)
	}

	return ""
}

func GetBuild() models.BuildScript {
	col := db.C("builds")
	var bs models.BuildScript
	change := mgo.Change{
		Update:    bson.M{"$set": bson.M{"status": 1, "modified": time.Now().Unix()}},
		ReturnNew: true,
	}
	_, err := col.Find(bson.M{"status": 0}).Apply(change, &bs)
	c3mcommon.CheckError("GetBuild script", err)
	return bs

}

func RemoveBuild(shopID string) string {
	col := db.C("builds")
	cond := bson.M{"status": 0, "shopid": shopID}

	_, err := col.RemoveAll(cond)
	c3mcommon.CheckError("insert build script", err)
	return ""
}
