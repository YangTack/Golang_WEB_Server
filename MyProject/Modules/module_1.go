package Modules

import "gopkg.in/mgo.v2"

var DbSession *mgo.Session
var Db *mgo.Database

type FileList struct {
	FileName string `bson:"name"`
	FileSize string `bson:"size"`
	FilePath string `bson:"file_path"`
	MD5      string `bson:"md5"`
}
