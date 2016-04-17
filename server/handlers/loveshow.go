package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"herefriend/config"
	"herefriend/lib"
)

const (
	g_LoveShowCount = 177
)

func getLoveShowPics(loveshowid int) []string {
	sentence := lib.SQLSentence(lib.SQLMAP_Select_LoveshowPicture)
	rows, err := lib.SQLQuery(sentence, loveshowid)
	if nil != err {
		return nil
	}
	defer rows.Close()

	tag := lib.GetQiniuLoveShowPicturePrefix(loveshowid)

	imgs := make([]string, 0)
	for rows.Next() {
		var image string

		err = rows.Scan(&image)
		if nil == err {
			imgs = append(imgs, tag+image)
		}
	}

	return imgs
}

func getLoveShowComments(loveshowid int) []loveshowcomment {
	sentence := lib.SQLSentence(lib.SQLMAP_Select_LoveshowBless)
	rows, err := lib.SQLQuery(sentence, loveshowid)
	if nil != err {
		return nil
	}
	defer rows.Close()

	var comments []loveshowcomment
	for rows.Next() {
		var comment loveshowcomment
		var timeSec int64

		err = rows.Scan(&comment.Id, &comment.Age, &timeSec, &comment.Name, &comment.District, &comment.Education, &comment.Text)
		if nil == err {
			comment.TimeUTC = lib.Int64_To_UTCTime(timeSec)
			comments = append(comments, comment)
		}
	}

	return comments
}

/*
 *
 *    Function: GetLoveShow
 *      Author: sunchao
 *        Date: 15/10/3
 * Description: 获取恋爱秀列表
 *
 */
func GetLoveShow(req *http.Request) (int, string) {
	exist, _, _ := getIdGenderByRequest(req)
	if true != exist {
		return 404, ""
	}

	count := lib.GetCountRequestArgument(req)
	sentence := lib.SQLSentence(lib.SQLMAP_Select_LoveshowByRows)
	rows, err := lib.SQLQuery(sentence, lib.Intn(g_LoveShowCount-count), count)
	if nil != err {
		return 404, ""
	}
	defer rows.Close()

	infos := make([]loveShow, 0)
	for rows.Next() {
		var info loveShow
		var timeSec int64

		err = rows.Scan(&info.Id, &timeSec, &info.Blessnum, &info.Daysfalllove, &info.Girl.Id, &info.Guy.Id, &info.Girl.Age, &info.Guy.Age,
			&info.Girl.Name, &info.Guy.Name, &info.Girl.Imgurl, &info.Guy.Imgurl, &info.Girl.District, &info.Guy.District, &info.Lovestatus,
			&info.Lovetitle, &info.Lovestory)
		if nil == err {
			tag := config.Conf_QiniuPre + "loveshow/" + strconv.Itoa(info.Id) + "/"

			if "" != info.Girl.Imgurl {
				info.Girl.Imgurl = tag + info.Girl.Imgurl
			}

			if "" != info.Guy.Imgurl {
				info.Guy.Imgurl = tag + info.Guy.Imgurl
			}

			info.TimeUTC = lib.Int64_To_UTCTime(timeSec)

			/* images */
			info.ShowPics = getLoveShowPics(info.Id)

			/* comments */
			info.Comments = getLoveShowComments(info.Id)

			infos = append(infos, info)
		}
	}

	jsonRlt, _ := json.Marshal(infos)
	return 200, string(jsonRlt)
}

/*
 *
 *    Function: LoveShowComment
 *      Author: sunchao
 *        Date: 15/10/3
 * Description: 恋爱秀送祝福
 *
 */
func LoveShowComment(req *http.Request) (int, string) {
	v := req.URL.Query()
	loveshowidstr := v.Get("loveshowid")
	if "" == loveshowidstr {
		return 404, ""
	}

	exist, id, gender := getIdGenderByRequest(req)
	if true != exist {
		return 404, ""
	}

	loveshowid, _ := strconv.Atoi(loveshowidstr)
	timeSec := lib.Time_To_UTCInt64(time.Now())
	_, info := GetUserInfo(id, gender)

	sentence := lib.SQLSentence(lib.SQLMAP_Insert_LoveshowBless)
	_, err := lib.SQLExec(sentence, loveshowid, id, info.Age, timeSec, info.Name, info.Province+info.District,
		info.Education, v.Get("bless"))
	if nil != err {
		return 404, ""
	}

	comments := getLoveShowComments(loveshowid)
	jsonRlt, _ := json.Marshal(comments)
	return 200, string(jsonRlt)
}

/*
 *
 *    Function: ListLoveShow
 *      Author: sunchao
 *        Date: 15/10/4
 * Description: 列出所有的恋爱秀
 *
 */
func ListLoveShow(req *http.Request) (int, string) {
	pageid, count := lib.Get_pageid_count_fromreq(req)
	sentence := lib.SQLSentence(lib.SQLMAP_Select_LoveshowWithHideByRows)
	rows, err := lib.SQLQuery(sentence, (pageid-1)*count, count)
	if nil != err {
		return 404, ""
	}
	defer rows.Close()

	infos := make([]loveShowList, 0)
	for rows.Next() {
		var info loveShowList
		var timeSec int64

		err = rows.Scan(&info.Id, &timeSec, &info.Blessnum, &info.Daysfalllove, &info.Girl.Id, &info.Guy.Id, &info.Girl.Age, &info.Guy.Age,
			&info.Girl.Name, &info.Guy.Name, &info.Girl.Imgurl, &info.Guy.Imgurl, &info.Girl.District, &info.Guy.District, &info.Lovestatus,
			&info.Lovetitle, &info.Lovestory, &info.Hide)
		if nil == err {
			tag := config.Conf_QiniuPre + "loveshow/" + strconv.Itoa(info.Id) + "/"

			if "" != info.Girl.Imgurl {
				info.Girl.Imgurl = tag + info.Girl.Imgurl
			}

			if "" != info.Guy.Imgurl {
				info.Guy.Imgurl = tag + info.Guy.Imgurl
			}

			info.TimeUTC = lib.Int64_To_UTCTime(timeSec)

			/* images */
			info.ShowPics = getLoveShowPics(info.Id)

			/* comments */
			info.Comments = getLoveShowComments(info.Id)

			infos = append(infos, info)
		}
	}

	jsonRlt, _ := json.Marshal(infos)
	return 200, string(jsonRlt)
}

/*
 *
 *    Function: ReplaceLoveShow
 *      Author: sunchao
 *        Date: 15/10/4
 * Description: 替换恋爱秀中的用户
 *
 */
func ReplaceLoveShow(req *http.Request) (int, string) {
	v := req.URL.Query()
	loveshowidstr := v.Get("loveshowid")
	replaceidstr := v.Get("replaceid")
	genderstr := v.Get("gender")

	if "" == loveshowidstr || "" == replaceidstr || "" == genderstr {
		return 404, ""
	}

	loveshowid, _ := strconv.Atoi(loveshowidstr)
	replaceid, _ := strconv.Atoi(replaceidstr)
	gender, _ := strconv.Atoi(genderstr)

	code, info := GetUserInfo(replaceid, gender)
	if 200 == code {
		var updatestr = lib.SQLSentence(lib.SQLMAP_Update_LoveshowGirl)
		if 1 == gender {
			updatestr = lib.SQLSentence(lib.SQLMAP_Update_LoveshowGirl)
		}

		imagename, _ := lib.DownloadImgAndRename(info.IconUrl, "loveshow/"+loveshowidstr)
		_, err := lib.SQLExec(updatestr, info.Id, info.Age, info.Name, imagename, info.Province+info.District, loveshowid)
		if nil != err {
			return 404, ""
		} else {
			return 200, ""
		}
	} else {
		return 404, ""
	}
}

/*
 *
 *    Function: HideLoveShow
 *      Author: sunchao
 *        Date: 15/10/4
 * Description: 隐藏恋爱秀
 *
 */
func HideLoveShow(req *http.Request) (int, string) {
	v := req.URL.Query()
	loveshowidstr := v.Get("loveshowid")
	hidestr := v.Get("hide")

	if "" == loveshowidstr || "" == hidestr {
		return 404, ""
	}

	loveshowid, _ := strconv.Atoi(loveshowidstr)
	hide, _ := strconv.Atoi(hidestr)

	sentence := lib.SQLSentence(lib.SQLMAP_Update_HideLoveshow)
	_, err := lib.SQLExec(sentence, hide, loveshowid)
	if nil != err {
		return 404, ""
	} else {
		return 200, ""
	}
}
