$(document).ready(function () {
	ConvertTimeStamps();
});

function ConvertTimeStamps() {
	$('#race-listing-table tr').each(function(){
		if ($(this).find('td').is('#racedate')) {
			dt = Math.floor($(this).find("#racedate").eq(0).html()/1000);
			var d = new Date(0);
			d.setUTCSeconds(dt);
			$(this).find("td#racedate").html((d.getUTCMonth()+1) + '/' + d.getUTCDate() + '/' + d.getUTCFullYear() + ' ' + d.getUTCHours() + ':' + ("0"+d.getUTCMinutes()).slice(-2) + ' UTC');
		}
	});
};
