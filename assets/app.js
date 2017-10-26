$(function(){
	var video = document.getElementById('video');
	if (navigator.mediaDevices && navigator.mediaDevices.getUserMedia) {
		navigator.mediaDevices.getUserMedia({ video: true }).then(function (stream) {
			video.src = window.URL.createObjectURL(stream);
			video.play();
		});
	}

	var canvas = document.getElementById('canvas');
	var context = canvas.getContext('2d');
	var video = document.getElementById('video');

	// Trigger photo take
	var button = $('#snap')
	button.click(function(){
		button.addClass('loading')
		$('.info.message').hide()
		context.drawImage(video, 0, 0, 400, 225);
		var dataURL = canvas.toDataURL();
		$.ajax({
			type: "POST",
			url: "/webFaceID",
			data: {
				imgBase64: dataURL
			},
			success: function(resp){
				console.info(resp)
				button.empty().append(
					$("<i>", {class:"camera icon"})
				).addClass('teal').removeClass('green')
				if (resp.faces_len == 0) {
					$('.info.message').text("We didn't see a face").fadeIn()
					return
				}
				if (resp.faces_len > 1) {
					$('.info.message').text("You must be alone to use Web Face ID securely").fadeIn()
					return
				}
				if (!resp.matched) {
					button.transition("shake")
					return
				}
				$('.info.message').text("Hello " + resp.name).fadeIn()
				button.empty().append(
					$("<i>", {class:"check icon"})
				).removeClass('teal').addClass('green').transition('tada')
			},
			error: function(){
				$('.info.message').text("Oops, something went wrong").fadeIn()
			},
			complete: function(){
				button.removeClass('loading')
			}
		})
		
	})

})
