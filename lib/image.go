package lib

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/net/context"
	"qiniupkg.com/api.v7/kodo"

	"herefriend/config"
)

var gKodoClient *kodo.Client
var gBucket kodo.Bucket
var gCtx context.Context

func init() {
	kodo.SetMac(config.ConfQiniuAccessKey, config.ConfQiniuSecretKey)
	gKodoClient = kodo.New(0, nil)
	gBucket = gKodoClient.Bucket(config.ConfQiniuScope)
	gCtx = context.Background()
}

func getQiniuUserImagePath(id int, filename string) string {
	return fmt.Sprintf("images/%d/%s", id, filename)
}

// GetQiniuUserImageURL .
func GetQiniuUserImageURL(id int, filename string) string {
	return config.ConfQiniuPre + fmt.Sprintf("images/%d/%s", id, filename)
}

// GetQiniuGiftImageURL .
func GetQiniuGiftImageURL(filename string) string {
	return config.ConfQiniuPre + fmt.Sprintf("images/gift/%s", filename)
}

// GetQiniuLoveShowPicturePrefix .
func GetQiniuLoveShowPicturePrefix(loveshowid int) string {
	return config.ConfQiniuPre + fmt.Sprintf("loveshow/%d/", loveshowid)
}

// PutImageToQiniuByPath .
func PutImageToQiniuByPath(path string, data io.Reader) error {
	buf, err := ioutil.ReadAll(data)
	if nil != err {
		return err
	}

	r := bytes.NewReader(buf)
	err = gBucket.Put(gCtx, nil, path, r, r.Size(), nil)
	if err != nil {
		fmt.Println(err)
	}

	return err
}

// PutImageToQiniu put image to Qiniu with path
func PutImageToQiniu(id int, filename string, data io.Reader) error {
	return PutImageToQiniuByPath(getQiniuUserImagePath(id, filename), data)
}

// DeleteImageFromQiniu delete image from Qiniu with path
func DeleteImageFromQiniu(id int, filename string) error {
	err := gBucket.Delete(gCtx, getQiniuUserImagePath(id, filename))
	if nil != err {
		fmt.Println(err)
	}

	return err
}

// DownloadImgAndRename download the url as image and rename to new image
func DownloadImgAndRename(url string, prefix string) (string, error) {
	if "" == url {
		return "", nil
	}

	resp, err := Get(url, nil)
	if nil != err {
		return "", err
	}

	defer resp.Body.Close()

	strslice := strings.Split(url, "/")
	imgname := strslice[len(strslice)-1]
	if "" != imgname {
		err = PutImageToQiniuByPath(filepath.Join(prefix, imgname), resp.Body)
	}

	return imgname, err
}

const letterBytes = "0123456789klmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// RandStringBytesMaskImprSrc 生成固定长度的随机字符串
func RandStringBytesMaskImprSrc(n int) string {
	var src = rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
