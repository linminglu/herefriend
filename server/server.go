package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"

	"herefriend/lib"
	"herefriend/server/cms"
	"herefriend/server/handlers"
)

func main() {
	defer lib.CloseSQL()

	f, err := os.Create("/var/run/herefriend.pid")
	if nil == err {
		f.WriteString(fmt.Sprintf("%d", os.Getpid()))
		f.Close()
	} else {
		fmt.Println(err)
	}

	r := gin.Default()
	/*
	 * Base
	 */
	userGroup := r.Group("/User")
	{
		userGroup.GET("/Register", handlers.Register)
		userGroup.GET("/Login", handlers.Login)
		userGroup.GET("/Logout", handlers.Logout)
		userGroup.GET("/WatchDog", handlers.WatchDog)
		userGroup.GET("/SetProfile", handlers.SetProfile)
		userGroup.GET("/GetDistrict", handlers.GetDistrict)
		userGroup.GET("/GetPersonInfo", handlers.GetPersonInfo)
		userGroup.GET("/Search", handlers.Search)
		userGroup.GET("/Heartbeat", handlers.Heartbeat)
		userGroup.GET("/DelImage", handlers.DeleteImage)
		userGroup.GET("/BuyVIP", handlers.BuyVip)
		userGroup.GET("/AppConfig", handlers.GetAppConfig)
		userGroup.GET("/WaterFlow", handlers.GetWaterFlow)
		userGroup.GET("/AllMessage", handlers.GetAllMessage)
		userGroup.GET("/UnreadMessage", handlers.GetUnreadMessage)
		userGroup.GET("/GetComments", handlers.GetComments)
		userGroup.GET("/GetVisits", handlers.GetVisits)
		userGroup.GET("/Visit", handlers.DoVisit)
		userGroup.GET("/ReadVisit", handlers.ReadVisit)
		userGroup.GET("/DelVisit", handlers.DeleteVisit)
		userGroup.GET("/Report", handlers.Report)
		userGroup.GET("/AddBlacklist", handlers.UserAddBlacklist)
		userGroup.GET("/DelBlacklist", handlers.UserDelBlacklist)
		userGroup.GET("/GetBlacklist", handlers.UserGetBlacklist)
		userGroup.GET("/CharmTopList", handlers.CharmTopList)
		userGroup.GET("/WealthTopList", handlers.WealthTopList)
		userGroup.POST("/PostImage", handlers.PostImage)
	}
	r.GET("/Action/Recommend", handlers.ActionRecommend)
	r.GET("/Action/DelRecommend", handlers.DelRecommend)
	r.GET("/vip/price", handlers.VipPrice)

	/*
	 * Gift
	 */
	giftGroup := r.Group("/Gift")
	{
		giftGroup.GET("/GoldPrice", handlers.GoldPrice)
		giftGroup.GET("/BuyBeans", handlers.BuyBeans)
		giftGroup.GET("/GiftList", handlers.GiftList)
		giftGroup.GET("/PresentGift", handlers.PresentGift)
		giftGroup.GET("/RecvListVerbose", handlers.RecvListVerbose)
		giftGroup.GET("/SendListVerbose", handlers.SendListVerbose)
	}

	/*
	 * html
	 */
	{
		r.StaticFile("/login", "public/signin.html")
		r.POST("/www/dashboard", cms.DashBoard)
		r.Static("/www", "public")
		r.Static("assets", "public/assets")
	}

	/*
	 * CMS
	 */
	cmsGroup := r.Group("/cms")
	{
		cmsGroup.GET("/log", cms.Log)
		cmsGroup.GET("/sysinfo", cms.SystemInfo)
		cmsGroup.GET("/cpuinfo", cms.CPUInfo)
		cmsGroup.GET("/sysuserinfo", cms.SystemUserInfo)
		cmsGroup.GET("/commentinfo", cms.CommentInfo)
		cmsGroup.GET("/recentComments", cms.RecentComments)
		cmsGroup.GET("/msgtemplate", cms.MsgTemplate)
		cmsGroup.GET("/msgtemplateadd", cms.MsgTemplateAdd)
		cmsGroup.GET("/msgtemplatedel", cms.MsgTemplateDel)
		cmsGroup.GET("/msgtemplatemodify", cms.MsgTemplateModify)
		cmsGroup.GET("/SearchUserInfos", cms.SearchUserInfos)
		cmsGroup.GET("/GetUserInfos", cms.GetUserInfos)
		cmsGroup.GET("/GetSingleUserInfo", cms.GetSingleUserInfo)
		cmsGroup.GET("/SetSingleUserInfo", cms.SetSingleUserInfo)
		cmsGroup.GET("/AdminGiveVipLevel", cms.AdminGiveVipLevel)
		cmsGroup.GET("/SetHeartbeat", cms.SetHeartbeat)
		cmsGroup.GET("/ChangeHeadImage", cms.ChangeHeadImage)
		cmsGroup.GET("/DeleteHeadImage", cms.DeleteHeadImage)
		cmsGroup.GET("/AddBlacklist", cms.AddBlacklist)
		cmsGroup.GET("/RefreshUserInfo", cms.RefreshUserInfo)
		cmsGroup.GET("/RegistUserInfo", cms.RegistUserInfo)
		cmsGroup.GET("/AdminChartsList", cms.AdminChartsList)
		cmsGroup.GET("/GetChartsList", cms.GetChartsList)
		cmsGroup.GET("/GetTalkHistory", cms.GetTalkHistory)
		cmsGroup.GET("/DoTalk", cms.DoTalk)
		cmsGroup.GET("/MessagePushSet", cms.MessagePushSet)
		cmsGroup.GET("/PresentGift", cms.PresentGift)
		cmsGroup.GET("/GetGiftList", cms.GetGiftList)
		cmsGroup.GET("/GetGiftVerbose", cms.GetGiftVerbose)
	}

	if os.Getenv("DEBUG") != "1" {
		gin.SetMode(gin.ReleaseMode)
	}

	r.Run(":8080")
}
