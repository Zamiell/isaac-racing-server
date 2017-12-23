$(document).ready(function () {
	ConvertTimeStamp('profiles', 'td.profile-date-created');
	ConvertTimeStamp('profiles', 'td.profile-last-race a');
	$('.tooltip').tooltipster({
		theme: 'tooltipster-shadow'
	});
});

function ConvertTimeStamp(table, tableData) {
	var m_names = new Array("Jan", "Feb", "Mar", "Apr", "May", "June", "July", "Aug", "Sept", "Oct", "Nov", "Dec");
	var d_names = new Array("Sun", "Mon", "Tue", "Wed", "Thur", "Fri", "Sat");
	$('#' + table + ' ' + tableData).each(function() {
		// Miserable hack to help with Safari's strict JS date restrictions
		if ($(this).html() != '') {
			dt = new Date($(this).html().trim().replace(/\s/, 'T').replace(' +0000 UTC', ''));
			var curr_date = dt.getDate();
			var sup = "";
			if (curr_date == 1 || curr_date == 21 || curr_date == 31) {
				sup = "st";
			} else if (curr_date == 2 || curr_date == 22) {
				sup = "nd";
			} else if (curr_date == 3 || curr_date == 23) {
				sup = "rd";
			} else {
				sup = "th";
			}
			// Write the timestamp back
			$(this).html(d_names[dt.getDay()] + ", " + m_names[dt.getMonth()] + " " + dt.getDate() + sup + ", " + dt.getFullYear());
		}
	});
};
