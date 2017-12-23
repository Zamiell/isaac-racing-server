$(document).ready(function () {
        ConvertTimeRaceStamps();
        $('.tooltip').tooltipster({
          theme: 'tooltipster-shadow'
        });

});

function ConvertTimeRaceStamps() {
        var m_names = new Array("Jan", "Feb", "Mar", "Apr", "May", "June", "July", "Aug", "Sept", "Oct", "Nov", "Dec");
        var d_names = new Array("Sun", "Mon", "Tue", "Wed", "Thur", "Fri", "Sat");

        $('.races-td-date').each(function(){
                //if ($(this).find('td').is('.race-td-date')) {
                        // Miserable hack to help with Safari's strict JS date restrictions
                        dt = new Date($(this).html().replace(/\s/, '').replace(/\s/, 'T').replace(' +0000 UTC', ''));
                        var curr_time = dt.getHours() + ":" + ((dt.getMinutes() < 10) ? "0" + dt.getMinutes() : dt.getMinutes());
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
                         $(this).html(d_names[dt.getDay()] + ", " + m_names[dt.getMonth()] + " " + curr_date + sup + ", " + curr_time);
             //   }
        });
};
