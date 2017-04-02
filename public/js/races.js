$(document).ready(function () {
	ConvertTimeStamps();
});

function ConvertTimeStamps() {
	$('#race-listing-table tr').each(function(){
		if ($(this).find('td').is('#racedate')) {
			dt = Math.floor($(this).find("#racedate").eq(0).html()/1000);
			console.log(dt);
			var d = new Date(0);
			d.setUTCSeconds(dt);
			$(this).find("td#racedate").html((d.getMonth()+1) + '/' + d.getDate() + '/' + d.getFullYear() + ' ' + d.getHours() + ':' + ("0"+d.getMinutes()).slice(-2));
		}
	});
};
