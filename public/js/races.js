$(document).ready(function () {
	ConvertTimeStamps();
});

function ConvertTimeStamps() {
	$('#race-listing-table tr').each(function(){
		if ($(this).find('td').is('#racedate')) {
			dt = $(this).find("#racedate").eq(0).html()/1000;
			var d = new Date(0);
			$(this).find("td#racedate").html((d.getMonth()+1) + '/' + d.getDate() + '/' + d.getFullYear() + ' ' + d.getHours() + ':' + d.getMinutes	());
		}
	});
};
