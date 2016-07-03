var gPage = 1
var gCount = 20
var gMaxPage = 0

var gTalkerCount = 10
var gTalkerPage = 0
var gTalkerMaxPage = 0
var focus_id = 0
var gCurUserId = 0
var gCurTalkerId = 0
var gLastMsgId = 0

var gShowprofile = true
var g_profile_id = 0

$.fn.scrollTo = function( target, options, callback ){
  if(typeof options == 'function' && arguments.length == 2){ callback = options; options = target; }
  var settings = $.extend({
    scrollTarget  : target,
    offsetTop     : 50,
    duration      : 500,
    easing        : 'swing'
  }, options);
  return this.each(function(){
    var scrollPane = $(this);
    var scrollTarget = (typeof settings.scrollTarget == "number") ? settings.scrollTarget : $(settings.scrollTarget);
    var scrollY = (typeof scrollTarget == "number") ? scrollTarget : scrollTarget.offset().top + scrollPane.scrollTop() - parseInt(settings.offsetTop);
    scrollPane.animate({scrollTop : scrollY }, parseInt(settings.duration), settings.easing, function(){
      if (typeof callback == 'function') { callback.call(this); }
    });
  });
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

function getUserInfo_ul(item) {
	var str = "<div class='picture'>"
	str += "<div class='tags' onclick='showprofile(" + item["Id"] + ")'>"
    str += "<div class='badge badge-info pull-right'>" + item["Name"] + "-" + item["Age"] + "</div><br>"
    str += "<div class='badge pull-right' style='opacity:0.5'>" + item["Province"] + "</div><br>"
    str += (0 != item["VipLevel"])? "<div class='badge badge-important pull-right'>vip:" + item["VipLevel"] + "</div>" : ""
    str += "</div>"
	str += "<img class='" + (true == item["Selected"] ? "selected" : "") +" lazy' style='width:120px;height:120px' "
    str += "id='img_" + item["Id"] + "'"
    if ("" != item["Img"]) {
        str += " data-original='" + item["Img"] + "'"
    } else {
        str += " src='assets/images/black.jpg'"
    }
    str += " onclick='showprofile(" + item["Id"] + ")'/>"
	str += "</div>"

	return str
}

function showprofile(userid) {
	if (false == gShowprofile) {
		animation()
	}

	if (0 != focus_id) {
		$("img#img_" + focus_id).removeClass("focus")
	}

	focus_id = userid
    refreshTalkWindows(1, userid)
	refreshProfileInfo(userid)
}

function animation() {
	if (false == gShowprofile) {
		gShowprofile = true
		$("#box-registlist").removeClass("span12").addClass("span5")
		$("#box-profile").toggle()
		$("#box-profile").animate({left:''}, "slow")
	} else {
		gShowprofile = false
		$("#box-profile").animate({
			left: '110%'
		},
		"slow", function() {
			$("#box-profile").toggle()
			$("#box-registlist").removeClass("span5").addClass("span12")
		})
	}
}

function refreshTalkWindows(talkerid, id) {
    gLastMsgId = 0
    $.getJSON("/cms/GetTalkHistory?id=" + id + "&talkid=" + talkerid + "&count=100&lastmsgid=" + gLastMsgId, function(data) {
        if (null != data) {
            $("#title-talkleft").html(data["UserName"])
            $("#title-talkright").html(data["TalkerName"])
		    var listStr = ""

            msgs = data["Comments"] || " "
            if (null != msgs && 0 != msgs.length) {
                var i = 0

                for (i = msgs.length - 1; i >= 0; i--) {
                    item = msgs[i]
                    if (item["MsgId"] > gLastMsgId) {
                        gLastMsgId = item["MsgId"]
                    }

                    t = new Date(item["TimeUTC"])
                    if (id == item["FromId"]) {
                        listStr += '<li class="left">' + item["MsgText"] + '<br><small class="date pull-left muted">' + t.formate("MM-dd HH:mm:ss") + '</small></li>'
                    } else {
                        listStr += '<li class="right">' + item["MsgText"] + '<br><small class="date pull-right muted">' + t.formate("MM-dd HH:mm:ss") + '</small></li>'
                    }
		        }
            }

		    $("#ul-talkwindow").html(listStr)

            if ("" == data["UserPic"]) {
                data["UserPic"] = "assets/images/black.jpg"
            }

            if ("" == data["TalkerPic"]) {
                data["TalkerPic"] = "assets/images/black.jpg"
            }

            $('#ul-talkwindow').append("<style>li.left:before{background-image:url('" + data["UserPic"] +"');background-size:contain}</style>");
            $('#ul-talkwindow').append("<style>li.right:before{background-image:url('" + data["TalkerPic"] +"');background-size:contain}</style>");
            lisize = $('#ul-talkwindow > li').length
            if (0 != lisize) {
                $('#ul-talkwindow'). scrollTo($('#ul-talkwindow > li')[lisize - 1])
            }

            gCurUserId = id
            gCurTalkerId = talkerid
        }
    }).fail(function() {
		alert("发生错误,请检查网络!")
	})
}

function refreshProfilePicture(userid) {
	$.getJSON("/cms/GetSingleUserInfo?id=" + userid, function(item) {
		if (null != item) {
			$("li#user_" + userid).html(getUserInfo_ul(item))
			$("li#user_" + userid + " img").lazyload();
		}
	})
}

function refreshProfileInfo(userid) {
	$.getJSON("/cms/GetSingleUserInfo?id=" + userid, function(item) {
		if (null != item) {
            if (null == item["VipSetAppVersion"]) {
			    $("#edit_appversion").editable("setValue", "")
            } else {
			    $("#edit_appversion").editable("setValue", item["VipSetAppVersion"])
            }

            if (null == item["VipLevel"]) {
			    $("#edit_viplevel").editable("setValue", 0)
            } else {
			    $("#edit_viplevel").editable("setValue", item["VipLevel"])
            }

			g_profile_id = userid
		}
	}).fail(function() {
		g_profile_id = 0
		alert("发生错误,请检查网络!")
	})
}

function listAdminChartsList(page, count, bscroll, beffect) {
	if (page <= 0 || (0 != gMaxPage && page > gMaxPage)) {
		return 0
	}

	var result = 0
	$.ajax({
		type: "GET",
		url: "/cms/AdminChartsList?page=" + page + "&count=" + count,
		dataType: "json",
		async: false,
		cache: false,
		success: function(data) {
			if ((null == data) || (0 == data.length)) {
				return 0
			}

			var listStr = ""
			$.each(data["Users"], function(i, item) {
				listStr += "<li id='user_" + item["Id"] + "'>" + getUserInfo_ul(item) + "</li>"
			});

			$("#gallery_list").html(listStr)

			if (true == bscroll) {
				scroll(0, 0)
			}

			if (true == beffect) {
				$("img.lazy").lazyload({
					effect: "fadeIn",
					failure_limit: 6,
					effect_speed: 1000
				});
			} else {
				$("img.lazy").lazyload();
			}

			gPage = page
	        gMaxPage = parseInt((data["Count"] + (gCount - 1)) / gCount)
			refreshPageBtn(page)

			result = 1
		},
		error: function() {
			alert("发生错误,请检查网络!")
			result = - 1
		}
	});

	return result
}

function refreshTalkWindow() {
    if (0 != gCurUserId && 0 != gCurTalkerId) {
        $.getJSON("/cms/GetTalkHistory?id=" + gCurUserId + "&talkid=" + gCurTalkerId + "&count=100&lastmsgid=" + gLastMsgId, function(data) {
            if (null != data) {
                $("#title-talkleft").html(data["UserName"])
                $("#title-talkright").html(data["TalkerName"])
                msgs = data["Comments"]
                if (null != msgs && 0 != msgs.length) {
	    	        var listStr = ""
                    var i = 0

                    for (i = msgs.length - 1; i >= 0; i--) {
                        item = msgs[i]
                        if (item["MsgId"] > gLastMsgId) {
                            gLastMsgId = item["MsgId"]
                        }

                        t = new Date(item["TimeUTC"])
                        if (gCurUserId == item["FromId"]) {
                            listStr += '<li class="left">' + item["MsgText"] + '<br><small class="date pull-left muted">' + t.formate("MM-dd HH:mm:ss") + '</small></li>'
                        } else {
                            listStr += '<li class="right">' + item["MsgText"] + '<br><small class="date pull-right muted">' + t.formate("MM-dd HH:mm:ss") + '</small></li>'
                        }
	    	        }

	    	        $("#ul-talkwindow").append(listStr)

                    lisize = $('#ul-talkwindow > li').length
                    if (0 != lisize) {
                        $('#ul-talkwindow'). scrollTo($('#ul-talkwindow > li')[lisize - 1])
                    }
                }
            }
        })
    }

	window.setTimeout(refreshTalkWindow, 10000)
}


function refreshPageBtn(page) {
	var btnlist = ''
	var start = 0
	var end = 0
	var i = 0

	if (0 == gMaxPage) {
		start = end = 1
	} else {
		if (gMaxPage <= 5) {
			start = 1
			end = gMaxPage
		} else {
			if (page <= 3) {
				start = 1
				end = 5
			} else if (page >= (gMaxPage - 2)) {
				start = gMaxPage - 4
				end = gMaxPage
			} else {
				start = page - 2
				end = page + 2
			}
		}
	}

	if (1 == page) {
		btnlist += '<li class="disabled"><a class="icon-step-backward"></a></li><li>'
		btnlist += '<li class="disabled"><a class="icon-chevron-left"></a></li><li>'
	} else {
		btnlist += '<li><a class="icon-step-backward" style="cursor:pointer" onclick="gotopage(1)"></a></li><li>'
		btnlist += '<li><a class="icon-chevron-left disabled" style="cursor:pointer" onclick="gotopage(' + String(page - 1) + ')"></a></li><li>'
	}

	for (i = start; i <= end; i++) {
		if (i == page) {
			btnlist += '<li class="active"><a>' + String(i) + '</a></li>'
		} else {
			btnlist += '<li><a style="cursor:pointer" onclick="gotopage(' + String(i) + ')">' + String(i) + '</a></li>'
		}
	}

	if (gMaxPage == page) {
		btnlist += '<li class="disabled"><a class="icon-chevron-right"></a></li><li>'
		btnlist += '<li class="disabled"><a class="icon-step-forward"></a></li><li>'
	} else {
		btnlist += '<li><a class="icon-chevron-right" style="cursor:pointer" onclick="gotopage(' + String(page + 1) + ')"></a></li><li>'
		btnlist += '<li><a class="icon-step-forward" style="cursor:pointer" onclick="gotopage(' + String(gMaxPage) + ')"></a></li><li>'
	}

	$("#gallery_pagebtn").html(btnlist)
}

function gotopage(num) {
	if (num) {
		listAdminChartsList(num, gCount, true, true)
	}
}

function chartspagejump() {
	gotopage(parseInt($("#charts_pagejump").val()))
}

function refreshEditable() {
	$('#edit_viplevel').editable({
		source: [{
			value: 0,
			text: "0"
		},
		{
			value: 1,
			text: "1"
		},
		{
			value: 2,
			text: "2"
		},
		{
			value: 3,
			text: "3"
		}],
		url: function(data) {
			result = $.ajax({
				type: "GET",
				url: "/cms/AdminGiveVipLevel?id=" + g_profile_id + "&level=" + data["value"],
				async: false,
				cache: false,
				success: function() {
					refreshProfilePicture(g_profile_id)
				},
			})

			return result
		}
	});

	$('#edit_appversion').editable({
		url: function(data) {
			result = $.ajax({
				type: "GET",
				url: "/cms/SetSingleUserInfo?id=" + g_profile_id + "&setvip_appversion=" + data["value"],
				async: false,
				cache: false,
			})

			return result
		}
	});
}

$(document).ready(function() {
    animation()
	refreshEditable()

	$('#charts_pagejump').bind('keypress', function(event) {
		if (event.keyCode == "13") {
			chartspagejump()
		}
	});

	listAdminChartsList(gPage, gCount, true, true)
    refreshTalkWindow()
});

