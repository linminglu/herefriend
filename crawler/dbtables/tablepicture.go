package dbtables

import (
	"fmt"

	"herefriend/config"
	"herefriend/lib"
)

// the fmt sentences
const g_createtable_picture_fmt = "create table if not exists %s_picture (" +
	"id int not null," +
	"filename varchar(255) character set 'utf8' not null," +
	"tag int," +
	"primary key (id, filename))"

const g_Girlspicture_OldDelete = "delete from girls_picture where id=? and tag=1"
const g_Guyspicture_OldDelete = "delete from guys_picture where id=? and tag=1"
const g_Girlspicture_OldHead = "select id,filename from girls_picture where tag=1 and renew=0 order by id desc limit 1000"
const g_Guyspicture_OldHead = "select id,filename from guys_picture where tag=1 and renew=0 order by id desc limit 1000"
const g_Girlspicture_Insert = "insert into girls_picture (id, filename, tag, renew) value (?,?,?,?)"
const g_Guyspicture_Insert = "insert into guys_picture (id, filename, tag, renew) value (?,?,?,?)"

// the sql sentences
var g_droptable_picture string
var g_createtable_picture string

func preparePictureSentences() {
	g_createtable_picture = fmt.Sprintf(g_createtable_picture_fmt, "girls")
}

/*
 *
 *    Function: PrepareTablePicture
 *      Author: sunchao
 *        Date: 15/7/12
 * Description: check and create the picture table
 *
 */
func PrepareTablePicture() {
	_, err := lib.SQLExec(g_createtable_picture)
	if nil != err {
		fmt.Println(err.Error())
	}

	return
}

/*
 * insert picture files' name to db by id
 */
func InsertPictureById(id, gender int, file string, bHead, bNew bool) {
	sqlstr := func() string {
		if 0 == gender {
			return g_Girlspicture_Insert
		} else {
			return g_Guyspicture_Insert
		}
	}()

	tag := func() int {
		if true == bHead {
			return 1
		} else {
			return 0
		}
	}()

	renew := func() int {
		if true == bNew {
			return 1
		} else {
			return 0
		}
	}()

	lib.SQLExec(sqlstr, id, file, tag, renew)
}

func DeleteOldHeadPictureById(id, gender int) {
	sqlstr := func() string {
		if 0 == gender {
			return g_Girlspicture_OldDelete
		} else {
			return g_Guyspicture_OldDelete
		}
	}()

	lib.SQLExec(sqlstr, id)
}

func DoRenewHeadPicture(gender int, callback func(id, gender int, filename string)) {
	var id int
	var filename string
	var iNum int

	sqlstr := func() string {
		if 0 == gender {
			return g_Girlspicture_OldHead
		} else {
			return g_Guyspicture_OldHead
		}
	}()

	for {
		rows, err := lib.SQLQuery(sqlstr)
		if nil != err {
			fmt.Println(err)
			break
		}

		iNum = 0

		for rows.Next() {
			iNum++
			err = rows.Scan(&id, &filename)
			if nil == err {
				callback(id, gender, filename)
			} else {
				fmt.Println(err)
			}
		}

		rows.Close()

		if 0 == iNum {
			break
		}
	}
}
