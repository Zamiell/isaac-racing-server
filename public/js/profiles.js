$(document).ready(function () {
    ConvertTimeStamps();
});

function ConvertTimeStamps() {
    var m_names = new Array("Jan", "Feb", "Mar", "Apr", "May", "June", "July", "Aug", "Sept", "Oct", "Nov", "Dec");
    var d_names = new Array("Sun", "Mon", "Tue", "Wed", "Thur", "Fri", "Sat");
    $('#profiles tr').each(function() {
        if ($(this).find('td').eq(1).html() != undefined) {
            // Miserable hack for Safari's JS strictness
            dt = new Date($(this).find("td").eq(1).html().replace(/\s/, '').replace(/\s/, 'T').replace(' +0000 UTC', ''));

            var curr_hour = dt.getHours();
            var curr_min = dt.getMinutes();
            var curr_time = ((curr_hour < 10) ? "0" + curr_hour : curr_hour) + ":" + ((curr_min < 10) ? "0" + curr_min : curr_min);
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

            $(this).find("td").eq(1).html(d_names[dt.getDay()] + ", " + m_names[dt.getMonth()] + " " + curr_date + sup + ", " + curr_time);
        }
    });
};
