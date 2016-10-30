var view=true;

function refreshEditable() {
    $('#inplaceediting-note').editable();
    $('#inplaceediting-pencil').click(function(e) {
        e.stopPropagation();
        e.preventDefault();
        $('#inplaceediting-note').editable('toggle');
    });
}

function editorBtn() {
    $("#article-content").toggle()

    if (view == true) {
        var options = {resizeType:1}
        KindEditor.create('#editor_id',  options);
        $('#editor_id').html($("#article-content").html())

        $("#article-btn").html("保存")
    } else {
        KindEditor.remove("#editor_id")
        content=$('#editor_id').val()
        $("#article-content").html(content)

        $("#article-btn").html("编辑")
    }
    view = !view
}

$(document).ready(function() {
    refreshEditable()
	$("#editor_id").hide()
});
