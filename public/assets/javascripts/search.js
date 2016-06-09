var g_page = 1
var g_count = 16
var g_maxpage = 0
var g_showprofile = true
var g_profile_id = 0

var search_gender = 0
var search_field = 0
var search_lastfield = 0
var search_key = ""
var search_lastkey = ""
var search_count = 0

var focus_id = 0

function setSearchResultNull() {
	$("#gallery_list").html('<div class="text-center control-group error"><label class="control-label">无数据</label></div>')
	$("#gallery_list").hide()
	$("#gallery_list").fadeIn("slow");
	$("#gallery_list").fadeOut("slow");
	$("#gallery_list").fadeIn("slow");
}

function animation() {
	if (false == g_showprofile) {
		g_showprofile = true
		$("#box-searchresult").removeClass("span12").addClass("span5")
		$("#box-profile").toggle()
		$("#box-profile").animate({
			left: ''
		},
		"slow")
	} else {
		g_showprofile = false
		$("#box-profile").animate({
			left: '110%'
		},
		"slow", function() {
			$("#box-profile").toggle()
			$("#box-searchresult").removeClass("span5").addClass("span12")
		})
	}
}

function getUserInfo_ul(item) {
	var str = "<div class='picture'>"
	str += "<div class='tags' onclick='showprofile(" + item["Id"] + ")'>"
	str += "<div class='badge badge-info pull-right'>" + item["Name"] + "-" + item["Age"] + "</div><br>"
	str += "<div class='badge pull-right' style='opacity:0.5'>" + item["Province"] + "</div><br>"
	str += ((true == item["Selected"]) ? "<div class='badge badge-important pull-right'><a class='icon-heart' style='color:white'/></div>": "") + "</div>"
	str += "<img class='" + (true == item["Selected"] ? "selected": "") + " lazy' style='width:120px;height:120px' "
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

function getUserProfilePicture(item) {
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
	str += "<img class='" + (true == item["Selected"] ? "selected": "") + " lazy' style='width:120px;height:120px' "
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

function refreshProfileInfo(userid) {
	$.getJSON("/User/GetPersonInfo?id=" + userid, function(item) {
		if (null != item) {
			$("#search_edit_userid").text(item["Id"])
			$("#search_edit_age").editable("setValue", item["Age"])
            if (null == item["Name"]) {
			    $("#search_edit_username").editable("setValue", "")
            } else {
			    $("#search_edit_username").editable("setValue", item["Name"])
            }
            if (null == item["Introduction"]) {
			    $("#search_edit_introduction").editable("setValue", "")
            } else {
			    $("#search_edit_introduction").editable("setValue", item["Introduction"])
            }
            if (null == item["VipLevel"]) {
			    $("#search_edit_viplevel").editable("setValue", 0)
            } else {
			    $("#search_edit_viplevel").editable("setValue", item["VipLevel"])
            }

			g_profile_id = userid
		}

		$.getJSON("/cms/GetSingleUserInfo?id=" + userid + "&gender=" + search_gender, function(item) {
			if (null != item) {
				$("#search_edit_delete").editable("setValue", (1 == item["Usertype"]) ? 0: 1)
			}
		})
	}).fail(function() {
		g_profile_id = 0
		alert("发生错误,请检查网络!")
	})
}

function refreshProfilePicture(userid) {
	$("#profile_picture").html("")
	$.getJSON("/cms/GetSingleUserInfo?id=" + userid + "&gender=" + search_gender, function(item) {
		if (null != item) {
			$("li#user_" + userid).html(getUserInfo_ul(item))
			$("li#user_" + userid + " img").lazyload();
			$("#profile_picture").html("")
			$("#profile_picture").html("<li id='profileuser_" + userid + "'>" + getUserProfilePicture(item) + "</li>")
			$("img#img_" + userid).addClass("focus")
			$("img#img_" + userid).lazyload({
				effect: "fadeIn",
				failure_limit: 6,
				effect_speed: 1000
			});
		}
	})
}

function showprofile(userid) {
	if (false == g_showprofile) {
		animation()
	}

	if (0 != focus_id) {
		$("img#img_" + focus_id).removeClass("focus")
	}

	focus_id = userid

	refreshProfilePicture(userid)
	refreshProfileInfo(userid)
}

function doHeartbeat(userid) {
	var isSelect = false

	if ($("img#img_" + userid).hasClass("selected")) {
		isSelect = true
	}

	$.get("/cms/SetHeartbeat?id=" + userid + "&action=" + (true == isSelect ? "0": "1") + "&gender=" + search_gender, "", function() {
		refreshProfilePicture(userid)
	}).fail(function() {
		alert("发生错误,请检查网络!")
	})
}

function doChangeHeatPic(userid) {
	$.get("/cms/ChangeHeadImage?id=" + userid + "&gender=" + search_gender, "", function() {
		refreshProfilePicture(userid)
	}).fail(function() {
		alert("发生错误,请检查网络!")
	})
}

function doDeleteHeatPic(userid) {
	$.get("/cms/DeleteHeadImage?id=" + userid + "&gender=" + search_gender, "", function() {
		$.getJSON("/cms/GetSingleUserInfo?id=" + userid + "&gender=" + search_gender, function(item) {
			if (null != item) {
				if ("" == item["Img"]) {
					doAddBlackList(userid)
				} else {
					$("li#user_" + userid).html(getUserInfo_ul(item))
					$("li#user_" + userid + " img").lazyload();
					$("#profile_picture").html("")
					$("#profile_picture").html("<li id='profileuser_" + userid + "'>" + getUserProfilePicture(item) + "</li>")
					$("img#img_" + userid).lazyload({
						effect: "fadeIn",
						failure_limit: 6,
						effect_speed: 1000
					});
				}
			}
		})
	}).fail(function() {
		alert("发生错误,请检查网络!")
	})
}

function doAddBlackList(userid) {
	$.get("/cms/AddBlacklist?id=" + userid + "&gender=" + search_gender, "", function() {
		if (true == g_showprofile) {
			animation()
		}

		if (0 == searchcallback(search_gender, g_page, search_lastfield, search_lastkey, false, false)) {
			if (1 < g_page) {
				searchcallback(search_gender, g_page - 1, search_lastfield, search_lastkey, true, true)
			} else {
				setSearchResultNull()
			}
		}
	}).fail(function() {
		alert("发生错误,请检查网络!")
	})
}

function listSelectByParam(data, page, bscroll, beffect) {
	if (null == data || page <= 0 || page > g_maxpage) {
		return
	}

	refreshPageBtn(page)

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

function gotopage(num) {
	if (0 == search_count) {
		return
	}

	if ((1 <= num) && (num <= g_maxpage)) {
		if (0 == searchcallback(search_gender, num, search_lastfield, search_lastkey, true, true)) {
			setSearchResultNull()
		}
	}
}

function searchpagejump() {
	gotopage(parseInt($("#search_pagejump").val()))
}

function searchcallback(gender, page, field, key, bscroll, beffect) {
	var result = 0

	$.ajax({
		type: "GET",
		url: "/cms/SearchUserInfos?gender=" + gender + "&page=" + page + "&count=" + g_count + "&field=" + field + "&key=" + key,
		dataType: "json",
		async: false,
		cache: false,
		success: function(data) {
			if (null != data) {
				search_count = data["Count"]
				if (0 == search_count) {
					return 0
				} else {
					search_lastfield = field
					search_lastkey = key
					g_page = page

					$("#gallery_list").html("")
					g_maxpage = parseInt((search_count + (g_count - 1)) / g_count)
					listSelectByParam(data["Users"], page, g_count, bscroll, beffect)

					result = 1
				}
			}
		},
		error: function() {
			alert("发生错误,请检查网络!")
			result = - 1
		}
	})

	return result
}

function dosearch(gender) {
	search_key = $("#search_input").val().trim()

	if (true == g_showprofile) {
		animation()
	}

	search_count = 0
	search_gender = gender
	if (0 == searchcallback(search_gender, 1, search_field, search_key, true, true)) {
		setSearchResultNull()
	}
}

function setsearchfield(field) {
	search_field = field
	$("#search_span").html(["姓名", "自我描述", "ID"][field])
}

function refreshEditable() {
	$('#search_edit_username').editable({
		url: function(data) {
			result = $.ajax({
				type: "GET",
				url: "/cms/SetSingleUserInfo?id=" + g_profile_id + "&gender=" + search_gender + "&name=" + data["value"],
				async: false,
				cache: false,
				success: function() {
					refreshProfilePicture(g_profile_id)
				}
			})

			return result
		}
	});

	$('#search_edit_delete').editable({
		source: [{
			value: 0,
			text: "否"
		},
		{
			value: 1,
			text: "是"
		}],
		url: function(data) {
			result = $.ajax({
				type: "GET",
				url: "/cms/SetSingleUserInfo?id=" + g_profile_id + "&gender=" + search_gender + "&delete=" + data["value"],
				async: false,
				cache: false,
			})
			return result
		}
	});

	$('#search_edit_age').editable({
		source: [{
			value: 18,
			text: '18'
		},
		{
			value: 19,
			text: '19'
		},
		{
			value: 20,
			text: '20'
		},
		{
			value: 21,
			text: '21'
		},
		{
			value: 22,
			text: '22'
		},
		{
			value: 23,
			text: '23'
		},
		{
			value: 24,
			text: '24'
		},
		{
			value: 25,
			text: '25'
		},
		{
			value: 26,
			text: '26'
		},
		{
			value: 27,
			text: '27'
		},
		{
			value: 28,
			text: '28'
		},
		{
			value: 29,
			text: '29'
		},
		{
			value: 30,
			text: '30'
		},
		{
			value: 31,
			text: '31'
		},
		{
			value: 32,
			text: '32'
		},
		{
			value: 33,
			text: '33'
		},
		{
			value: 34,
			text: '34'
		},
		{
			value: 35,
			text: '35'
		},
		{
			value: 36,
			text: '36'
		},
		{
			value: 37,
			text: '37'
		},
		{
			value: 38,
			text: '38'
		},
		{
			value: 39,
			text: '39'
		},
		{
			value: 40,
			text: '40'
		},
		{
			value: 41,
			text: '41'
		},
		{
			value: 42,
			text: '42'
		},
		{
			value: 43,
			text: '43'
		},
		{
			value: 44,
			text: '44'
		},
		{
			value: 45,
			text: '45'
		},
		{
			value: 46,
			text: '46'
		},
		{
			value: 47,
			text: '47'
		},
		{
			value: 48,
			text: '48'
		},
		{
			value: 49,
			text: '49'
		},
		{
			value: 50,
			text: '50'
		},
		{
			value: 51,
			text: '51'
		},
		{
			value: 52,
			text: '52'
		},
		{
			value: 53,
			text: '53'
		},
		{
			value: 54,
			text: '54'
		},
		{
			value: 55,
			text: '55'
		},
		{
			value: 56,
			text: '56'
		},
		{
			value: 57,
			text: '57'
		},
		{
			value: 58,
			text: '58'
		},
		{
			value: 59,
			text: '59'
		},
		{
			value: 60,
			text: '60'
		},
		{
			value: 61,
			text: '61'
		},
		{
			value: 62,
			text: '62'
		},
		{
			value: 63,
			text: '63'
		},
		{
			value: 64,
			text: '64'
		},
		{
			value: 65,
			text: '65'
		},
		{
			value: 66,
			text: '66'
		},
		{
			value: 67,
			text: '67'
		},
		{
			value: 68,
			text: '68'
		},
		{
			value: 69,
			text: '69'
		},
		{
			value: 70,
			text: '70'
		},
		{
			value: 71,
			text: '71'
		},
		{
			value: 72,
			text: '72'
		},
		{
			value: 73,
			text: '73'
		},
		{
			value: 74,
			text: '74'
		},
		{
			value: 75,
			text: '75'
		},
		{
			value: 76,
			text: '76'
		},
		{
			value: 77,
			text: '77'
		},
		{
			value: 78,
			text: '78'
		},
		{
			value: 79,
			text: '79'
		},
		{
			value: 80,
			text: '80'
		},
		{
			value: 81,
			text: '81'
		},
		{
			value: 82,
			text: '82'
		},
		{
			value: 83,
			text: '83'
		},
		{
			value: 84,
			text: '84'
		},
		{
			value: 85,
			text: '85'
		}],
		url: function(data) {
			result = $.ajax({
				type: "GET",
				url: "/cms/SetSingleUserInfo?id=" + g_profile_id + "&gender=" + search_gender + "&age=" + data["value"],
				async: false,
				cache: false,
				success: function() {
					refreshProfilePicture(g_profile_id)
				}
			})

			return result
		}
	});

	$('#search_edit_introduction').editable({
		showbuttons: 'bottom',
		url: function(data) {
			result = $.ajax({
				type: "GET",
				url: "/cms/SetSingleUserInfo?id=" + g_profile_id + "&gender=" + search_gender + "&introduction=" + data["value"],
				async: false,
				cache: false,
			})

			return result
		}
	});

	$('#search_edit_viplevel').editable({
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
				url: "/cms/SetSingleUserInfo?id=" + g_profile_id + "&gender=" + search_gender + "&viplevel=" + data["value"],
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

	$('#search_pagejump').bind('keypress', function(event) {
		if (event.keyCode == "13") {
			searchpagejump()
		}
	});

	$('#search_input').bind('keypress', function(event) {
		if (event.keyCode == "13") {
			dosearch(search_gender)
		}
	});
});
