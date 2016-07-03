#HereFriend接口列表





<!-- toc -->

* [1.说明](#1说明)
* [2.账户接口](#2账户接口)
  * [2.1 用户注册](#21-用户注册)
  * [2.2 用户登录](#22-用户登录)
  * [2.3 用户资料修改](#23-用户资料修改)
  * [2.4 上传图片](#24-上传图片)
  * [2.5 删除图片](#25-删除图片)
  * [2.6 用户退出](#26-用户退出)
  * [2.7 用户购买服务](#27-用户购买服务)
  * [2.8 获取用户信息](#28-获取用户信息)
  * [2.9 用户在线狗叫握手](#29-用户在线狗叫握手)
* [3.消息接口](#3消息接口)
  * [3.1 发送消息(打招呼/聊天)](#31-发送消息打招呼聊天)
  * [3.2 删除消息](#32-删除消息)
  * [3.3 查看资料](#33-查看资料)
  * [3.4 设置"查看资料"信息为已读](#34-设置查看资料信息为已读)
  * [3.5 删除"查看资料"信息](#35-删除查看资料信息)
  * [3.6 获取与指定用户的聊天消息](#36-获取与指定用户的聊天消息)
  * [3.7 获取聊天记录和访问记录](#37-获取聊天记录和访问记录)
  * [3.8 举报](#38-举报)
  * [3.9 添加用户黑名单](#39-添加用户黑名单)
  * [3.10 删除用户黑名单](#310-删除用户黑名单)
  * [3.11 查询用户黑名单](#311-查询用户黑名单)
  * [3.12 获取未读消息数量](#312-获取未读消息数量)
  * [3.13 获取消息列表](#313-获取消息列表)
  * [3.14 获取未读消息数量](#314-获取未读消息数量)
* [4.功能接口](#4功能接口)
  * [4.1 搜索用户](#41-搜索用户)
  * [4.2 心动列表](#42-心动列表)
  * [4.3 恋爱秀列表](#43-恋爱秀列表)
  * [4.4 恋爱秀送祝福](#44-恋爱秀送祝福)
  * [4.5 地区列表](#45-地区列表)
  * [4.6 VIP价格表](#46-vip价格表)
  * [4.7 检查新版本](#47-检查新版本)
* [5.推送消息](#5推送消息)
  * [5.1 透传消息推送](#51-透传消息推送)
* [6.JSON结构体](#6json结构体)
  * [6.1 personInfo](#61-personinfo)
  * [6.2 registerInfo](#62-registerinfo)
  * [6.3 聊天信息与访问资料信息](#63-聊天信息与访问资料信息)
  * [6.4 loveShow](#64-loveshow)
  * [6.5 DistrictJson](#65-districtjson)
  * [6.6 VipLevel信息](#66-viplevel信息)
  * [6.7 用户黑名单](#67-用户黑名单)
  * [6.8 GoldPrice信息](#68-goldprice信息)
  * [6.9 GoldList信息](#69-goldlist信息)
  * [6.10 送礼物后的变化信息](#610-送礼物后的变化信息)
  * [6.11 礼物列表详情](#611-礼物列表详情)
  * [6.12 未读消息信息](#612-未读消息信息)
  * [6.13 魅力排行](#613-魅力排行)
  * [6.14 财富排行(花费的财富)](#614-财富排行花费的财富)
* [7.管理页面](#7管理页面)
  * [7.1 管理页面](#71-管理页面)
* [8.礼物系统](#8礼物系统)
  * [8.1 金币价格](#81-金币价格)
  * [8.2 购买金币](#82-购买金币)
  * [8.3 礼物列表](#83-礼物列表)
  * [8.4 送出礼物](#84-送出礼物)
  * [8.5 收到礼物详情](#85-收到礼物详情)
  * [8.6 送出礼物详情](#86-送出礼物详情)
  * [8.7 异性魅力排行榜](#87-异性魅力排行榜)
  * [8.8 富豪榜](#88-富豪榜)

<!-- toc stop -->



##1.说明
本文档对各接口进行说明，最新请参考源码：[https://github.com/gemail/herefriend](https://github.com/gemail/herefriend)

##2.账户接口
### 2.1 用户注册
- [x] 已实现

- API

	**`GET`** *`/User/Register`*

- 参数

    |参数名|参数类型|取值范围|可选|默认值|参数说明|
    |:-- |:-- |:-- |:-- |:-- |:-- |
    |age|int|非0|否|/|注册用户年龄|
    |gender|int|0 (女)<br>1 (男)|否|/|注册用户性别|
    |cid|string|/|否|/|手机客户端个推clientid|

- 返回值

    |数据名称|数据类型/范围|说明|
    |:-- |:-- |:-- |
    |HTTPCODE|200 (OK)<br>404 (FAILED)|HTTP Get 返回值|
    |HTTP内容|registerInfo, [参考6.2]|- registerInfo.Id 累计增长<br>- registerInfo.PassWord 为随机6为数字<br>- registerInfo.Member 中的地区为自动检测IP获取|

### 2.2 用户登录
- [x] 已实现

- API

	**`GET`** *`/User/Login`*

- 参数

    |参数名|参数类型|取值范围|可选|默认值|参数说明|
    |:-- |:-- |:-- |:-- |:-- |:-- |
    |id|int|非0|否|/|用户id|
    |password|int|非0|否|/|用户密码|
    |cid|string|/|否|/|手机客户端个推clientid|

- 返回值

    |数据名称|数据类型/范围|说明|
    |:-- |:-- |:-- |
    |HTTPCODE|200 (OK)<br>404 (FAILED)|HTTP Get 返回值|
	|HTTP内容|personInfo, [参考6.1]|用户信息|

### 2.3 用户资料修改
- [x] 已实现

- API

	**`GET`** *`/User/SetProfile`*

- 参数

    |参数名|参数类型|取值范围|可选|默认值|参数说明|
    |:-- |:-- |:-- |:-- |:-- |:-- |
    |id|int|非0|否|/|用户id|
    |password|int|非0|否|/|用户密码|
    |newpassword|int|非0|是|/|用户新密码|
    |其他属性|参考personInfo结构体|类型相关|是|/|非(id, password, newpassword)字段|

- 返回值

    |数据名称|数据类型/范围|说明|
    |:-- |:-- |:-- |
    |HTTPCODE|200 (OK)<br>404 (FAILED)|HTTP Get 返回值|
    |HTTP内容|personInfo, [参考6.1]|修改后的用户信息|

### 2.4 上传图片
- [x] 已实现

- API

	**`POST`** *`/User/PostImage`*

- 参数

    |参数名|参数类型|取值范围|可选|默认值|参数说明|
    |:-- |:-- |:-- |:-- |:-- |:-- |
    |id|int|非0|否|/|用户id|
    |password|int|非0|否|/|用户密码|
    |pictype|int|0 (头像)<br>1 (相册)|否|/|照片类型|

- FormData说明

	`Content-Disposition form-data; name="file"; filename="%s"`

- 返回值

    |数据名称|数据类型/范围|说明|
    |:-- |:-- |:-- |
    |HTTPCODE|200 (OK)<br>404 (FAILED)|HTTP Get 返回值|
    |HTTP内容|personInfo, [参考6.1]|修改后的用户信息|

### 2.5 删除图片
- [x] 已实现

- API

	**`GET`** *`/User/DelImage`*

- 参数

    |参数名|参数类型|取值范围|可选|默认值|参数说明|
    |:-- |:-- |:-- |:-- |:-- |:-- |
    |id|int|非0|否|/|用户id|
    |password|int|非0|否|/|用户密码|
    |filename|string|非空|否|/|要删除的图片的名字，也可以是完整的URL地址|

- 返回值

    |数据名称|数据类型/范围|说明|
    |:-- |:-- |:-- |
    |HTTPCODE|200 (OK)<br>404 (FAILED)|HTTP Get 返回值|
    |HTTP内容|personInfo, [参考6.1]|修改后的用户信息|

### 2.6 用户退出
- [x] 已实现

- API

	**`GET`** *`/User/Logout`*

- 参数

    |参数名|参数类型|取值范围|可选|默认值|参数说明|
    |:-- |:-- |:-- |:-- |:-- |:-- |
    |id|int|非0|否|/|用户id|
    |password|int|非0|否|/|用户密码|

- 返回值

    |数据名称|数据类型/范围|说明|
    |:-- |:-- |:-- |
    |HTTPCODE|200 (OK)<br>404 (FAILED)|HTTP Get 返回值|
    |HTTP内容|/|/|

### 2.7 用户购买服务
- [x] 已实现

- API

	**`GET`** *`/User/BuyVIP`*

- 参数

    |参数名|参数类型|取值范围|可选|默认值|参数说明|
    |:-- |:-- |:-- |:-- |:-- |:-- |
    |id|int|非0|否|/|用户id|
    |password|int|非0|否|/|用户密码|
    |level|int|1 (写信会员)<br>2 (钻石会员)<br>3 (至尊会员)|否|/|VIP级别|
    |days|int|30 (月)<br>90 (季度)<br>180 (半年)<br>360 (年)|否|/|VIP时间|

- 返回值

    |数据名称|数据类型/范围|说明|
    |:-- |:-- |:-- |
    |HTTPCODE|200 (OK)<br>404 (FAILED)|HTTP Get 返回值|
    |HTTP内容|personInfo, [参考6.1]|用户信息|

### 2.8 获取用户信息
- [x] 已实现

- API

	**`GET`** *`/User/GetPersonInfo`*

- 参数

    |参数名|参数类型|取值范围|可选|默认值|参数说明|
    |:-- |:-- |:-- |:-- |:-- |:-- |
    |id|int|非0|否|/|要获取的用户id|

- 返回值

    |数据名称|数据类型/范围|说明|
    |:-- |:-- |:-- |
    |HTTPCODE|200 (OK)<br>404 (FAILED)|HTTP Get 返回值|
    |HTTP内容|personInfo, [参考6.1]||

### 2.9 用户在线狗叫握手
- [x] 已实现

- API

	**`GET`** *`/User/WatchDog`*

- 参数

    |参数名|参数类型|取值范围|可选|默认值|参数说明|
    |:-- |:-- |:-- |:-- |:-- |:-- |
    |id|int|非0|否|/|用户id|
    |password|int|非0|否|/|用户密码|

- 返回值

    |数据名称|数据类型/范围|说明|
    |:-- |:-- |:-- |
    |HTTPCODE|200 (OK)<br>404 (FAILED)|HTTP Get 返回值|
    |HTTP内容|/|/|

##3.消息接口
### 3.1 发送消息(打招呼/聊天)
- [x] 已实现

- API

	**`GET`** *`/Action/Recommend`*

- 参数

    |参数名|参数类型|取值范围|可选|默认值|参数说明|
    |:-- |:-- |:-- |:-- |:-- |:-- |
    |id|int|非0|否|/|用户id|
    |password|int|非0|否|/|用户密码|
    |toid|int|非0|否|/|打招呼对象的id|
    |type|int|1 (打招呼)<br>2 (聊天)<br>3 (心动)<br>4 (索要)|否|/|打招呼和心动可以不附带消息|
    |msg|string|/|是|空字符串|附带消息|

- 返回值

    |数据名称|数据类型/范围|说明|
    |:-- |:-- |:-- |
    |HTTPCODE|200 (OK)<br>404 (FAILED)|HTTP Get 返回值|
    |HTTP内容|当HTTPCODE为200: messageInfo, [参考6.3]|HTTPCODE为404时，HTTP内容只是描述错误，无其他参考价值|

### 3.2 删除消息
- [x] 已实现

- API

	**`GET`** *`/Action/DelRecommend`*

- 参数

    |参数名|参数类型|取值范围|可选|默认值|参数说明|
    |:-- |:-- |:-- |:-- |:-- |:-- |
    |id|int|非0|否|/|用户id|
    |password|int|非0|否|/|用户密码|
    |talkid|int|非0|否|/|聊天对象id|
    |msgid|int|非0|否|/|有效的messageInfo.MsgId|

- 返回值

    |数据名称|数据类型/范围|说明|
    |:-- |:-- |:-- |
    |HTTPCODE|200 (OK)<br>404 (FAILED)|HTTP Get 返回值|
    |HTTP内容|(空)|/|

### 3.3 查看资料
- [x] 已实现

- API

	**`GET`** *`/User/Visit`*

- 参数

    |参数名|参数类型|取值范围|可选|默认值|参数说明|
    |:-- |:-- |:-- |:-- |:-- |:-- |
    |id|int|非0|否|/|用户id|
    |password|int|非0|否|/|用户密码|
    |toid|int|非0|否|/|要访问的用户id|

- 返回值

    |数据名称|数据类型/范围|说明|
    |:-- |:-- |:-- |
    |HTTPCODE|200 (OK)<br>404 (FAILED)|HTTP Get 返回值|
    |HTTP内容|personInfo, [参考6.1]|要访问的用户信息|

### 3.4 设置"查看资料"信息为已读
- [x] 已实现

- API

	**`GET`** *`/User/ReadVisit`*

- 参数

    |参数名|参数类型|取值范围|可选|默认值|参数说明|
    |:-- |:-- |:-- |:-- |:-- |:-- |
    |id|int|非0|否|/|用户id|
    |password|int|非0|否|/|用户密码|
    |visitid|int|非0|否|/|有效的messageInfo.MsgId|

- 返回值

    |数据名称|数据类型/范围|说明|
    |:-- |:-- |:-- |
    |HTTPCODE|200 (OK)<br>404 (FAILED)|HTTP Get 返回值|
    |HTTP内容|/|/|

### 3.5 删除"查看资料"信息
- [x] 已实现

- API

	**`GET`** *`/User/DelVisit`*

- 参数

    |参数名|参数类型|取值范围|可选|默认值|参数说明|
    |:-- |:-- |:-- |:-- |:-- |:-- |
    |id|int|非0|否|/|用户id|
    |password|int|非0|否|/|用户密码|
    |visitid|int|非0|否|/|有效的visitInfo.VisitId|

- 返回值

    |数据名称|数据类型/范围|说明|
    |:-- |:-- |:-- |
    |HTTPCODE|200 (OK)<br>404 (FAILED)|HTTP Get 返回值|
    |HTTP内容|/|/|

### 3.6 获取与指定用户的聊天消息
- [x] 已实现

- API

	**`GET`** *`/User/WaterFlow`*

- 参数

    |参数名|参数类型|取值范围|可选|默认值|参数说明|
    |:-- |:-- |:-- |:-- |:-- |:-- |
    |id|int|非0|否|/|用户id|
    |password|int|非0|否|/|用户密码|
    |talkid|int|非0|否|/|聊天对象id|
    |lastmsgid|int|非0|是|0|服务端将返回lastmsgid之后的消息|
    |page|int|非0|是|1|页数|
    |count|int|非0|是|10|每页信息数目|

- 返回值

    |数据名称|数据类型/范围|说明|
    |:-- |:-- |:-- |
    |HTTPCODE|200 (OK)<br>404 (FAILED)|HTTP Get 返回值|
    |HTTP内容|messageInfo数组, [参考6.3]|/|

### 3.7 获取聊天记录和访问记录
- [x] 已实现

- API

	**`GET`** *`/User/AllMessage`*

- 参数

    |参数名|参数类型|取值范围|可选|默认值|参数说明|
    |:-- |:-- |:-- |:-- |:-- |:-- |
    |id|int|非0|否|/|用户id|
    |password|int|非0|否|/|用户密码|
    |lasttime|string|格林尼治时间|是|0代表的格林尼治时间|用于指定某个日期之后的信息,类似格式：<br>2016-01-01T06:03:32Z|
    |page|int|非0|是|1|页数|
    |count|int|非0|是|10|每页信息数目|

- 返回值

    |数据名称|数据类型/范围|说明|
    |:-- |:-- |:-- |
    |HTTPCODE|200 (OK)<br>404 (FAILED)|HTTP Get 返回值|
    |HTTP内容|allMessageInfo, [参考6.3]|/|

### 3.8 举报
- [x] 已实现

- API

	**`GET`** *`/User/Report`*

- 参数

    |参数名|参数类型|取值范围|可选|默认值|参数说明|
    |:-- |:-- |:-- |:-- |:-- |:-- |
    |id|int|非0|否|/|用户id|
    |password|int|非0|否|/|用户密码|
    |reportedid|int|非0|否|/|被举报的用户|
    |reason|string|/|是|空|举报原因|

- 返回值

    |数据名称|数据类型/范围|说明|
    |:-- |:-- |:-- |
    |HTTPCODE|200 (OK)<br>404 (FAILED)|HTTP Get 返回值|
    |HTTP内容|/|/|

### 3.9 添加用户黑名单
- [x] 已实现

- API

	**`GET`** *`/User/AddBlacklist`*

- 参数

    |参数名|参数类型|取值范围|可选|默认值|参数说明|
    |:-- |:-- |:-- |:-- |:-- |:-- |
    |id|int|非0|否|/|用户id|
    |password|int|非0|否|/|用户密码|
    |blacklistid|int|非0|否|/|被拉黑的用户|

- 返回值

    |数据名称|数据类型/范围|说明|
    |:-- |:-- |:-- |
    |HTTPCODE|200 (OK)<br>404 (FAILED)|HTTP Get 返回值|
    |HTTP内容|/|/|

### 3.10 删除用户黑名单
- [x] 已实现

- API

	**`GET`** *`/User/DelBlacklist`*

- 参数

    |参数名|参数类型|取值范围|可选|默认值|参数说明|
    |:-- |:-- |:-- |:-- |:-- |:-- |
    |id|int|非0|否|/|用户id|
    |password|int|非0|否|/|用户密码|
    |blacklistid|int|非0|否|/|被拉黑的用户|

- 返回值

    |数据名称|数据类型/范围|说明|
    |:-- |:-- |:-- |
    |HTTPCODE|200 (OK)<br>404 (FAILED)|HTTP Get 返回值|
    |HTTP内容|/|/|

###3.11 查询用户黑名单
- [x] 已实现

- API

	**`GET`** *`/User/GetBlacklist`*

- 参数

    |参数名|参数类型|取值范围|可选|默认值|参数说明|
    |:-- |:-- |:-- |:-- |:-- |:-- |
    |id|int|非0|否|/|用户id|
    |password|int|非0|否|/|用户密码|

- 返回值

    |数据名称|数据类型/范围|说明|
    |:-- |:-- |:-- |
    |HTTPCODE|200 (OK)<br>404 (FAILED)|HTTP Get 返回值|
    |HTTP内容|[参考6.7]|/|

### 3.12 获取未读消息数量
- [x] 已实现

- API

	**`GET`** *`/User/UnreadMessage`*

- 参数

    |参数名|参数类型|取值范围|可选|默认值|参数说明|
    |:-- |:-- |:-- |:-- |:-- |:-- |
    |id|int|非0|否|/|用户id|
    |password|int|非0|否|/|用户密码|
    |lasttime|string|格林尼治时间|是|0代表的格林尼治时间|用于指定某个日期之后的信息,类似格式：<br>2016-01-01T06:03:32Z|

- 返回值

    |数据名称|数据类型/范围|说明|
    |:-- |:-- |:-- |
    |HTTPCODE|200 (OK)<br>404 (FAILED)|HTTP Get 返回值|
    |HTTP内容|unreadMessageInfo, [参考6.12]|/|

### 3.13 获取消息列表
- [x] 已实现

- API

	**`GET`** *`/User/GetComments`*

- 参数

    |参数名|参数类型|取值范围|可选|默认值|参数说明|
    |:-- |:-- |:-- |:-- |:-- |:-- |
    |id|int|非0|否|/|用户id|
    |password|int|非0|否|/|用户密码|
    |lasttime|string|格林尼治时间|是|0代表的格林尼治时间|用于指定某个日期之后的信息,类似格式：<br>2016-01-01T06:03:32Z|

- 返回值

    |数据名称|数据类型/范围|说明|
    |:-- |:-- |:-- |
    |HTTPCODE|200 (OK)<br>404 (FAILED)|HTTP Get 返回值|
    |HTTP内容|messageInfo, [参考6.2]|/|

### 3.14 获取来访记录列表
- [x] 已实现

- API

	**`GET`** *`/User/GetVisits`*

- 参数

    |参数名|参数类型|取值范围|可选|默认值|参数说明|
    |:-- |:-- |:-- |:-- |:-- |:-- |
    |id|int|非0|否|/|用户id|
    |password|int|非0|否|/|用户密码|
    |lasttime|string|格林尼治时间|是|0代表的格林尼治时间|用于指定某个日期之后的信息,类似格式：<br>2016-01-01T06:03:32Z|

- 返回值

    |数据名称|数据类型/范围|说明|
    |:-- |:-- |:-- |
    |HTTPCODE|200 (OK)<br>404 (FAILED)|HTTP Get 返回值|
    |HTTP内容|messageInfo, [参考6.2]|/|

##4.功能接口
###4.1 搜索用户
- [x] 已实现

- API

	**`GET`** *`/User/Search`*

- 参数

    |参数名|参数类型|取值范围|可选|默认值|参数说明|
    |:-- |:-- |:-- |:-- |:-- |:-- |
    |id|int|非0|否|/|用户id|
    |password|int|非0|否|/|用户密码|
    |count|int|非0|是|10|每页信息数目||
    |Status|string|"在线"、"不限"|是|所有状态|考虑无实际效果，去掉"离线"搜索|
    |AgeMax|int|非0|是|/|最大年龄(包括)|
    |AgeMin|int|非0|是|/|最小年龄(包括)|
    |HeightMax|int|非0|是|/|最高身高(包括)|
    |HeightMin|int|非0|是|/|最低身高(包括)|
    |IncomeMax|int|非0|是|/|最高收入(包括)|
    |IncomeMin|int|非0|是|/|最低收入(包括)|
    |Study|string|非0|是|/|教育程度(完全匹配)|
    |Work|string|非0|是|/|工作(完全匹配)|
    |Province|string|非0|是|/|所在省(完全匹配, 字段值域可以通过/User/GetDistrict获得)|

- 返回值

    |数据名称|数据类型/范围|说明|
    |:-- |:-- |:-- |
    |HTTPCODE|200 (OK)<br>404 (FAILED)|HTTP Get 返回值|
    |HTTP内容|personInfo数组, [参考6.1]|/|

###4.2 心动列表
- [x] 已实现

- API

	**`GET`** *`/User/Heartbeat`*

- 参数

    |参数名|参数类型|取值范围|可选|默认值|参数说明|
    |:-- |:-- |:-- |:-- |:-- |:-- |
    |id|int|非0|否|/|用户id|
    |count|int|非0|是|10|每页信息数目|

- 返回值

    |数据名称|数据类型/范围|说明|
    |:-- |:-- |:-- |
    |HTTPCODE|200 (OK)<br>404 (FAILED)|HTTP Get 返回值|
    |HTTP内容|personInfo数组, [参考6.1]|/|

###4.3 恋爱秀列表
- [x] 已实现

- API

	**`GET`** *`/User/LoveShow`*

- 参数

    |参数名|参数类型|取值范围|可选|默认值|参数说明|
    |:-- |:-- |:-- |:-- |:-- |:-- |
    |id|int|非0|否|/|用户id|
    |password|int|非0|否|/|用户密码|
    |count|int|非0|是|10|每页信息数目|

- 返回值

    |数据名称|数据类型/范围|说明|
    |:-- |:-- |:-- |
    |HTTPCODE|200 (OK)<br>404 (FAILED)|HTTP Get 返回值|
    |HTTP内容|loveShow数组, [参考6.4]|/|

###4.4 恋爱秀送祝福
- [x] 已实现

- API

	**`GET`** *`/Action/LoveShowComment`*

- 参数

    |参数名|参数类型|取值范围|可选|默认值|参数说明|
    |:-- |:-- |:-- |:-- |:-- |:-- |
    |id|int|非0|否|/|用户id|
    |password|int|非0|否|/|用户密码|
    |loveshowid|int|非0|否|/|恋爱秀id|
    |bless|string|/|是|/|恋爱祝福|

- 返回值

    |数据名称|数据类型/范围|说明|
    |:-- |:-- |:-- |
    |HTTPCODE|200 (OK)<br>404 (FAILED)|HTTP Get 返回值|
    |HTTP内容|comment数组, [参考6.4]|/|

###4.5 地区列表
- [x] 已实现

- API

	**`GET`** *`/User/GetDistrict`*

- 参数

    |参数名|参数类型|取值范围|可选|默认值|参数说明|
    |:-- |:-- |:-- |:-- |:-- |:-- |
    |(无参数)|/|/|/|/|/|

- 返回值

    |数据名称|数据类型/范围|说明|
    |:-- |:-- |:-- |
    |HTTPCODE|200 (OK)<br>404 (FAILED)|HTTP Get 返回值|
    |HTTP内容|DistrictJson数组, [参考6.5]|/|

###4.6 VIP价格表
- [x] 已实现

- API

	**`GET`** *`/vip/price`*

- 参数

    |参数名|参数类型|取值范围|可选|默认值|参数说明|
    |:-- |:-- |:-- |:-- |:-- |:-- |
    |(无参数)|/|/|/|/|/|

- 返回值

    |数据名称|数据类型/范围|说明|
    |:-- |:-- |:-- |
    |HTTPCODE|200 (OK)<br>404 (FAILED)|HTTP Get 返回值|
    |HTTP内容|[参考6.6]|/|

###4.7 检查新版本
- [x] 已实现

- API

	**`GET`** *`/Version/Check`*

- 参数

    |参数名|参数类型|取值范围|可选|默认值|参数说明|
    |:-- |:-- |:-- |:-- |:-- |:-- |
    |version|string|'YYYYMMDD'|否|/|客户端版本，用于和服务器版本进行对比|

- 返回值

    |数据名称|数据类型/范围|说明|
    |:-- |:-- |:-- |
    |HTTPCODE|200 (有新版本)<br>404 (无新版本)|HTTP Get 返回值|
    |HTTP内容|[参考6.7]|/|

##5.推送消息
###5.1 透传消息推送
- 透传消息部分为如下结构体：

    ```go
    type PushMsgEvaluation struct {
    	Enable      bool   //是否要弹出评价对话框
    	ShowMessage string //弹出对话框显示的信息
    }

    type PushMsgRecvGift struct {
    	SenderId    int    //赠送者ID
    	GiftId      int    //礼物ID
    	GiftNum     int    //礼物数量
    	GiftName    string //礼物名称
    	ShowMessage string //弹出对话框显示的信息
    }
    
    type PushMessageInfo struct {
	    /*
	     * 根据类型不同，消息实体的结构体不同，如下为具体对应关系:
	     * ------------------------------------------------------
	     * |   Type |				Value						|
	     * ------------------------------------------------------
	     * |    1   |          PushMsgUnread (NOT Avaliable)    |
	     * ------------------------------------------------------
	     * |    2   |          PushMsgEvaluation                |
	     * ------------------------------------------------------
	     * |    3   |          PushMsgRecvGift                  |
	     * ------------------------------------------------------
	     */
    	Type  int    //消息类型
    	Value string //消息实体, 可以解析为对应的数据结构
    }
    ```

##6.JSON结构体
###6.1 personInfo
1. **结构体说明**

	存放用户信息，无信息字段为空值。

2. **CODE**

    ```go
    /*
     * Infomation shows to the clients
     */
    type PersonInfo struct {
    	Id              int                `json:",omitempty"` //ID号
    	Height          int                `json:",omitempty"` //身高
    	Weight          int                `json:",omitempty"` //体重
    	Age             int                `json:",omitempty"` //年龄
    	Gender          int                `json:",omitempty"` //性别: 0(女) 1(男)
    	OnlineStatus    int                `json:",omitempty"` //在线状态
    	VipLevel        int                `json:",omitempty"` //Vip级别
    	VipExpireTime   time.Time          `json:",omitempty"` //会员到期时间
    	Name            string             `json:",omitempty"` //姓名
    	Province        string             `json:",omitempty"` //所在省/直辖市/自治区
    	District        string             `json:",omitempty"` //所在区域
    	Native          string             `json:",omitempty"` //家乡
    	LoveType        string             `json:",omitempty"` //恋爱类型
    	BodyType        string             `json:",omitempty"` //体型
    	BloodType       string             `json:",omitempty"` //体型
    	Animal          string             `json:",omitempty"` //属相
    	Constellation   string             `json:",omitempty"` //星座
    	Lang            string             `json:",omitempty"` //语言
    	Introduction    string             //自我介绍
    	Selfjudge       string             `json:",omitempty"` //自评
    	Education       string             `json:",omitempty"` //教育程度
    	Income          string             `json:",omitempty"` //收入情况
    	IncomeMin       int                `json:",omitempty"` //收入最低
    	IncomeMax       int                `json:",omitempty"` //收入最高
    	School          string             `json:",omitempty"` //毕业学校
    	Occupation      string             `json:",omitempty"` //职业
    	Housing         string             `json:",omitempty"` //购房情况
    	Carstatus       string             `json:",omitempty"` //购车情况
    	Speciality      string             `json:",omitempty"` //技能
    	Marriage        string             `json:",omitempty"` //婚姻状况
    	Companytype     string             `json:",omitempty"` //公司类型
    	Companyindustry string             `json:",omitempty"` //公司领域
    	Nationnality    string             `json:",omitempty"` //民族
    	Religion        string             `json:",omitempty"` //信仰
    	Charactor       string             `json:",omitempty"` //性格类型
    	Hobbies         string             `json:",omitempty"` //兴趣爱好
    	CityLove        int                `json:",omitempty"` //是否接受异地恋: 0(视情况而定) 1(接受) 2(不接受)
    	Naken           int                `json:",omitempty"` //是否接受婚前性行为: 0(视情况而定) 1(接受) 2(不接受)
    	Allow_age       string             `json:",omitempty"` //择偶条件:年龄
    	Allow_residence string             `json:",omitempty"` //择偶条件:居住地
    	Allow_height    string             `json:",omitempty"` //择偶条件:身高
    	Allow_marriage  string             `json:",omitempty"` //择偶条件:婚姻状况
    	Allow_education string             `json:",omitempty"` //择偶条件:教育程度
    	Allow_housing   string             `json:",omitempty"` //择偶条件:购房情况
    	Allow_income    string             `json:",omitempty"` //择偶条件:收入
    	Allow_kidstatus string             `json:",omitempty"` //择偶条件:子女情况
    	IconUrl         string             `json:",omitempty"` //头像url
    	Pics            []string           `json:",omitempty"` //照片列表
    	GoldBeans       int                `json:",omitempty"` //用户的金币数量
    	RecvGiftList    []GiftSendRecvInfo `json:",omitempty"` //收到的礼物列表
    	SendGiftList    []GiftSendRecvInfo `json:",omitempty"` //送出的礼物列表
    }
    ```

###6.2 registerInfo
1. **结构体说明**

	存放用户注册信息，无信息字段为空值。

2. **CODE**

    ```go
    type reviewInfo struct {
        FoceShowReviewAlert bool
        ShowReviewAlert     bool
        ReviewAlertMsg      string
        ReviewAlertCancel   string
        ReviewAlertGo       string
    }

    type registerInfo struct {
        Id              int        //用户ID
        PassWord        string     //用户密码
        ShowGift        bool       //是否显示礼物
        ClientVersion   int        //客户端版本
        Member          personInfo //用户详细信息
        ReviewAlertInfo reviewInfo
    }
    ```

###6.3 聊天信息与访问资料信息
1. **结构体说明**

	存放打招呼信息，无信息字段为空值。

2. **CODE**

    ```go
	type messageInfo struct {
		MsgId     int        //消息Id
		MsgText   string     //消息内容, 无内容时此字段会自动隐藏
		UserId    int        //用户Id
		UserInfo  personInfo //用户信息, 无内容此字段会自动隐藏
		Direction int        //消息方向, 0: UserId发送给我的消息, 1: 我发送给UserId的消息
		Readed    bool       //客户端是否显示为已读
		TimeUTC   time.Time  //标准时间,用来参考转换为本地时间
	}

	type allMessageInfo struct {
		RecommendArray []messageInfo    
		VisitArray     []messageInfo    
	}
    ```

###6.4 loveShow
1. **结构体说明**

	存放恋爱秀信息，无信息字段自动隐藏。

2. **CODE**

    ```go
    type comment struct {
        Id        int       //用户id
        Age       int       //年龄
        Name      string    //姓名
        District  string    //地区
        Education string    //教育程度
        Text      string    //评论内容
        TimeUTC   time.Time //评论时间(标准时间)
    }

    type loverInfo struct {
        Id       int    //用户id
        Age      int    //年龄
        Name     string //姓名
        Imgurl   string //小头像地址
        District string //地区
    }

    type loveShow struct {
        Id           int       //爱情秀编号
        Girl         loverInfo //女生信息
        Guy          loverInfo //男生信息
        Daysfalllove int       //见面多少天之后就恋爱了
        Blessnum     int       //收到的祝福数量
        Lovestatus   string    //目前的恋爱状态
        Lovetitle    string    //爱情秀主题
        Lovestory    string    //爱情故事
        TimeUTC      time.Time //故事时间(标准时间)
        ShowPics     []string  //恋爱秀照片
        Comments     []comment //评论
    }
    ```

###6.5 DistrictJson
1. **结构体说明**

	存放全国地区信息，包括省和地区。

2. **CODE**

    ```go
    type DistrictJson struct {
        Province string
        District []string
    }
    ```

###6.6 VipLevel信息
1. **结构体说明**

	存放vip price信息。

2. **CODE**
    ```go
    type vipPriceInfo struct {
    	Days         int
    	Month        string
    	OldPrice     int
    	Price        int
    	Discount     string
    	DailyAverage string
    	ProductId    string
    	IabItemType  string
    }
    
    type vipInfo struct {
    	VipLevel     int
    	DiscountDesc string
    	PriceList    []vipPriceInfo
    }
    ```

###6.7 用户黑名单
1. **结构体说明**

	存放用户自己的黑名单信息。

2. **CODE**

    ```go
    type userBlacklist struct {
    	Id        int   //用户id
    	Blacklist []int //用户id的黑名单
    }
    ```

###6.8 GoldPrice信息
1. **结构体说明**

	存放金币信息。

2. **CODE**

    ```go
    type goldBeansPrice struct {
    	Price     int    //价格
    	Count     int    //普通会员购买数量
    	Song      int    //赠送
    	ProductId string //产品ID
    }
    ```

###6.9 GoldList信息
1. **结构体说明**

	存放礼物信息。

2. **CODE**

    ```go
    type giftInfo struct {
    	Id int //礼物固定id
    	/* 礼物类型：
    	 * 0 免费
    	 * 1 普通礼物
    	 * 2 折扣礼物
    	 * 3 名人礼物
    	 * ...
    	 */
    	Type                int
    	Name                string //礼物名称
    	ValidNum            int    //库存数量
    	Description         string //礼物描述
    	ImageUrl            string //礼物图片URL
    	Effect              int    //礼物特效，需要客户端支持
    	Price               int    //价格(beans)
    	OriginPrice         int    //原价(beans)，对于折扣礼物和Price不同
    	DiscountDescription string //折扣描述信息，对折扣作说明
    }
    ```

###6.10 送礼物后的变化信息
1. **结构体说明**

	送出礼物后，送礼物和收礼物的人的个人信息发生变化，APP需要更新这两个用户的个人信息。

2. **CODE**

    ```go
    type presentGiftInfo struct {
    	UserInfo    personInfo //个人信息
    	WhoRecvGift personInfo //收到礼物的人的信息
    }
    ```

###6.11 礼物列表详情
1. **结构体说明**

	对于收到或者送出的礼物详情，如果返回的数据代表收到的礼物详情，那么UserId代表送礼物的用户ID。如果返回的数据代表送出的礼物详情，那么UserId代表收到礼物的用户ID。

2. **CODE**

    ```go
    /*
     * 礼物列表详情
     */
    type giftListVerbose struct {
    	Person  personInfo //赠送礼物或者收到礼物的用户信息
    	GiftId  int        //礼物ID
    	GiftNum int        //礼物数量
    	Message string     //礼物留言
    	TimeUTC time.Time  //送礼物的时间
    }
    ```

###6.12 未读消息信息
1. **结构体说明**

	未读消息数量信息，用户替代未读消息推送消息

2. **CODE**

    ```go
    type unreadMessageInfo struct {
    	UnreadRecommend int //未读的聊天消息
    	UnreadVisit     int //未读的访问消息
    	Badge           int //badge: the icon number of app
    }
    ```

###6.13 魅力排行
1. **结构体说明**

	用户魅力信息

2. **CODE**

    ```go
    type userCharmInfo struct {
    	Person      personInfo //用户信息
    	GiftValue   int        //收到礼物的总价值
    	AdmireCount int        //被心仪的数量,暂无统计
    }
    ```

###6.14 财富排行(花费的财富)
1. **结构体说明**

	用户花费的财富排行

2. **CODE**

    ```go
    type userWealthInfo struct {
    	Person        common.PersonInfo //用户信息
    	ConsumedBeans int               //花费金币的总数量
    }
    ```

##7.管理页面
###7.1 管理页面
- [x] 已实现

- URL

	*`/login`*

- 说明

	登陆后进入CMS(content manage system, 内容管理系统)

##8.礼物系统
###8.1 金币价格
- [x] 已实现

- API

	**`GET`** *`/Gift/GoldPrice`*

- 参数

    |参数名|参数类型|取值范围|可选|默认值|参数说明|
    |:-- |:-- |:-- |:-- |:-- |:-- |
    |(无参数)|/|/|/|/|/|

- 返回值

    |数据名称|数据类型/范围|说明|
    |:-- |:-- |:-- |
    |HTTPCODE|200 (OK)<br>404 (FAILED)|HTTP Get 返回值|
    |HTTP内容|goldBeansPrice数组[参考6.8]|/|

###8.2 购买金币
- [x] 已实现

- API

	**`GET`** *`/Gift/BuyBeans`*

- 参数

    |参数名|参数类型|取值范围|可选|默认值|参数说明|
    |:-- |:-- |:-- |:-- |:-- |:-- |
    |id|int|非0|否|/|用户id|
    |password|int|非0|否|/|用户密码|
    |beans|int|非0|否|/|购买的金币数量|

- 返回值

    |数据名称|数据类型/范围|说明|
    |:-- |:-- |:-- |
    |HTTPCODE|200 (OK)<br>404 (FAILED)|HTTP Get 返回值|
    |HTTP内容|personInfo, [参考6.1]|/|

###8.3 礼物列表
- [x] 已实现

- API

	**`GET`** *`/Gift/GiftList`*

- 参数

    |参数名|参数类型|取值范围|可选|默认值|参数说明|
    |:-- |:-- |:-- |:-- |:-- |:-- |
    |(无参数)|/|/|/|/|/|

- 返回值

    |数据名称|数据类型/范围|说明|
    |:-- |:-- |:-- |
    |HTTPCODE|200 (OK)<br>404 (FAILED)|HTTP Get 返回值|
    |HTTP内容|giftInfo数组[参考6.9]|/|

###8.4 送出礼物
- [x] 已实现

- API

	**`GET`** *`/Gift/PresentGift`*

- 参数

    |参数名|参数类型|取值范围|可选|默认值|参数说明|
    |:-- |:-- |:-- |:-- |:-- |:-- |
    |id|int|非0|否|/|用户id|
    |password|int|非0|否|/|用户密码|
    |toid|int|非0|否|/|被赠送用户id|
    |giftid|int|非0|否|/|赠送的礼物id|
    |num|int|非0|否|/|赠送礼物的数量|
    |message|string|/|是|/|留言|

- 返回值

    |数据名称|数据类型/范围|说明|
    |:-- |:-- |:-- |
    |HTTPCODE|200 (OK)<br>404 (FAILED)|HTTP Get 返回值|
    |HTTP内容|presentGiftInfo, [参考6.10]|/|

###8.5 收到礼物详情
- [x] 已实现

- API

	**`GET`** *`/Gift/RecvListVerbose`*

- 参数

    |参数名|参数类型|取值范围|可选|默认值|参数说明|
    |:-- |:-- |:-- |:-- |:-- |:-- |
    |id|int|非0|否|/|用户id|
    |password|int|非0|否|/|用户密码|
    |queryid|int|非0|否|/|要查询的id|

- 返回值

    |数据名称|数据类型/范围|说明|
    |:-- |:-- |:-- |
    |HTTPCODE|200 (OK)<br>404 (FAILED)|HTTP Get 返回值|
    |HTTP内容|giftListVerbose数组, [参考6.11]|/|

###8.6 送出礼物详情
- [x] 已实现

- API

	**`GET`** *`/Gift/SendListVerbose`*

- 参数

    |参数名|参数类型|取值范围|可选|默认值|参数说明|
    |:-- |:-- |:-- |:-- |:-- |:-- |
    |id|int|非0|否|/|用户id|
    |password|int|非0|否|/|用户密码|
    |queryid|int|非0|否|/|要查询的id|

- 返回值

    |数据名称|数据类型/范围|说明|
    |:-- |:-- |:-- |
    |HTTPCODE|200 (OK)<br>404 (FAILED)|HTTP Get 返回值|
    |HTTP内容|giftListVerbose数组, [参考6.11]|/|

###8.7 异性魅力排行榜
- [x] 已实现

- API

	**`GET`** *`/User/CharmTopList`*

- 参数

    |参数名|参数类型|取值范围|可选|默认值|参数说明|
    |:-- |:-- |:-- |:-- |:-- |:-- |
    |id|int|非0|否|/|用户id|
    |password|int|非0|否|/|用户密码|
    |gender|int|0 (女)<br>1 (男)|否|/|用户性别|
    |page|int|非0|是|1|页数|
    |count|int|非0|是|10|每页信息数目|

- 返回值

    |数据名称|数据类型/范围|说明|
    |:-- |:-- |:-- |
    |HTTPCODE|200 (OK)<br>404 (FAILED)|HTTP Get 返回值|
    |HTTP内容|userCharmInfo数组, [参考6.13]|/|

###8.8 富豪榜
- [x] 已实现

- API

	**`GET`** *`/User/WealthList`*

- 参数

    |参数名|参数类型|取值范围|可选|默认值|参数说明|
    |:-- |:-- |:-- |:-- |:-- |:-- |
    |id|int|非0|否|/|用户id|
    |password|int|非0|否|/|用户密码|
    |page|int|非0|是|1|页数|
    |count|int|非0|是|10|每页信息数目|

- 返回值

    |数据名称|数据类型/范围|说明|
    |:-- |:-- |:-- |
    |HTTPCODE|200 (OK)<br>404 (FAILED)|HTTP Get 返回值|
    |HTTP内容|userWealthInfo数组, [参考6.14]|/|
