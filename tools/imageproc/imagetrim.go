package main

import (
	"fmt"

	"herefriend/crawler/image"
	"herefriend/lib"
)

func trimImagesByGender(gender int, baseid int) {
	var count int
	var id int
	var filename string
	var tag int

	sentense := lib.SQLSentence(lib.SQLMapSelectSearchPicturesByFlag, gender)
	updatesentense := lib.SQLSentence(lib.SQLMapUpdateSetPictureFlag, gender)
	deletesentence := lib.SQLSentence(lib.SQLMapDeletePicture, gender)
	for {
		count = 0
		rows, err := lib.SQLQuery(sentense, baseid, 0, 1000)
		if nil != err {
			fmt.Println(err)
			break
		}

		for rows.Next() {
			err = rows.Scan(&id, &filename, &tag)
			if nil != err {
				fmt.Println(err)
				continue
			}

			count = count + 1
			err = image.DownloadImageAndPutToQiniu(lib.GetQiniuUserImageURL(id, filename), true, id, filename)
			if nil == err {
				_, err = lib.SQLExec(updatesentense, id, filename, tag)
				if nil != err {
					fmt.Println(err)
				}
			} else {
				fmt.Println(err)
				_, err = lib.SQLExec(deletesentence, id, filename)
				if nil != err {
					fmt.Println(err)
				} else {
					lib.DeleteImageFromQiniu(id, filename)
					fmt.Println(fmt.Sprintf("delete image %s of %d\n", filename, id))
				}
			}
		}

		rows.Close()

		if 0 == count {
			break
		}
	}
}

func main() {
	trimImagesByGender(0, 124493868)
	trimImagesByGender(1, 124493897)
}
