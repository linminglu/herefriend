var MAX_CPU_SIZE = 100
var MAX_MSG_SIZE = 20
var g_cpuarray = []
var g_cpuindex = 0
var g_msgarray = []
var g_msgindex = 0
var g_lastmsgid = 0
var g_msgcount = 0

function setcpuline(dot) {
	g_cpuarray[g_cpuindex] = parseFloat(dot)
	g_cpuindex = (g_cpuindex + 1) % MAX_CPU_SIZE

	var Sales = [];
	var index = g_cpuindex

	for (var i = 0; i < MAX_CPU_SIZE; i++) {
		Sales.push([i, g_cpuarray[(index + i) % MAX_CPU_SIZE]]);
	}

	$.plot($("#sys-cpu-chart"), [{
		data: Sales
	}], {
		series: {
			lines: {
				show: true,
				lineWidth: 3
			},
			shadowSize: 0
		},
		legend: {
			show: false
		},
		grid: {
			clickable: true,
			hoverable: true,
			borderWidth: 0,
			tickColor: "#f4f7f9"
		},
		colors: ["#00acec"],
		yaxis: {
			min: 0,
			max: 100
		},
		xaxis: {
			show: false
		}

	});
}

function refreshCpuinfo() {
	$.getJSON("/cms/cpuinfo", function(info) {
		if (null != info) {
			setcpuline(info["CpuUsage"])
		}
	})

	window.setTimeout(refreshCpuinfo, 2000)
}

function refreshSysinfo() {
	$.getJSON("/cms/sysinfo", function(info) {
		if (null != info) {
			$("#sys_describe").html("<br/>【OS】" + info["OSDescribe"] + "<br/>【CPU】" + info["CpuDescribe"])
			$("#sys_hdd").html(info["HDUsed"] + "G/" + info["HDTotal"] + "G (" + info["HDUsage"] + "%)")
			$("#sys_mem").html(info["MemUsed"] + "M/" + info["MemTotal"] + "M (" + info["MemUsage"] + "%)")
		}
	})

	window.setTimeout(refreshSysinfo, 10000)
}

function refreshUserinfo() {
	$.getJSON("/cms/sysuserinfo", function(info) {
		if (null != info) {
			$("#user_girlsnum").html(info["GirlsNum"])
			$("#user_guysnum").html(info["GuysNum"])
			$("#user_onlinenum").html(info["OnlineNum"] + " (15分钟在线:" + info["ActiveNum"] + ")")
			$("#user_registnum").html(info["RegistNum"])
		}
	})

	window.setTimeout(refreshUserinfo, 10000)
}

function refreshCommentinfo() {
	$.getJSON("/cms/commentinfo", function(info) {
		if (null != info) {
			$("#comment_talknum").html(info["TalkNum"])
			$("#comment_pushnum").html(info["PushNum"])
			$("#comment_buyvipnum").html(info["BuyVIPNum"])
		}
	})

	window.setTimeout(refreshCommentinfo, 10000)
}

function getmsgtypestring(type) {
	if (1 == type) {
		return "打招呼"
	} else if (2 == type) {
		return "聊天"
	} else if (3 == type) {
		return "心动"
	} else {
		return "未知消息"
	}
}

/**      
* 对Date的扩展，将 Date 转化为指定格式的String      
* 月(M)、日(d)、12小时(h)、24小时(H)、分(m)、秒(s)、周(E)、季度(q) 可以用 1-2 个占位符      
* 年(y)可以用 1-4 个占位符，毫秒(S)只能用 1 个占位符(是 1-3 位的数字)      
* eg:      
* (new Date()).formate("yyyy-MM-dd hh:mm:ss.S") ==> 2006-07-02 08:09:04.423      
* (new Date()).formate("yyyy-MM-dd E HH:mm:ss") ==> 2009-03-10 二 20:09:04      
* (new Date()).formate("yyyy-MM-dd EE hh:mm:ss") ==> 2009-03-10 周二 08:09:04      
* (new Date()).formate("yyyy-MM-dd EEE hh:mm:ss") ==> 2009-03-10 星期二 08:09:04      
* (new Date()).formate("yyyy-M-d h:m:s.S") ==> 2006-7-2 8:9:4.18      
*/
Date.prototype.formate = function(fmt) {
	var o = {
		"M+": this.getMonth() + 1,
		//月份         
		"d+": this.getDate(),
		//日         
		"h+": this.getHours() % 12 == 0 ? 12: this.getHours() % 12,
		//小时         
		"H+": this.getHours(),
		//小时         
		"m+": this.getMinutes(),
		//分         
		"s+": this.getSeconds(),
		//秒         
		"q+": Math.floor((this.getMonth() + 3) / 3),
		//季度         
		"S": this.getMilliseconds() //毫秒         
	};
	var week = {
		"0": "\u65e5",
		"1": "\u4e00",
		"2": "\u4e8c",
		"3": "\u4e09",
		"4": "\u56db",
		"5": "\u4e94",
		"6": "\u516d"
	};
	if (/(y+)/.test(fmt)) {
		fmt = fmt.replace(RegExp.$1, (this.getFullYear() + "").substr(4 - RegExp.$1.length));
	}
	if (/(E+)/.test(fmt)) {
		fmt = fmt.replace(RegExp.$1, ((RegExp.$1.length > 1) ? (RegExp.$1.length > 2 ? "\u661f\u671f": "\u5468") : "") + week[this.getDay() + ""]);
	}
	for (var k in o) {
		if (new RegExp("(" + k + ")").test(fmt)) {
			fmt = fmt.replace(RegExp.$1, (RegExp.$1.length == 1) ? (o[k]) : (("00" + o[k]).substr(("" + o[k]).length)));
		}
	}
	return fmt;
}

function getcomment_li(info) {
	//var timestr = info["TimeUTC"].replace('Z', '').replace('T', ' ')
	var t = new Date(info["TimeUTC"])
	var str = "<li><div class='avatar pull-left'>" +
	 "<img class='lazy' data-original='" + info["FromPic"] + "' height='23' width='23'/></div>" +
	 "<div class='body'><div class='name'>" + "<a class='text-contrast'>" + info["From"] + " </a><a class='icon-volume-up'>" +
	 "<a class='text-contrast'> " + info["To"] + "</a></div>" +
	 "<div class='text'><a class='text-contrast'>" + getmsgtypestring(info["MsgType"]) + " </a>" + info["MsgText"] + "</div></div>" +
	 "<div class='text-right'><small><a class='text-contrast'>" + t.formate("yyyy-MM-dd EE HH:mm:ss") + " </a>" + "<i class='icon-time'></i></small></div></li>"

	return str
}

function refreshCommendHistory() {
	$.getJSON("/cms/recommendhistory?lastmsgid=" + g_lastmsgid, function(data) {
		if (null != data && data.length) {
			$.each(data, function(i, info) {
				g_msgarray[g_msgindex] = info
				g_msgindex = (g_msgindex + 1) % MAX_MSG_SIZE
				g_lastmsgid = parseInt(info["MsgId"])

				if (g_msgcount < MAX_MSG_SIZE) {
					g_msgcount++
				}
			})

			var ulbodystr = ""
			for (var i = 0; i < g_msgcount; i++) {
				ulbodystr += getcomment_li(g_msgarray[(MAX_MSG_SIZE + g_msgindex - 1 - i) % MAX_MSG_SIZE])
			}

			$("#recent_comments").html(ulbodystr);
			$("#recent_comments img.lazy").lazyload();
		}
	})

	window.setTimeout(refreshCommendHistory, 10000)
}

$(document).ready(function() {
	refreshSysinfo()
	refreshCpuinfo()
	refreshUserinfo()
	refreshCommentinfo()
	refreshCommendHistory()
});

