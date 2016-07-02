var gGender = 0
var gPage = 1
var gCount = 20
var gMaxPage = 0

var gTalkerCount = 10
var gTalkerPage = 0
var gTalkerMaxPage = 0

var gShowWindow = true

var focus_id = 0
var gCurUserId = 0
var gCurTalkerId = 0
var gLastMsgId = 0

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

function setResultNull() {
	$("#gallery_list").html('<div class="text-center control-group error"><label class="control-label">无数据</label></div>')
	$("#gallery_list").hide()
	$("#gallery_list").fadeIn("slow");
	$("#gallery_list").fadeOut("slow");
	$("#gallery_list").fadeIn("slow");
}

function getUserInfo_ul(item) {
	var str = "<div class='picture'>"
	str += "<div class='tags' onclick='getTalkerList(" + item["Id"] + ")'>"
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
    str += " onclick='getTalkerList(" + item["Id"] + ")'/>"
	str += "</div>"

	return str
}

function animation() {
	if (false == gShowWindow) {
		gShowWindow = true
		$("#box-registlist").removeClass("span12").addClass("span8")
		$("#box-chartswindows").toggle()
		$("#box-chartswindows").animate({left:''}, "slow")
	} else {
		gShowWindow = false
		$("#box-chartswindows").animate({
			left: '110%'
		},
		"slow", function() {
			$("#box-chartswindows").toggle()
			$("#box-registlist").removeClass("span8").addClass("span12")
		})
	}
}

function getTalkerInfo_ul(id, item) {
    var talkerid = 0
    var talkerinfo = ""
    var picture = ""
	var t = new Date(item["TimeUTC"])

    if (0 == item["Direction"]) {
        talkerid = item["FromId"]
        talkerinfo = item["From"]
        picture = item["FromPic"]
    } else {
        talkerid = item["ToId"]
        talkerinfo = item["To"]
        picture = item["ToPic"]
    }

	var str = '<li><div id="talker_' + talkerid + '">'
            + '<div class="pull-left">'
	        + '    <img class="lazy" style="width:50px;height:50px;border-radius:50px;cursor:pointer;"'

    if ("" != picture) {
        str += ' data-original="' + picture + '"'
    } else {
        str += ' src="assets/images/black.jpg"'
    }

    //str += ' <a class="btn btn-mini right" data-toggle="modal" onclick="refreshGiftVerboseDlg(' + item["Id"] + ')" href="#gift_verbose_dlg" role="button">详情<span>&gt;</span></a>'
    str += ' data-toggle="modal" onclick="getTalkWindow(' + id + ',' + talkerid + ')" href="#talk_window_dlg"/>'
         + '</div>'
         + '<div class="text-right pull-right">'
         + '    <p>' + item["MsgText"] + '</p>'
         + '</div>'
         + '<div class="clearfix"></div>'
         + '<div class="pull-left">'
         + '    <p>'
         + '        <span class="text-contrast">' + talkerinfo + '</span>'
         + '    </p>'
         + '</div>'
         + '<div class="text-right pull-right">'
         + '    <p>'
         + '        <i class="icon-time text-muted"></i>'
         + '        <span class="text-muted">' + t.formate("yyyy-MM-dd HH:mm:ss") + '</span>'
         + '    </p>'
         + '</div>'
         + '<div class="clearfix"></div>'
	     + '</div></li>'

	return str
}

function getTalkWindow(id, talkerid) {
    gLastMsgId = 0
    $.getJSON("/cms/GetTalkHistory?id=" + id + "&talkid=" + talkerid + "&count=100&lastmsgid=" + gLastMsgId, function(data) {
        if (null != data) {
            $("#title-talkleft").html(data["UserName"])
            $("#title-talkright").html(data["TalkerName"])
		    var listStr = ""

            msgs = data["Comments"]
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

function dotalk() {
    if (0 == gCurUserId || 0 == gCurTalkerId) {
        return
    }

    msg = $("#talk_input").val().trim()
    if (0 == msg.length) {
        return
    }

	var result = 0
	$.ajax({
		type: "GET",
		url: "/cms/DoTalk?fromid=" + gCurTalkerId + "&toid=" + gCurUserId + "&msg=" + msg,
		dataType: "json",
		async: false,
		cache: false,
		success: function(data) {
            $("#talk_input").val("")
            if (data["MsgId"] > gLastMsgId) {
                gLastMsgId = data["MsgId"]
            }

            t = new Date(data["TimeUTC"])
            $('#ul-talkwindow').append('<li class="right">' + msg + '<br><small class="date pull-right muted">' + t.formate("MM-dd HH:mm:ss") + '</small></li>')

            lisize = $('#ul-talkwindow > li').length
            if (0 != lisize) {
                $('#ul-talkwindow'). scrollTo($('#ul-talkwindow > li')[lisize - 1])
            }
		},
		error: function() {
			alert("发生错误,请检查网络!")
		}
	});
}

function getTalkerList(id) {
	if (false == gShowWindow) {
		animation()
	}

    if (0 != focus_id) {
        $("img#img_"+focus_id).removeClass("focus")
    }

    focus_id = id
    $("img#img_"+id).addClass("focus")
    
    $.getJSON("/cms/GetChartsList?id=" + id + "&count=100", function(data) {
		var listStr = ""
        if (null != data) {
		    $.each(data, function(i, item) {
		    	listStr += getTalkerInfo_ul(id, item)
		    });
        }

		$("#talker_list").html(listStr)
		$("#talker_list img.lazy").lazyload({effect: "fadeIn", failure_limit: 6, effect_speed: 1000});
    }).fail(function() {
		alert("发生错误,请检查网络!")
	})
}

function listRegistUserInfo(gender, page, count, bscroll, beffect) {
	if (page <= 0 || (0 != gMaxPage && page > gMaxPage)) {
		return 0
	}

	var result = 0
	$.ajax({
		type: "GET",
		url: "/cms/RegistUserInfo?page=" + page + "&count=" + count + "&gender=" + gender,
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
		listRegistUserInfo(gGender, num, gCount, true, true)
	}
}

function changegender() {
	gGender = 1 - gGender
	gPage = 1

	var result = listRegistUserInfo(gGender, gPage, gCount, true, true)
	if ( - 1 == result) {
		gGender = 1 - gGender
	} else {
		if (0 == result) {
			setResultNull()
		}

		$("#gallery_changegender").html("切换为 [" + (0 == gGender ? "男生": "女生") + "]")
	}
}

function chartspagejump() {
	gotopage(parseInt($("#charts_pagejump").val()))
}

$(document).ready(function() {
    animation()

	$('#charts_pagejump').bind('keypress', function(event) {
		if (event.keyCode == "13") {
			chartspagejump()
		}
	});

	$('#talk_input').bind('keypress', function(event) {
		if (event.keyCode == "13") {
			dotalk()
		}
	});

	listRegistUserInfo(gGender, gPage, gCount, true, true)
    refreshTalkWindow()
});

