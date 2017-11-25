$(document).ready(function () {
        // Season 1 R+9 functions
        ConvertTimeStamps('td.s1r9-td-date');
        ConvertTimes('td.s1r9-td-time');
        $('#season1r9-table').tablesorter({
          headers:{
            '.hof-th-date, .hof-th-proof' : {
              sorter: false
            }
          }
        });
        // Season 1 R+14 functions
        // Season 2 R+7 functions
        // Season 3 R+7 functions
});

function pad(n, width, z) {
    z = z || '0';
    n = n + '';
    return n.length >= width ? n : new Array(width - n.length + 1).join(z) + n;
}

function ConvertTimeStamps(td) {
    var m_names = new Array("Jan", "Feb", "Mar", "Apr", "May", "June", "July", "Aug", "Sept", "Oct", "Nov", "Dec");
    var d_names = new Array("Sun", "Mon", "Tue", "Wed", "Thur", "Fri", "Sat");

    // Miserable hack to help with Safari's strict JS date restrictions
    $(td).each(function(){
      //dt = new Date($(this).html().replace(/\s?/, '').replace(/\s/, 'T').replace(' +0000 UTC', ''));
      dt = new Date($(this).html());
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
      $(this).html(d_names[dt.getDay()] + ", " + m_names[dt.getMonth()] + " " + curr_date + sup + ', ' + dt.getFullYear());
    });
};

function ConvertTimes (td) {
  $(td).each(function(){
    t = $(this).html();
    s = pad(Math.floor(t % 60), 2);
    m = pad(Math.floor(t / 60 % 60), 2);
    h = Math.floor(t / 60 / 60 % 24);
    $(this).html(h + "h " + m + "m " + s + "s")
  });
};
