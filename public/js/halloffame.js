const LAST_RANKED_SOLO_SEASON = 2;

const tableIDs = [
  "season1R9AB",
  "season1R14AB",
  "season2R7AB",
  "season3R7AB",
  "season4R7AB",
  "season5R7AB",
  "season6R7AB",
  "season7R7AB",
  "season8R7AB",
  "season1R7Rep",
  "season2R7Rep",
];
for (i = 1; i <= LAST_RANKED_SOLO_SEASON; i++) {
  tableIDs.push(`season${i}RankedSolo`);
}

let activeLeaderboard = tableIDs[0];
let transition = false;

$(document).ready(() => {
  ConvertTimeStamps("td.td-date");
  ConvertTimes("td.td-time");
  ConvertForfeitRate("td.td-forfeit-rate-50", 50);
  ConvertForfeitRate("td.td-forfeit-rate-100", 100);

  hideAllBoards();
  selectLeaderboard(activeLeaderboard);
});

function ConvertTimeStamps(td) {
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
    $(this).html(
      d_names[dt.getDay()] +
        ", " +
        m_names[dt.getMonth()] +
        " " +
        curr_date +
        sup +
        ", " +
        dt.getFullYear()
    );
  });
}

function ConvertTimes(td) {
  $(td).each(function () {
    const t = $(this).html();
    const seconds = pad(Math.floor(t % 60), 2);
    const minutes = pad(Math.floor((t / 60) % 60), 2);
    const hours = Math.floor((t / 60 / 60) % 24);
    // const timeString = h + "h " + m + "m " + s + "s"
    const timeString =
      hours === 0 ? `${minutes}:${seconds}` : `${hours}:${minutes}:${seconds}`;
    $(this).html(timeString);
  });
}

function ConvertForfeitRate(td, numRaces) {
  $(td).each(function () {
    const numForfeits = $(this).html();
    const forfeitRate = numForfeits / numRaces;
    const forfeitPercent = Math.round(forfeitRate * 100);
    const forfeitRateString = `${forfeitPercent}% (${numForfeits}/${numRaces})`;
    $(this).html(forfeitRateString);
  });
}

function hideAllBoards() {
  for (const tableID of tableIDs) {
    $(`#hof-${tableID}`).css("display", "none");
  }
}

function selectLeaderboard(type) {
  transition = true;

  for (const tableID of tableIDs) {
    if (type === tableID) {
      $("#hof-" + activeLeaderboard).fadeOut(350, () => {
        $("#hof-" + type).fadeIn(350, () => {
          activeLeaderboard = type;
          transition = false;
        });
      });
    }
  }
}
