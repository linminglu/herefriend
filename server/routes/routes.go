package routes

import (
	"github.com/go-martini/martini"

	"herefriend/server/cms"
	"herefriend/server/handlers"
)

/*
 *
 *    Function: InstallRoutes
 *      Author: sunchao
 *        Date: 15/6/20
 * Description: install the routes of the web server
 *
 */
func InstallRoutes(m *martini.ClassicMartini) {
	/*
	 * Base
	 */
	m.Get("/User/Register", handlers.Register)
	m.Get("/User/Login", handlers.Login)
	m.Get("/User/Logout", handlers.Logout)
	m.Get("/User/WatchDog", handlers.WatchDog)
	m.Get("/User/SetProfile", handlers.SetProfile)
	m.Get("/User/GetDistrict", handlers.GetDistrict)
	m.Get("/User/GetPersonInfo", handlers.GetPersonInfo)
	m.Get("/User/Search", handlers.Search)
	m.Get("/User/Heartbeat", handlers.Heartbeat)
	m.Post("/User/PostImage", handlers.PostImage)
	m.Get("/User/DelImage", handlers.DeleteImage)
	m.Get("/vip/price", handlers.VipPrice)
	m.Get("/User/BuyVIP", handlers.BuyVip)

	/*
	 * Comments
	 */
	m.Get("/Action/Recommend", handlers.ActionRecommend)
	m.Get("/Action/DelRecommend", handlers.DelRecommend)
	m.Get("/User/Visit", handlers.DoVisit)
	m.Get("/User/ReadVisit", handlers.ReadVisit)
	m.Get("/User/DelVisit", handlers.DeleteVisit)
	m.Get("/User/WaterFlow", handlers.GetWaterFlow)
	m.Get("/User/AllMessage", handlers.GetAllMessage)
	m.Get("/User/UnreadMessage", handlers.GetUnreadMessage)
	m.Get("/User/Report", handlers.Report)
	m.Get("/User/AddBlacklist", handlers.UserAddBlacklist)
	m.Get("/User/DelBlacklist", handlers.UserDelBlacklist)
	m.Get("/User/GetBlacklist", handlers.UserGetBlacklist)

	/*
	 * Gift
	 */
	m.Get("/Gift/GoldPrice", handlers.GoldPrice)
	m.Get("/Gift/BuyBeans", handlers.BuyBeans)
	m.Get("/Gift/GiftList", handlers.GiftList)
	m.Get("/Gift/PresentGift", handlers.PresentGift)
	m.Get("/Gift/RecvListVerbose", handlers.RecvListVerbose)
	m.Get("/Gift/SendListVerbose", handlers.SendListVerbose)
	m.Get("/User/CharmTopList", handlers.CharmTopList)

	/*
	 * CMS
	 */
	m.Get("/login", cms.CmsLogin)
	m.Get("/cms/log", cms.Log)
	m.Post("/dashboard", cms.CmsDashBoard)
	m.Get("/cms/sysinfo", cms.SystemInfo)
	m.Get("/cms/cpuinfo", cms.CpuInfo)
	m.Get("/cms/sysuserinfo", cms.SystemUserInfo)
	m.Get("/cms/commentinfo", cms.CommentInfo)
	m.Get("/cms/recommendhistory", cms.Recommendhistory)
	m.Get("/cms/msgtemplate", cms.MsgTemplate)
	m.Get("/cms/msgtemplateadd", cms.MsgTemplateAdd)
	m.Get("/cms/msgtemplatedel", cms.MsgTemplateDel)
	m.Get("/cms/msgtemplatemodify", cms.MsgTemplateModify)
	m.Get("/cms/SearchUserInfos", cms.SearchUserInfos)
	m.Get("/cms/GetUserInfos", cms.GetUserInfos)
	m.Get("/cms/GetSingleUserInfo", cms.GetSingleUserInfo)
	m.Get("/cms/SetSingleUserInfo", cms.SetSingleUserInfo)
	m.Get("/cms/SetHeartbeat", cms.SetHeartbeat)
	m.Get("/cms/ChangeHeadImage", cms.ChangeHeadImage)
	m.Get("/cms/DeleteHeadImage", cms.DeleteHeadImage)
	m.Get("/cms/AddBlacklist", cms.AddBlacklist)
	m.Get("/cms/RefreshUserInfo", cms.RefreshUserInfo)
	m.Get("/cms/RegistUserInfo", cms.RegistUserInfo)
	m.Get("/cms/GetChartsList", cms.GetChartsList)
	m.Get("/cms/GetTalkHistory", cms.GetTalkHistory)
	m.Get("/cms/DoTalk", cms.DoTalk)
	m.Get("/cms/MessagePushSet", cms.MessagePushSet)
	m.Get("/cms/PresentGift", cms.PresentGift)
	m.Get("/cms/RefreshGiftConsume", cms.RefreshGiftConsume)
}
