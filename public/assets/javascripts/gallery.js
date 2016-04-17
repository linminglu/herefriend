var g_gender = 0
var g_page = 1
var g_count = 100
var g_maxpage = 0
var g_maxusernum = [0, 0]

function setResultNull() {
	$("#gallery_list").html('<div class="text-center control-group error"><label class="control-label">无数据</label></div>')
	$("#gallery_list").hide()
	$("#gallery_list").fadeIn("slow");
	$("#gallery_list").fadeOut("slow");
	$("#gallery_list").fadeIn("slow");
}

function getUserInfo_ul(item) {
	var str = "<div class='picture'>"
	str += "<div class='tags' onclick='doHeartbeat(" + item["Id"] + ")'>"
    str += "<div class='badge badge-info pull-right'>" + item["Name"] + "-" + item["Age"] + "</div><br>"
    str += "<div class='badge pull-right' style='opacity:0.5'>" + item["Province"] + "</div><br>"
    str += ((true == item["Selected"]) ? "<div class='badge badge-important pull-right'><a class='icon-heart' style='color:white'/></div>": "") + "</div>"
	str += "<div class='actions'>"
	str += "<div class='pull-left'>"
    str += "<a class='btn btn-link' onclick='doChangeHeatPic(" + item["Id"] + ")'>"
    str += "    <small><i class='icon-retweet'></i></small></a>"
    str += "<a class='btn btn-link' onclick='doDeleteHeatPic(" + item["Id"] + ")'>"
    str += "    <i class='icon-cut'></i></a>"
    str += "<a class='btn btn-link' onclick='doAddBlackList(" + item["Id"] + ")'>"
    str += "    <i class='icon-remove-sign'></i></a></div>"
	str += "</div>"
	str += "<img class='" + (true == item["Selected"] ? "selected" : "") +" lazy' style='width:120px;height:120px' "
    str += "id='img_" + item["Id"] + "'"
    if ("" != item["Img"]) {
        str += " data-original='" + item["Img"] + "'"
    } else {
        str += " src='assets/images/black.jpg'"
    }
    str += " onclick='doHeartbeat(" + item["Id"] + ")'/>"
	str += "</div>"

	return str
}

function doHeartbeat(userid) {
	var isSelect = false

	if ($("img#img_" + userid).hasClass("selected")) {
		isSelect = true
	}

	$.get("/cms/SetHeartbeat?id=" + userid + "&action=" + (true == isSelect ? "0": "1") + "&gender=" + g_gender, "", function() {
		$.getJSON("/cms/GetSingleUserInfo?id=" + userid + "&gender=" + g_gender, function(item) {
			if (null != item) {
				$("li#user_" + userid).html(getUserInfo_ul(item))
				$("li#user_" + userid + " img").lazyload();
			}
		})
	}).fail(function(x) {
	    if (403 == x["status"]) {
		    alert("禁止操作!")
	    } else {
		    alert("发生错误,请检查网络!")
	    }
	})
}

function doChangeHeatPic(userid) {
	$.get("/cms/ChangeHeadImage?id=" + userid + "&gender=" + g_gender, "", function() {
		$.getJSON("/cms/GetSingleUserInfo?id=" + userid + "&gender=" + g_gender, function(item) {
			if (null != item) {
				$("li#user_" + userid).html(getUserInfo_ul(item))
				$("li#user_" + userid + " img").lazyload();
			}
		})
	}).fail(function(x) {
	    if (403 == x["status"]) {
		    alert("禁止操作!")
	    } else {
		    alert("发生错误,请检查网络!")
	    }
	})
}

function doDeleteHeatPic(userid) {
	$.get("/cms/DeleteHeadImage?id=" + userid + "&gender=" + g_gender, "", function() {
		$.getJSON("/cms/GetSingleUserInfo?id=" + userid + "&gender=" + g_gender, function(item) {
			if (null != item) {
				if ("" == item["Img"]) {
                    doAddBlackList(userid)
				} else {
					$("li#user_" + userid).html(getUserInfo_ul(item))
					$("li#user_" + userid + " img").lazyload();
				}
			}
		})
	}).fail(function(x) {
	    if (403 == x["status"]) {
		    alert("禁止操作!")
	    } else {
		    alert("发生错误,请检查网络!")
	    }
	})
}

function doAddBlackList(userid) {
	$.get("/cms/AddBlacklist?id=" + userid + "&gender=" + g_gender, "", function() {
		g_maxusernum[g_gender] = g_maxusernum[g_gender] - 1
		g_maxpage = parseInt((g_maxusernum[g_gender] + (g_count - 1)) / g_count)

		if (0 == listSelectByParam(g_gender, g_page, g_count, false, false)) {
			if (1 < g_page) {
				listSelectByParam(g_gender, g_page - 1, g_count, false, false)
			} else {
				setResultNull()
			}
		}
	}).fail(function(x) {
	    if (403 == x["status"]) {
		    alert("禁止操作!")
	    } else {
		    alert("发生错误,请检查网络!")
	    }
	})
}

function listSelectByParam(gender, page, count, bscroll, beffect) {
	if (page <= 0 || page > g_maxpage) {
		return 0
	}

	var result = 0
	$.ajax({
		type: "GET",
		url: "/cms/GetUserInfos?page=" + page + "&count=" + count + "&gender=" + gender,
		dataType: "json",
		async: false,
		cache: false,
		success: function(data) {
			if ((null == data) || (0 == data.length)) {
				return 0
			}

			var listStr = ""
			$.each(data, function(i, item) {
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

			g_page = page
			g_count = count
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

function refreshPageBtn(page) {
	var btnlist = ''
	var start = 0
	var end = 0
	var i = 0

	if (0 == g_maxpage) {
		start = end = 1
	} else {
		if (g_maxpage <= 5) {
			start = 1
			end = g_maxpage
		} else {
			if (page <= 3) {
				start = 1
				end = 5
			} else if (page >= (g_maxpage - 2)) {
				start = g_maxpage - 4
				end = g_maxpage
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

	if (g_maxpage == page) {
		btnlist += '<li class="disabled"><a class="icon-chevron-right"></a></li><li>'
		btnlist += '<li class="disabled"><a class="icon-step-forward"></a></li><li>'
	} else {
		btnlist += '<li><a class="icon-chevron-right" style="cursor:pointer" onclick="gotopage(' + String(page + 1) + ')"></a></li><li>'
		btnlist += '<li><a class="icon-step-forward" style="cursor:pointer" onclick="gotopage(' + String(g_maxpage) + ')"></a></li><li>'
	}

	$("#gallery_pagebtn").html(btnlist)
}

function home() {
	g_page = 1
	listSelectByParam(g_gender, g_page, g_count, true, true)
}

function gotopage(num) {
	if (num) {
		listSelectByParam(g_gender, num, g_count, true, true)
	}
}

function jumppage() {
	var num = parseInt($("#gallery_pagejump").val())
	if (num) {
		listSelectByParam(g_gender, num, g_count, true, true)
	}
}

function changegender() {
	g_gender = 1 - g_gender
	g_page = 1
	g_maxpage = parseInt((g_maxusernum[g_gender] + (g_count - 1)) / g_count)

	var result = listSelectByParam(g_gender, g_page, g_count, true, true)
	if ( - 1 == result) {
		g_gender = 1 - g_gender
		g_maxpage = parseInt((g_maxusernum[g_gender] + (g_count - 1)) / g_count)
	} else {
		if (0 == result) {
			setResultNull()
		}

		$("#gallery_changegender").html("切换为 [" + (0 == g_gender ? "男生": "女生") + "]")
	}
}

$(document).ready(function() {
	$('#gallery_pagejump').bind('keypress', function(event) {
		if (event.keyCode == "13") {
			jumppage()
		}
	});

	$.getJSON("/cms/sysuserinfo", function(info) {
		if (null != info) {
			g_maxusernum[0] = info["GirlsNum"]
			g_maxusernum[1] = info["GuysNum"]

			g_maxpage = parseInt((g_maxusernum[g_gender] + (g_count - 1)) / g_count)
			listSelectByParam(g_gender, g_page, g_count, true, true)
		} else {
			setResultNull()
		}
	}).fail(function(x) {
	    if (403 == x["status"]) {
		    alert("禁止操作!")
	    } else {
		    alert("发生错误,请检查网络!")
	    }
	})
});

