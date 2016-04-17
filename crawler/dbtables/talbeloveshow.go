package dbtables

import (
	"fmt"
	"time"

	"herefriend/lib"
)

// the fmt sentences
const g_createtable_loveshow = "create table if not exists loveshow (" +
	"id int not null auto_increment," +
	"loveshowid int not null," +
	"time bigint not null," +
	"blessnum int not null," +
	"falldays int not null," +
	"girl_id int not null," +
	"guy_id int not null," +
	"girl_age int not null," +
	"guy_age int not null," +
	"girl_name varchar(45) character set 'utf8'," +
	"guy_name varchar(45) character set 'utf8'," +
	"girl_headimg varchar(255) character set 'utf8'," +
	"guy_headimg varchar(255) character set 'utf8'," +
	"girl_district varchar(255) character set 'utf8'," +
	"guy_district varchar(255) character set 'utf8'," +
	"status varchar(20) character set 'utf8'," +
	"title varchar(50) character set 'utf8'," +
	"story varchar(2048) character set 'utf8'," +
	"primary key (id))"

const g_createtalbe_loveshowpicture = "create table if not exists loveshowpicture (" +
	"id int not null auto_increment," +
	"loveshowid int not null," +
	"filename varchar(255) character set 'utf8'," +
	"primary key (id))"

const g_createtalbe_loveshowbless = "create table if not exists loveshowbless (" +
	"id int not null auto_increment," +
	"loveshowid int not null," +
	"uid int not null," +
	"age int not null," +
	"time bigint not null," +
	"name varchar(45) character set 'utf8'," +
	"district varchar(255) character set 'utf8'," +
	"education varchar(255) character set 'utf8'," +
	"bless varchar(512) character set 'utf8'," +
	"primary key (id))"

const g_insert_loveshow = "insert into loveshow (loveshowid, time, blessnum, falldays, girl_id, guy_id, girl_age," +
	"guy_age, girl_name, guy_name, girl_headimg, guy_headimg, girl_district, guy_district, status, title, story) values " +
	"(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
const g_insert_loveshowpicture = "insert into loveshowpicture (loveshowid, filename) values (?,?)"
const g_insert_loveshowbless = "insert into loveshowbless (loveshowid, uid, age, time, name, district, education, bless) values " +
	"(?, ?, ?, ?, ?, ?, ?, ?)"
const g_check_loveshowexist = "select id from loveshow where loveshow_id=?"
const g_check_loveshowpictureexist = "select id from loveshowpicture where loveshowid=? and filename=?"
const g_check_loveshowblessexist = "select id from loveshowbless where loveshowid=? and uid=? and time=?"

/*
 *
 *    Function: PrepareTableLoveshow
 *      Author: sunchao
 *        Date: 15/7/12
 * Description: check and create the loveshow table
 *
 */
func PrepareTableLoveshow() {
	_, err := lib.SQLExec(g_createtable_loveshow)
	if nil != err {
		fmt.Println(err.Error())
	}

	_, err = lib.SQLExec(g_createtalbe_loveshowpicture)
	if nil != err {
		fmt.Println(err.Error())
	}

	_, err = lib.SQLExec(g_createtalbe_loveshowbless)
	if nil != err {
		fmt.Println(err.Error())
	}

	return
}

func IsLoveshowExist(loveshowid int) bool {
	var id int
	var bExist = false

	err := lib.SQLQueryRow(g_check_loveshowexist, loveshowid).Scan(&id)
	if nil == err && 0 != id {
		bExist = true
	}

	return bExist
}

func isLoveshowpictureExist(loveshowid int, filename string) bool {
	var id int
	var bExist = false

	err := lib.SQLQueryRow(g_check_loveshowpictureexist, loveshowid, filename).Scan(&id)
	if nil == err && 0 != id {
		bExist = true
	}

	return bExist
}

func isLoveshowblessExist(loveshowid, uid int, timeSec int64) bool {
	var id int
	var bExist = false

	err := lib.SQLQueryRow(g_check_loveshowblessexist, loveshowid, uid, timeSec).Scan(&id)
	if nil == err && 0 != id {
		bExist = true
	}

	return bExist
}

func InsertToLoveshow(loveshowid int, timeUTC time.Time, blessnum, falldays, girl_id, guy_id, girl_age, guy_age int,
	girl_name, guy_name, girl_headimg, guy_headimg, girl_district, guy_district, status, title, story string) error {
	_, err := lib.SQLExec(g_insert_loveshow, loveshowid, timeUTC.UTC().Unix(), blessnum, falldays, girl_id, guy_id,
		girl_age, guy_age, girl_name, guy_name, girl_headimg, guy_headimg, girl_district, guy_district, status, title, story)
	if nil != err {
		fmt.Println(err)
	}

	return err

}

func InsertToLoveshowpicture(loveshowid int, filename string) error {
	if true != isLoveshowpictureExist(loveshowid, filename) {
		_, err := lib.SQLExec(g_insert_loveshowpicture, loveshowid, filename)
		if nil != err {
			fmt.Println(err)
		}

		return err
	}

	return nil
}

func InsertToLoveshowBless(loveshowid, uid, age int, timeUTC time.Time, name, district, education, bless string) error {
	timeSec := timeUTC.UTC().Unix()

	if true != isLoveshowblessExist(loveshowid, uid, timeSec) {
		_, err := lib.SQLExec(g_insert_loveshowbless, loveshowid, uid, age, timeSec, name, district, education, bless)
		if nil != err {
			fmt.Println(err)
		}

		return err
	}

	return nil
}
