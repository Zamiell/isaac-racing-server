$(document).ready(function() {
  ConvertTimeStamp('.td-startdate');
});

function ConvertTimeStamp(tableData) {
    var m_names = new Array('Jan', 'Feb', 'Mar', 'Apr', 'May', 'June', 'July', 'Aug', 'Sept', 'Oct', 'Nov', 'Dec');
    var d_names = new Array('Sun', 'Mon', 'Tue', 'Wed', 'Thur', 'Fri', 'Sat');
    $(tableData).each(function() {
        // Miserable hack to help with Safari's strict JS date restrictions
        dt = new Date($(this).html().replace(' +0000 UTC', '').replace(/\s/, 'T'));
        var curr_hours = dt.getHours();
        var curr_min = dt.getMinutes();
        var curr_time = curr_hours + ':' + ((curr_min < 10) ? '0' + curr_min : curr_min);
        var curr_date = dt.getDate();
        var sup = '';
        if (curr_date == 1 || curr_date == 21 || curr_date == 31) {
            sup = 'st';
        } else if (curr_date == 2 || curr_date == 22) {
            sup = 'nd';
        } else if (curr_date == 3 || curr_date == 23) {
            sup = 'rd';
        } else {
            sup = 'th';
        }

        $(this).html(d_names[dt.getDay()] + ', ' + m_names[dt.getMonth()] + ' ' + dt.getDate() + sup + ', ' + dt.getFullYear() + " @ " + pad(dt.getHours(),2) + ":" + pad(dt.getMinutes(),2));
    });
};
