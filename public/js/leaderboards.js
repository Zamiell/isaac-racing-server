let activeLeaderboard = 'unseeded';
let transition = false;

function showLeaderboard(type) {
    // Header buttons
    $('#leaderboard-seeded-button').addClass('inactive');
    $('#leaderboard-unseeded-button').addClass('inactive');
    $('#leaderboard-diversity-button').addClass('inactive');
    $('#leaderboard-other-button').addClass('inactive');
    $('#leaderboard-' + type + '-button').removeClass('inactive');

    // Fade out the old leaderboard and fade in the new one
    transition = true;


    if (type == 'unseeded') {
        $('#leaderboard-' + activeLeaderboard).fadeOut(350, function() {
            $('#leaderboard-' + type).add('#unseeded-notes-banner').add('#unseeded-notes').fadeIn(350, function() {
                activeLeaderboard = type;
                transition = false;
            });
        });
    } else {
        $('#leaderboard-' + activeLeaderboard).add('#unseeded-notes-banner').add('#unseeded-notes').fadeOut(350, function() {
            $('#leaderboard-' + type).fadeIn(350, function() {
                activeLeaderboard = type;
                transition = false;
            });
        });        
    }
}

function pad(n, width, z) {
    z = z || '0';
    n = n + '';
    return n.length >= width ? n : new Array(width - n.length + 1).join(z) + n;
}

$(document).ready(function() {
    CheckForHash();
    // Unseeded things
    $('#leaderboard-unseeded-table').tablesorter(); 
    AdjustRank('unseeded');
    ConvertTime('unseeded','lb-adj-avg');
    ConvertTime('unseeded','lb-real-avg');
    ConvertTime('unseeded','lb-fastest');
    ConvertTime('unseeded','lb-for-pen');
    ConvertTimeStamp('unseeded','td.lb-last-race a');
    ConvertForfeitRate('unseeded','lb-num-for');

    //Diversity things
    $('#leaderboard-diversity-table').tablesorter(); 
    AdjustRank('diversity');
    ConvertTime('diversity','lb-fastest');
    ConvertTimeStamp('diversity','td.lb-last-race a');
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
        //console.log(window.location.hash.substr(1));  
        type = window.location.hash.substr(1);
        if (type == 'diversity' || type == 'unseeded') {
            showLeaderboard(type);
        } else {
            showLeaderboard('unseeded');
        }
    } 
}


/*$('#leaderboard-seeded-button').click(function() {
    if (activeLeaderboard !== 'seeded' && transition === false) {
        showLeaderboard('seeded');
    }
});*/
$('#leaderboard-unseeded-button').click(function() {
    if (activeLeaderboard !== 'unseeded' && transition === false) {
        showLeaderboard('unseeded');
    }
});
$('#leaderboard-diversity-button').click(function() {
    if (activeLeaderboard !== 'diversity' && transition === false) {
        showLeaderboard('diversity');
    }
});

