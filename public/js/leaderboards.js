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
