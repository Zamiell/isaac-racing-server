const numRankedSoloUnseeded = 100;

$(document).ready(() => {
  ConvertTimeProfileStamps("td.races-td-date");
  ConvertTimeProfileStampsWithYear("span#join-date");
  ConvertTimeProfileStamps("td.ranked-racedate");
  ConvertTotalTime("#misc-wasted-time");
  ConvertForfeitRate("#ranked-solo-numfor-val");
  ConvertRaceTime("#ranked-solo-adjavg-val");
  ConvertRaceTime("#ranked-solo-realavg-val");
  ConvertRaceTime("#ranked-solo-forpen-val");
  ConvertRaceTime("#ranked-solo-fastest-val");
  ConvertRaceTime(".races-td-time");
  BannedUser();
  $(".tooltip").tooltipster({
    theme: "tooltipster-shadow",
  });
});

function ConvertRaceTime(td) {
  $(td).each(function () {
    runtime = Math.floor($(this).html() / 1000);
    if (runtime) {
      sec = Math.floor(runtime % 60);
      min = Math.floor((runtime / 60) % 60);
      hour = Math.floor((runtime / 60 / 60) % 24);
      time_converted = "";
      if (hour > 0) {
        time_converted = hour + ":";
      }
      time_converted =
        time_converted + (hour > 0 ? pad(min, 2) : min) + ":" + pad(sec, 2);
      $(this).html(time_converted);
    }
  });
}

function ConvertTotalTime(td) {
  $(td).each(function () {
    if ($(this).html() > 0) {
      const time = $(this).html();
      const seconds = Math.floor((time / 1000) % 60);
      const minutes = Math.floor((time / 1000 / 60) % 60);
      const hours = Math.floor((time / 1000 / 60 / 60) % 24);
      const days = Math.floor(time / 1000 / 60 / 60 / 24);
      const elapsed =
        days +
        " Day" +
        (days != 1 ? "s " : " ") +
        hours +
        " Hour" +
        (hours != 1 ? "s " : " ") +
        minutes +
        " Minute" +
        (minutes != 1 ? "s " : " ") +
        seconds +
        " Second" +
        (seconds != 1 ? "s " : " ");
      $(this).html(elapsed);
    }
  });
}

function ConvertTimeProfileStamps(td) {
  var m_names = new Array(
    "Jan",
    "Feb",
    "Mar",
    "Apr",
    "May",
    "June",
    "July",
    "Aug",
    "Sept",
    "Oct",
    "Nov",
    "Dec"
  );
  var d_names = new Array("Sun", "Mon", "Tue", "Wed", "Thur", "Fri", "Sat");

  // Miserable hack to help with Safari's strict JS date restrictions
  $(td).each(function () {
    dt = new Date(
      $(this)
        .html()
        .replace(/\s?/, "")
        .replace(/\s/, "T")
        .replace(" +0000 UTC", "")
    );
    var curr_time =
      dt.getHours() +
      ":" +
      (dt.getMinutes() < 10 ? "0" + dt.getMinutes() : dt.getMinutes());
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
    $(this).html(
      d_names[dt.getDay()] +
        ", " +
        m_names[dt.getMonth()] +
        " " +
        curr_date +
        sup +
        ", " +
        curr_time
    );
  });
}

function ConvertTimeProfileStampsWithYear(td) {
  var m_names = new Array(
    "Jan",
    "Feb",
    "Mar",
    "Apr",
    "May",
    "June",
    "July",
    "Aug",
    "Sept",
    "Oct",
    "Nov",
    "Dec"
  );
  var d_names = new Array("Sun", "Mon", "Tue", "Wed", "Thur", "Fri", "Sat");

  // Miserable hack to help with Safari's strict JS date restrictions
  $(td).each(function () {
    dt = new Date(
      $(this)
        .html()
        .replace(/\s?/, "")
        .replace(/\s/, "T")
        .replace(" +0000 UTC", "")
    );
    var curr_time =
      dt.getHours() +
      ":" +
      (dt.getMinutes() < 10 ? "0" + dt.getMinutes() : dt.getMinutes());
    var curr_date = dt.getDate();
    var curr_year = dt.getFullYear();
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
    $(this).html(
      d_names[dt.getDay()] +
        ", " +
        m_names[dt.getMonth()] +
        " " +
        curr_date +
        "<sup>" +
        sup +
        "</sup>" +
        " " +
        curr_year +
        " @ " +
        curr_time
    );
  });
}

function ConvertForfeitRate(td) {
  $(td).each(function () {
    if ($(this).html() > 0) {
      num = $(this).html();
      // I have no idea how this works, but somehow it does
      total =
        $(this)
          .closest("tr")
          .next("tr")
          .find("#ranked-solo-numraces-val")
          .html() > numRankedSoloUnseeded
          ? numRankedSoloUnseeded
          : $(this)
              .closest("tr")
              .next("tr")
              .find("#ranked-solo-numraces-val")
              .html();
      rate = (num / total) * 100;
      rate = Math.round(rate); // Round it to the nearest whole number
      $(this).html(rate + "% (" + num + "/" + total + ")");
    }
  });
}

function BannedUser() {
  if ($("div#banned").html() == "true") {
    var docWidth = $(document).width();
    var docHeight = $(document).height();
    var navHeight = $("#header").height();
    var winHeight = $(window).height();
    var winPercent = navHeight / winHeight;
    var overlayDiv = '<div id="overlay-div"></div>';
    $("span#span-ban").append(overlayDiv);
    $("#overlay-div").css({
      opacity: "1.0",
      position: "fixed",
      width: docWidth + "px",
      height: docHeight + "px",
      "text-align": "center",
      "z-index": "10",
      "margin-top": "0.3em",
    });
    $("#overlay-div").append('<div id="image-div"></div>');
    $("#image-div").css(
      "position",
      "relative",
      "left",
      docWidth / 4 + "px",
      "width",
      docWidth / 4
    );
    $("#image-div").append('<img src="/public/img/no.png"id="zoomed-img" />');
    $("#zoomed-img").css({ height: winHeight - winHeight * winPercent + "px" });
    var imgWidth = $("#image-div").width();
    var imgHeight = $("#image-height").height();
  }
}
