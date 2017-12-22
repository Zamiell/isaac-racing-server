$(document).ready(function() {
    ConvertTimeProfileStamps('td.races-td-date');
    ConvertTimeProfileStamps('span#join-date');
    ConvertTimeProfileStamps('td.ranked-racedate');
    ConvertTime('#unseeded-adjavg-val');
    ConvertTime('#unseeded-realavg-val');
    ConvertTime('#unseeded-forpen-val');
    ConvertTime('#unseeded-fastest-val');
    ConvertTotalTime('#misc-wasted-time');
    ConvertForfeitRate('#unseeded-numfor-val');
    BannedUser();
});

function ConvertTime(td) {
    $(td).each(function() {
      if ($(this).html() > 0) {
        time = $(this).html();
        $(this).html(Math.floor(time / 1000 / 60) + ":" + pad(Math.floor(time / 1000 % 60), 2));
      }
    });
};

function ConvertTotalTime(td) {
    $(td).each(function() {
      if ($(this).html() > 0) {
        time = $(this).html();
        seconds = Math.floor(time / 1000 % 60);
        minutes = Math.floor(time / 1000 / 60 % 60);
        hours = Math.floor(time / 1000 / 60 / 60 % 24);
        days = Math.floor(time / 1000 / 60 / 60 / 24 );
        elasped = days + ' Day' + ( days != 1 ? 's ' : ' ' ) + hours + ' Hour' + ( hours != 1 ? 's ' : ' ' ) + minutes + ' Minute' + ( minutes != 1 ? 's ' : ' ' ) + seconds + ' Second' + ( seconds != 1 ? 's ' : ' ' );
        $(this).html(elasped);
      }
    });
};

function ConvertTimeProfileStamps(td) {
    var m_names = new Array("Jan", "Feb", "Mar", "Apr", "May", "June", "July", "Aug", "Sept", "Oct", "Nov", "Dec");
    var d_names = new Array("Sun", "Mon", "Tue", "Wed", "Thur", "Fri", "Sat");

    // Miserable hack to help with Safari's strict JS date restrictions
    $(td).each(function(){
      dt = new Date($(this).html().replace(/\s?/, '').replace(/\s/, 'T').replace(' +0000 UTC', ''));
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
    });
}

function ConvertForfeitRate(td) {
    $(td).each(function() {
      if ($(this).html() > 0) {
        num = $(this).html();
        // I have no idea how this works, but somehow it does
        total = ($(this).closest('tr').next('tr').find('#unseeded-numraces-val').html() > 50) ? 50 : $(this).closest('tr').next('tr').find('#unseeded-numraces-val').html();
        rate = num / total * 100;
        rate = Math.round(rate); // Round it to the nearest whole number
        $(this).html(rate + "% (" + num + "/" + total + ")");
      }
    });
};

function BannedUser() {
    if ($('div#banned').html() == "true") {
        var docWidth = $(document).width();
        var docHeight = $(document).height();
        var navHeight = $('#header').height();
        var winHeight = $(window).height();
        var winPercent = navHeight / winHeight;
        var banHeight = $('#banned').height();
        var overlayDiv = "<div id=\"overlay-div\"></div>";
        $("span#span-ban").append(overlayDiv);
        $("#overlay-div").css({"opacity":"1.0", "position":"fixed","width": docWidth + "px","height":docHeight + "px","text-align":"center","z-index":"10"});
        $("#overlay-div").append("<div id=\"image-div\"></div>");
        $("#image-div").css("position","relative", "left",docWidth/4 + "px","width", docWidth/4);
        $("#image-div").append("<img src=\"/public/img/no.png\"id=\"zoomed-img\" />");
        $("#zoomed-img").css({"height": (winHeight - (winHeight * winPercent )) - banHeight + "px"});
        var imgWidth = $("#image-div").width();
        var imgHeight = $("#image-height").height();
    };

};
