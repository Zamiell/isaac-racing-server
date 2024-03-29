let activeLeaderboard = "unseeded";
let transition = false;
let button_array = ["seeded", "unseeded", "diversity", "ranked-solo"];
const fadeTime = 350;
const numRankedSoloRaces = 100;

function hideAllNotes() {
  $("#notes-multiplayer").hide(0);
  $("#notes-solo").hide(0);
}

function hideAllBoards() {
  for (const button of button_array) {
    $(`#leaderboard-${button}`).hide(0);
  }
}

function showLeaderboard(type) {
  transition = true;
  hideAllNotes();

  // Set all the buttons inactive
  for (const button of button_array) {
    $(`#leaderboard-${button}-button`).addClass("inactive");
  }

  // Show the current leaderboard button
  $(`#leaderboard-${type}-button`).removeClass("inactive");
  $(`#leaderboard-${activeLeaderboard}`).fadeOut(fadeTime, () => {
    $(`#leaderboard-${type}`).fadeIn(fadeTime);
    const notesID = `#notes-${type.endsWith("-solo") ? type : "multiplayer"}`;
    $(notesID).fadeIn(fadeTime, () => {
      activeLeaderboard = type;
      transition = false;
    });
  });
}

$(document).ready(() => {
  hideAllNotes();
  hideAllBoards();

  // Seeded things
  $("#leaderboard-seeded-table").tablesorter({
    headers: {
      2: { sorter: false },
      6: { sorter: false },
    },
  });
  AdjustRank("seeded");
  ConvertTime("seeded", "lb-fastest");
  ConvertTimeStamp("seeded", "td.lb-last-race a");

  // Unseeded things
  $("#leaderboard-unseeded-table").tablesorter({
    headers: {
      2: { sorter: false },
      6: { sorter: false },
    },
  });
  AdjustRank("unseeded");
  ConvertTime("unseeded", "lb-fastest");
  ConvertTimeStamp("unseeded", "td.lb-last-race a");

  // Ranked solo things
  $("#leaderboard-ranked-solo-table").tablesorter({
    headers: {
      5: { sorter: false },
      9: { sorter: false },
    },
  });
  AdjustRank("ranked-solo");
  ConvertTime("ranked-solo", "lb-adj-avg");
  ConvertTime("ranked-solo", "lb-real-avg");
  ConvertTime("ranked-solo", "lb-fastest");
  ConvertTime("ranked-solo", "lb-for-pen");
  ConvertTimeStamp("ranked-solo", "td.lb-last-race a");
  ConvertForfeitRate("ranked-solo", "lb-num-for");

  // Diversity things
  $("#leaderboard-diversity-table").tablesorter({
    headers: {
      2: { sorter: false },
      6: { sorter: false },
    },
  });
  AdjustRank("diversity");
  ConvertTime("diversity", "lb-fastest");
  ConvertTimeStamp("diversity", "td.lb-last-race a");

  // Starting functions
  showLeaderboard("unseeded");
  CheckForHash();
});

function ConvertTime(leaderboard, tableData) {
  $("#leaderboard-" + leaderboard + " td." + tableData).each(function () {
    time = $(this).html();
    $(this).html(
      Math.floor(time / 1000 / 60) +
        ":" +
        pad(Math.floor((time / 1000) % 60), 2)
    );
  });
}

function AdjustRank(leaderboard) {
  $("#leaderboard-" + leaderboard + " td.lb-rank").each(function () {
    $(this).html(parseInt($(this).html()) + 1);
  });
}

function ConvertForfeitRate(leaderboard, tableData) {
  $("#leaderboard-" + leaderboard + " td." + tableData).each(function () {
    num = $(this).html();
    total =
      $(this).next().html() > numRankedSoloRaces
        ? numRankedSoloRaces
        : $(this).next().html();
    rate = (num / total) * 100;
    rate = Math.round(rate); // Round it to the nearest whole number
    $(this).html(rate + "% (" + num + "/" + total + ")");
  });
}

function ConvertTimeStamp(leaderboard, tableData) {
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
  $("#leaderboard-" + leaderboard + " " + tableData).each(function () {
    // Miserable hack to help with Safari's strict JS date restrictions
    dt = new Date($(this).html().replace(" +0000 UTC", "").replace(/\s/, "T"));
    var curr_hours = dt.getHours();
    var curr_min = dt.getMinutes();
    var curr_time =
      curr_hours + ":" + (curr_min < 10 ? "0" + curr_min : curr_min);
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

    $(this).html(
      d_names[dt.getDay()] +
        ", " +
        m_names[dt.getMonth()] +
        " " +
        dt.getDate() +
        sup +
        ", " +
        dt.getFullYear()
    );
  });
}

function CheckForHash() {
  if (window.location.hash) {
    type = window.location.hash.substr(1);
    if (
      type == "seeded" ||
      type == "diversity" ||
      type == "unseeded" ||
      type == "ranked-solo"
    ) {
      showLeaderboard(type);
    } else {
      showLeaderboard("unseeded");
    }
  }
}

$("#leaderboard-seeded-button").click(() => {
  if (activeLeaderboard !== "seeded" && transition === false) {
    showLeaderboard("seeded");
  }
});

$("#leaderboard-unseeded-button").click(() => {
  if (activeLeaderboard !== "unseeded" && transition === false) {
    showLeaderboard("unseeded");
  }
});

$("#leaderboard-diversity-button").click(() => {
  if (activeLeaderboard !== "diversity" && transition === false) {
    showLeaderboard("diversity");
  }
});

$("#leaderboard-ranked-solo-button").click(() => {
  if (activeLeaderboard !== "ranked-solo" && transition === false) {
    showLeaderboard("ranked-solo");
  }
});
