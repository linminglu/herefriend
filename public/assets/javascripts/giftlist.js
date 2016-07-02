function refreshGiftVerboseDlg(id) {
	$.ajax({
		type: "GET",
		url: "/cms/GetGiftVerbose?id=" + id,
		dataType: "json",
		async: false,
		cache: false,
		success: function(data) {
			if ((null == data) || (0 == data.length)) {
				return 0
			}
            $("#gift_verbose_name").html(data["Name"])
		},
		error: function() {
			alert("发生错误,请检查网络!")
		}
	});
}

function getGiftItemHtml(item) {
    var str = '<li class="respl-item category-10">'
            + '<div class="unit">'
            + '     <p class="hoverline"></p>'
            + '     <a class="images lazy"><img src="' + item["ImageUrl"] +'"></a>'
            + '     <h3><a>' + item["Name"] + '</a></h3>'
            + '     <dl>'
            + '         <dd>价格：' + item["Price"] + '金币</dd>'
            + '         <dd>已售/剩余：' + 0 + '/' + item["ValidNum"] +'个</dd>'
            + '     </dl>'
            + '     <div class="view">'
            + '         <a class="btn btn-mini right" data-toggle="modal" onclick="refreshGiftVerboseDlg(' + item["Id"] + ')" href="#gift_verbose_dlg" role="button">详情<span>&gt;</span></a>'
            + '     </div>'
            + '</div>'
            + '</li>'

    return str
}

function load_GiftItems() {
	$.ajax({
		type: "GET",
		url: "/cms/GetGiftList",
		dataType: "json",
		async: false,
		cache: false,
		success: function(data) {
			if ((null == data) || (0 == data.length)) {
				return 0
			}

			var listStr = ""
			$.each(data, function(i, item) {
				listStr += getGiftItemHtml(item)
			});

			$("#giftlist-items").html(listStr)
			scroll(0, 0)

			$("img.lazy").lazyload({
				effect: "fadeIn",
				failure_limit: 6,
				effect_speed: 1000
			});
		},
		error: function() {
			alert("发生错误,请检查网络!")
		}
	});
}

jQuery(document).ready(function($) {
load_GiftItems();
;(function(element) {
		var $respl = $(element);
		var $container = $('.respl-items', $respl);

		$container.imagesLoaded(function() {
			$container.isotope({
				containerStyle: {
					position: 'relative',
					height: 'auto',
					overflow: 'visible'
				},
				itemSelector: '.respl-item',
				sortAscending: true
			});
			_opTionSets();	

			function _opTionSets(){
				var $optionSets = $('.respl-header .respl-option', $respl),
				$optionLinks = $optionSets.find('a');
				$optionLinks.each(function(){
					$(this).click(function(){
						var $this = $(this);
						var $optionSet = $this.parents('.respl-option');
						$this.parent().addClass('select').siblings().removeClass('select');
						var options = {},key = $optionSet.attr('data-option-key'),value = $this.attr('data-rl_value');
						value = value === 'false' ? false: value;
						options[key] = value;
						$container.isotope(options);
						return false;
					});
				});
			}
		});
	})('#giftlist-content');
	$(".ProductList li").hover(function(){
		$(this).find(".hoverline").fadeIn().end().find(".images").addClass("mousehover").stop().animate({
			"height":"160px"															  
		},500).find("img").stop().animate({
			"margin-top":"-30px"	
		},500);	
		$(this).find(".view").stop().animate({
			"bottom":"0"									 
		},500);
	},function(){
		$(this).find(".hoverline").fadeOut().end().find(".images").removeClass("mousehover").stop().animate({
			"height":"222px"															  
		},500).find("img").stop().animate({
			"margin-top":"0"	
		},500);	
		$(this).find(".view").stop().animate({
			"bottom":"-60px"									 
		},500);
	});
});

