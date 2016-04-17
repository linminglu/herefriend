package pageweibo

import (
	"herefriend/crawler/request"
)

type PagePic struct {
	crawled bool
	headimg string
	imgs    []string
	request *request.Request
}

func NewPagePic(request *request.Request) *PagePic {
	return &PagePic{request: request}
}

func GetHeadPic(this *PagePic) string {
	return this.headimg
}

/*
 * the crawl action, get all the <a> and <img> from the page
 */
func (this *PagePic) Crawl() *PagePic {
	if true == this.crawled {
		return this
	}

	err := this.request.Download(nil)
	if nil != err {
		return this
	}

	query, err := this.request.CreateQuery()
	if nil != err {
		return this
	}

	var val string
	var exist bool

	/* headpic */
	sel := query.Find(".picface > a > img")
	if sel.Length() > 0 {
		val, exist = sel.Attr("src")
		if true == exist {
			this.headimg = val
		}
	}

	/*
	 * crawl the pictures
	 */
	//轮播pictures
	sel = query.Find(".photo_turn > ul > li")
	if n := sel.Length(); n > 0 {
		for i := 0; i < n; i++ {
			val, exist = sel.Eq(i).Find("a > img").Attr("src")
			if true == exist {
				this.imgs = append(this.imgs, val)
			}
		}
	}

	//双开pictures
	sel = query.Find(".vg_img_box")
	if n := sel.Length(); n > 0 {
		for i := 0; i < n; i++ {
			val, exist = sel.Eq(i).Find("img").Attr("src")
			if true == exist {
				this.imgs = append(this.imgs, val)
			}
		}
	}

	//photo_list
	sel = query.Find(".photo_list > li")
	if n := sel.Length(); n > 0 {
		for i := 0; i < n; i++ {
			val, exist = sel.Eq(i).Find("a > img").Attr("src")
			if true == exist {
				this.imgs = append(this.imgs, val)
			}
		}
	}

	if "" == this.headimg && len(this.imgs) > 0 {
		this.headimg = this.imgs[0]
	}

	this.crawled = true
	return this
}
