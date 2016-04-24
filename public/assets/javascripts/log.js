function refreshLog() {
	$.get("/cms/log", "", function(data) {
		$("#log-content").html("<pre>"+ data + "</pre>")
	}).fail(function() {
		alert("发生错误,请检查网络!")
	})

	window.setTimeout(refreshLog, 5000)
}

$(document).ready(function() {
	refreshLog()
});

