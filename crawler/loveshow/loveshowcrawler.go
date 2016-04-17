package loveshow

import (
	"strconv"

	"herefriend/crawler/request"
)

func StartCrawl() {
	pageid := 1
	end := false

	for true != end {
		list := loveshow.NewSpacelist(request.NewRequest(request.REQUESTURL_SPACELIST, pageid, nil)).Crawl()
		spaceids := list.GetSpaceIds()
		if nil == spaceids || 0 == len(spaceids) {
			end = true
		} else {
			for _, idStr := range spaceids {
				id, _ := strconv.Atoi(idStr)
				loveshow.NewSpace(request.NewRequest(request.REQUESTURL_SPACE, id, nil)).Crawl().Save()
			}
		}

		pageid = pageid + 1
	}

	return
}
