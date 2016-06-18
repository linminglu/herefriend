package image

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"os"

	"herefriend/lib"

	"github.com/gographics/imagick/imagick"
)

var g_mw *imagick.MagickWand

func init() {
	imagick.Initialize()
	g_mw = imagick.NewMagickWand()
}

/*
 |    Function: DownloadImageAndPutToQiniu
 |      Author: Mr.Sancho
 |        Date: 2016-01-12
 |   Arguments:
 |      Return:
 | Description: download the given url as image and put it to Qiniu
 |
*/
func DownloadImageAndPutToQiniu(url string, cut bool, id int, filename string) error {
	fmt.Println("【Download】" + url)
	resp, err := lib.Get(url, nil)
	if nil != err {
		return err
	}
	defer resp.Body.Close()

	if 200 != resp.StatusCode {
		fmt.Printf("【Download】failed, statuscode=%d\n", resp.StatusCode)
		return errors.New("")
	}

	var data io.Reader
	if true == cut {
		var m image.Image
		var subImg image.Image

		m, _, err = image.Decode(resp.Body)
		if err != nil {
			fmt.Println(err)
			return err
		}

		bounds := m.Bounds()
		switch m.(type) {
		case *image.Gray:
			grayImg := m.(*image.Gray)
			subImg = grayImg.SubImage(image.Rect(0, 0, bounds.Max.X, (bounds.Max.Y - 60)))
		case *image.YCbCr:
			rgbImg := m.(*image.YCbCr)
			subImg = rgbImg.SubImage(image.Rect(0, 0, bounds.Max.X, (bounds.Max.Y - 60)))
		}

		var buf []byte
		buffer := bytes.NewBuffer(buf)
		jpeg.Encode(buffer, subImg, nil)

		data = io.Reader(buffer)
	} else {
		data = io.Reader(resp.Body)
	}

	//create a template file to buffer the image
	tmpfile, err := ioutil.TempFile(os.TempDir(), "herefriend_lib_tmpfile")
	if err != nil {
		fmt.Println(err)
		return err
	}

	tmpfilename := tmpfile.Name()
	defer os.Remove(tmpfilename)

	_, err = io.Copy(tmpfile, data)
	if nil != err {
		return err
	}
	tmpfile.Close()

	//remove the edges of the image
	err = g_mw.ReadImage(tmpfilename)
	if nil != err {
		fmt.Println(err)
		return err
	}
	defer g_mw.Clear()
	err = g_mw.TrimImage(15)
	if nil != err {
		fmt.Println(err)
		return err
	}
	tmpfile.Seek(0, os.SEEK_SET)
	err = g_mw.WriteImage(tmpfilename)
	if nil != err {
		fmt.Println(err)
		return err
	}

	tmpfile, err = os.Open(tmpfilename)
	if nil != err {
		fmt.Println(err)
		return err
	}

	//put the image to Qiniu
	if url == lib.GetQiniuUserImageURL(id, filename) {
		lib.DeleteImageFromQiniu(id, filename)
	}

	err = lib.PutImageToQiniu(id, filename, tmpfile)

	tmpfile.Close()

	return err
}

func GetImageWidthHight(url string) (error, int, int) {
	resp, err := lib.Get(url, nil)
	if nil != err {
		return err, 0, 0
	}
	defer resp.Body.Close()

	if 200 != resp.StatusCode {
		fmt.Printf("【Download】failed, statuscode=%d\n", resp.StatusCode)
		return errors.New(""), 0, 0
	}

	var m image.Image
	m, _, err = image.Decode(resp.Body)
	if err != nil {
		fmt.Println(err)
		return err, 0, 0
	}

	bounds := m.Bounds()
	return nil, bounds.Max.X, bounds.Max.Y
}
