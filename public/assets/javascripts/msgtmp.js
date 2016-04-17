function doAddMessage() {
	var dt = $("#msgtemplate_table").dataTable();
	if (null == dt) {
		return
	}

	var strbuf = $("#msgtemplate_textarea").val();
	var words = strbuf.split(/[\r|\n]+/)
	var haserr = false

	$.each(words, function(i, w) {
		w = w.trim()
		if (0 != w.length) {
			$.getJSON("/cms/msgtemplateadd?template=" + w, function(item) {
				if (null != item) {
					dt.fnAddData(['<input id="dt-tr-' + item["Id"] + '" type="checkbox"/>',
                                  item["Id"], item["Template"],
                                  '<input class="btn btn-mini" type="button" style="float:left" value="修改" onclick="doModify(' + item["Id"] + ')">' +
                                  '<input class="btn btn-mini" type="button" style="float:right" value="删除"' +
                                  ' onclick="doDelete(' + item["Id"] + ')">'])
				} else {
					haserr = true
					alert("发生错误,请检查网络!")
				}
			}).fail(function() {
				haserr = true
				alert("发生错误,请检查网络!")
			})
		}
	})

	if (false == haserr) {
		$("#msgtemplate_textarea").val("")
	}
}

function doDeleteSelected() {
	var dt = $("#msgtemplate_table").dataTable();
	if (null == dt) {
		return
	}

	var nodes = dt.fnGetNodes()
	if (null != nodes) {
		for (var i = nodes.length - 1; i >= 0; i--) {
			if (true == nodes[i].childNodes[0].childNodes[0].checked) {
                var idstr = nodes[i].childNodes[0].childNodes[0].id
	            var words = idstr.split('-')

                $.ajax({
                    type: "GET",
                    url: "/cms/msgtemplatedel?id=" + words[2],
                    async: false, //设为false就是同步请求
                    cache: false,
                    success: function () {
				        dt.fnDeleteRow(i)
                    }
                });
            }
		}
	}
}

function doDelete(templateId) {
	var dt = $("#msgtemplate_table").dataTable();
	if (null == dt) {
		return
	}

	$.get("/cms/msgtemplatedel?id=" + templateId, "", function() {
		var trid = "dt-tr-" + templateId
		var nodes = dt.fnGetNodes()
		if (null != nodes) {
		    for (var i = nodes.length - 1; i >= 0; i--) {
				if (trid == nodes[i].childNodes[0].childNodes[0].id) {
					dt.fnDeleteRow(i)
					break
				}
			}
		}
	}).fail(function() {
		alert("发生错误,请检查网络!")
	})
}

function doModify(templateId) {
    var tds = $("#dt-tr-" + templateId).parent().parent().find("td")
    tds[1].innerHTML = '<input type="text" value="' + tds[1].innerHTML +'" class="span12"/>'
    tds[2].innerHTML = '<input class="btn btn-mini" type="button" style="float:left" value="保存" onclick="doModifySave(' + templateId +
                      ')"><input class="btn btn-mini" type="button" style="float:right" value="删除"' +
                      ' onclick="doDelete(' + templateId + ')">'
}

function doModifySave(templateId) {
    var tds = $("#dt-tr-" + templateId).parent().parent().find("td")
    var template = tds[1].childNodes[0].value.trim()

    if (0 != template.length) {
        $.ajax({
            type: "GET",
            url: "/cms/msgtemplatemodify?id=" + templateId + '&template=' + template,
            async: false, //设为false就是同步请求
            cache: false,
            success: function () {
                tds[1].innerHTML = template;
                tds[2].innerHTML = '<input class="btn btn-mini" type="button" style="float:left" value="修改" onclick="doModify(' + templateId + ')">' +
                      '<input class="btn btn-mini" type="button" style="float:right" value="删除"' +
                      ' onclick="doDelete(' + templateId + ')">'
            }
        });
    }
}

$(document).ready(function() {
	$.getJSON("/cms/msgtemplate", function(data) {
		if (null != data) {
			var bodystr = ""

			$.each(data, function(i, item) {
				str = '<tr>' +
                      '<td><input id="dt-tr-' + item["Id"] + '" type="checkbox"/></td>' +
                      '<td>' + item["Id"] + '</td>' +
                      '<td>' + item["Template"] + '</td>' +
                      '<td><input class="btn btn-mini" type="button" style="float:left" value="修改" onclick="doModify(' + item["Id"] + ')">' +
                      '<input class="btn btn-mini" type="button" style="float:right" value="删除"' +
                      ' onclick="doDelete(' + item["Id"] + ')"></td></tr>'
				bodystr += str
			})

			$("#msgtemplate_table_tbody").html(bodystr);
			$("#msgtemplate_table").dataTable({
                bPaginate: true,
				bAutoWidth: false,
				aLengthMenu: [[10, 20, 50, 100, - 1], [10, 20, 50, 100, "All"]],
				aoColumnDefs: [{
					"bVisible": false,
					"aTargets": [1]
				}],
				bDestroy: true,
				aaSortingFixed: [[1, "desc"]],
				sPaginationType: "bootstrap"
			});
		} else {
			alert("发生错误,请检查网络!")
		}
	}).fail(function() {
		alert("发生错误,请检查网络!")
	})

	$("#msgtemplate_textarea").focus();
})
