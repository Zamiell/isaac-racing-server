let activeLeaderboard = 'unseeded-solo';
let transition = false;
let button_array = ["seeded","seeded-solo","unseeded","unseeded-solo","diversity","other"];

function hideAllNotes() {
  $('#unseeded-notes-banner').css("display","none");
  $('#unseeded-notes').css("display","none");
  $('#diversity-notes-banner').css("display","none");
  $('#diversity-notes').css("display","none");
}

function showLeaderboard(type) {
  transition = true;
  console.log(type);
  hideAllNotes();

  for (var i = 0, len = button_array.length; i < len; i++) {
    $('#leaderboard-' + button_array[i] + '-button').addClass('inactive');
  }
  $('#leaderboard-' + type + '-button').removeClass('inactive');
  $('#leaderboard-' + activeLeaderboard).fadeOut(350, function() {
    $('#leaderboard-' + type).add('#' + type + '-notes-banner').add('#' + type + '-notes').fadeIn(350, function() {
      activeLeaderboard = type;
      transition = false;
    });
  });
};

function pad(n, width, z) {
    z = z || '0';
    n = n + '';
    return n.length >= width ? n : new Array(width - n.length + 1).join(z) + n;
}

$(document).ready(function() {

    hideAllNotes();
    // Unseeded things

    $('#leaderboard-unseeded-solo-table').tablesorter();
    AdjustRank('unseeded-solo');
    ConvertTime('unseeded-solo','lb-adj-avg');
    ConvertTime('unseeded-solo','lb-real-avg');
    ConvertTime('unseeded-solo','lb-fastest');
    ConvertTime('unseeded-solo','lb-for-pen');
    ConvertTimeStamp('unseeded-solo','td.lb-last-race a');
    ConvertForfeitRate('unseeded-solo','lb-num-for');

    //Diversity things
    $('#leaderboard-diversity-table').tablesorter();
    AdjustRank('diversity');
    ConvertTime('diversity','lb-fastest');
    ConvertTimeStamp('diversity','td.lb-last-race a');

    // Starting functions
    showLeaderboard('unseeded-solo');
    CheckForHash();

});

function ConvertTime(leaderboard, tableData) {
    $('#leaderboard-' + leaderboard + ' td.' + tableData).each(function() {
        time = $(this).html();
        $(this).html(Math.floor(time / 1000 / 60) + ":" + pad(Math.floor(time / 1000 % 60), 2));
    });
};

function AdjustRank(leaderboard) {
    $('#leaderboard-' + leaderboard + ' td.lb-rank').each(function() {
        $(this).html(parseInt($(this).html()) + 1);
    });
};

function ConvertForfeitRate(leaderboard, tableData) {
    $('#leaderboard-' + leaderboard + ' td.' + tableData).each(function() {
        num = $(this).html();
        total = ($(this).next().html() > 50) ? 50 : $(this).next().html();
        rate = num / total * 100;
        rate = Math.round(rate); // Round it to the nearest whole number
        $(this).html(rate + "% (" + num + "/" + total + ")");
    });
};

function ConvertTimeStamp(leaderboard, tableData) {
    var m_names = new Array("Jan", "Feb", "Mar", "Apr", "May", "June", "July", "Aug", "Sept", "Oct", "Nov", "Dec");
    var d_names = new Array("Sun", "Mon", "Tue", "Wed", "Thur", "Fri", "Sat");

    $('#leaderboard-' + leaderboard + ' ' + tableData).each(function() {
        // Miserable hack to help with Safari's strict JS date restrictions
        dt = new Date($(this).html().replace(/\s/, 'T').replace(' +0000 UTC', ''));
        var curr_hours = dt.getHours();
        var curr_min = dt.getMinutes();
        var curr_time = curr_hours + ":" + ((curr_min < 10) ? "0" + curr_min : curr_min);
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

        $(this).html(d_names[dt.getDay()] + ", " + m_names[dt.getMonth()] + " " + dt.getDate() + sup + ", " + dt.getFullYear());
    });
};

function CheckForHash() {
    if (window.location.hash) {
        type = window.location.hash.substr(1);
        if (type == 'diversity' || type == 'unseeded' || type == 'unseeded-solo') {
            showLeaderboard(type);
        } else {
            showLeaderboard('unseeded-solo');
        }
    }
}


/*
$('#leaderboard-seeded-button').click(function() {
    if (activeLeaderboard !== 'seeded' && transition === false) {
        showLeaderboard('seeded');
    }
});
*/
/*
$('#leaderboard-seeded-solo-button').click(function() {
    if (activeLeaderboard !== 'seeded-solo' && transition === false) {
        showLeaderboard('seeded-solo');
    }
});
*/
$('#leaderboard-unseeded-button').click(function() {
    if (activeLeaderboard !== 'unseeded' && transition === false) {
        showLeaderboard('unseeded');
    }
});
$('#leaderboard-unseeded-solo-button').click(function() {
    if (activeLeaderboard !== 'unseeded-solo' && transition === false) {
        showLeaderboard('unseeded-solo');
    }
});
$('#leaderboard-diversity-button').click(function() {
    if (activeLeaderboard !== 'diversity' && transition === false) {
        showLeaderboard('diversity');
    }
});
