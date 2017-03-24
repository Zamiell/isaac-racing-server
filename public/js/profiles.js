$(document).ready(function () {
	ConvertTimeStamps();
});

function ConvertTimeStamps() {
	$('#profiles tr').each(function(){
		dt = $(this).find("td").eq(1).html()/1000;
		var d = new Date(0);
		d.setUTCSeconds(dt);
		$(this).find("td").eq(1).html(d.getMonth() + '/' + d.getDate() + '/' + d.getFullYear());
	});
};
