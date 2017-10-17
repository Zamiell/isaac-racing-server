let activeLeaderboard = 'unseeded';
let transition = false;

function showLeaderboard(type) {
    // Header buttons
    $('#leaderboard-seeded-button').addClass('inactive');
    $('#leaderboard-unseeded-button').addClass('inactive');
    $('#leaderboard-other-button').addClass('inactive');
    $('#leaderboard-' + type + '-button').removeClass('inactive');

    // Fade out the old leaderboard and fade in the new one
    transition = true;
    $('#leaderboard-' + activeLeaderboard).fadeOut(350, function() {
        $('#leaderboard-' + type).fadeIn(350, function() {
            activeLeaderboard = type;
            transition = false;
        });
    });
}

function pad(n, width, z) {
  z = z || '0';
  n = n + '';
  return n.length >= width ? n : new Array(width - n.length + 1).join(z) + n;
}

$(document).ready(function () {
        ConvertAvgTime();
        ConvertRealTime();
        ConvertFastestTime();
        AdjustRank();
        ConvertTimeStamp();
        ConvertForfeitPenalty();    
});
        
function ConvertAvgTime() {
    $('#leaderboard-' + activeLeaderboard + ' td.lb-adj-avg').each(function(){ 
        time = $(this).html();
        $(this).html(Math.floor(time/1000/60) + ":" + pad(Math.floor(time/1000%60),2));
    });

};
function ConvertRealTime() {
    $('#leaderboard-' + activeLeaderboard + ' td.lb-real-avg').each(function(){ 
        time = $(this).html();
        $(this).html(Math.floor(time/1000/60) + ":" + pad(Math.floor(time/1000%60),2));
    });

};
function ConvertFastestTime(){
    $('#leaderboard-' + activeLeaderboard + ' td.lb-fastest').each(function(){ 
        time = $(this).html();
        $(this).html(Math.floor(time/1000/60) + ":" + pad(Math.floor(time/1000%60),2));
    });

};
function AdjustRank() {
    $('#leaderboard-' + activeLeaderboard + ' td.lb-rank').each(function(){ 
        $(this).html(parseInt($(this).html()) + 1);
    });

};

function ConvertForfeitPenalty() {
    $('#leaderboard-' + activeLeaderboard + ' td.lb-for-pen').each(function(){ 
        time = $(this).html();
        $(this).html(Math.floor(time/1000/60) + ":" + pad(Math.floor(time/1000%60),2));
    });

};
function ConvertTimeStamp() {
        var m_names = new Array("Jan", "Feb", "Mar", "Apr", "May", "June", "July", "Aug", "Sept", "Oct", "Nov", "Dec");
        var d_names = new Array("Sun", "Mon", "Tue", "Wed", "Thur", "Fri", "Sat");

    $('#leaderboard-' + activeLeaderboard + ' td.lb-last-race').each(function(){ 

                        // Miserable hack to help with Safari's strict JS date restrictions 
                        dt = new Date($(this).html().replace(/\s/, 'T').replace(' +0000 UTC', ''));
                        
                        var curr_hours = dt.getHours();
                        var curr_min = dt.getMinutes();
                        var curr_time = curr_hours + ":" + ((curr_min < 10) ? "0" + curr_min : curr_min);
                        var curr_date = dt.getDate();
                        var sup = "";
                        if (curr_date == 1 || curr_date == 21 || curr_date == 31)
                           {
                           sup = "st";
                           }
                        else if (curr_date == 2 || curr_date == 22)
                           {
                           sup = "nd";
                           }
                        else if (curr_date == 3 || curr_date == 23)
                           {
                           sup = "rd";
                           }
                        else
                           {
                           sup = "th";
                           }
                         
                         $(this).html(d_names[dt.getDay()] + ", " + m_names[dt.getMonth()] + " " + dt.getDate() + sup + ", " + dt.getFullYear());

        });
};


/*
Disabling this for now since there is only one leaderboard currently

$('#leaderboard-seeded-button').click(function() {
    if (activeLeaderboard !== 'seeded' && transition === false) {
        showLeaderboard('seeded');
    }
});
$('#leaderboard-unseeded-button').click(function() {
    if (activeLeaderboard !== 'unseeded' && transition === false) {
        showLeaderboard('unseeded');
    }
});
$('#leaderboard-other-button').click(function() {
    if (activeLeaderboard !== 'other' && transition === false) {
        showLeaderboard('other');
    }
});
*/
