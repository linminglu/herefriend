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
  * [6.7 新版本信息](#67-新版本信息)
  * [6.8 用户黑名单](#68-用户黑名单)
* [7.管理页面](#7管理页面)
  * [7.1 管理页面](#71-管理页面)

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
    |type|int|1 (打招呼)<br>2 (聊天)<br>3 (心动)|否|/|打招呼和心动可以不附带消息|
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
    |HTTP内容|[参考6.8]|/|

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
    type PushMsgUnread struct {
    	UnreadRecommend int //未读的聊天消息
    	UnreadVisit     int //未读的访问消息
    	Badge           int //badge: the icon number of app
    }
    
    type PushMsgEvaluation struct {
    	Enable      bool   //是否要弹出评价对话框
    	ShowMessage string //弹出对话框显示的信息
    }
    
    type PushMessageInfo struct {
    	/*
    	 * 根据类型不同，消息实体的结构体不同，如下为具体对应关系:
    	 * 目前只有接收到Type=1的时候，APP应该修改Badge图标显示值
    	 * ------------------------------------------------------
    	 * | Type值 |         Value对应的数据结构               |
    	 * ------------------------------------------------------
    	 * |    1   |          PushMsgUnread                    |
    	 * ------------------------------------------------------
    	 * |    2   |          PushMsgEvaluation                |
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
    	Id              int       //ID号
    	Height          int       //身高
    	Weight          int       //体重
    	Age             int       //年龄
    	Gender          int       //性别: 0(女) 1(男)
    	OnlineStatus    int       //在线状态
    	VipLevel        int       //Vip级别
    	VipExpireTime   time.Time //会员到期时间
    	Name            string    //姓名
    	Province        string    //所在省/直辖市/自治区
    	District        string    //所在区域
    	Native          string    //家乡
    	LoveType        string    //恋爱类型
    	BodyType        string    //体型
    	BloodType       string    //体型
    	Animal          string    //属相
    	Constellation   string    //星座
    	Lang            string    //语言
    	Introduction    string    //自我介绍
    	Selfjudge       string    //自评
    	Education       string    //教育程度
    	Income          string    //收入情况
    	IncomeMin       int       //收入最低
    	IncomeMax       int       //收入最高
    	School          string    //毕业学校
    	Occupation      string    //职业
    	Housing         string    //购房情况
    	Carstatus       string    //购车情况
    	Speciality      string    //技能
    	Marriage        string    //婚姻状况
    	Companytype     string    //公司类型
    	Companyindustry string    //公司领域
    	Nationnality    string    //民族
    	Religion        string    //信仰
    	Charactor       string    //性格类型
    	Hobbies         string    //兴趣爱好
    	CityLove        int       //是否接受异地恋: 0(视情况而定) 1(接受) 2(不接受)
    	Naken           int       //是否接受婚前性行为: 0(视情况而定) 1(接受) 2(不接受)
    	Allow_age       string    //择偶条件:年龄
    	Allow_residence string    //择偶条件:居住地
    	Allow_height    string    //择偶条件:身高
    	Allow_marriage  string    //择偶条件:婚姻状况
    	Allow_education string    //择偶条件:教育程度
    	Allow_housing   string    //择偶条件:购房情况
    	Allow_income    string    //择偶条件:收入
    	Allow_kidstatus string    //择偶条件:子女情况
    	IconUrl         string    //头像url
    	Pics            []string  //照片列表
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

###6.7 新版本信息
1. **结构体说明**
	存放新版本信息。

2. **CODE**
    ```go
    type versionInfo struct {
    	Version string //版本
    	Url     string //channel渠道
    	Msg     string //版本说明
    	Force   bool
    }
    ```

###6.8 用户黑名单
1. **结构体说明**
	存放用户自己的黑名单信息。

2. **CODE**
    ```go
    type userBlacklist struct {
    	Id        int   //用户id
    	Blacklist []int //用户id的黑名单
    }
    ```

##7.管理页面
###7.1 管理页面
- [x] 已实现

- URL

	*`/login`*

- 说明

	登陆后进入CMS(content manage system, 内容管理系统)
