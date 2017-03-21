$(document).ready(function () {
    ConvertTimeStamps();
});

function ConvertTimeStamps() {
    //$('#profiles').each(function() {
    $('#profiles tr').each(function(){
        dt = $(this).find("td").eq(1).html()/1000;
        var d = new Date(0);
        d.setUTCSeconds(dt);
        month = d.getMonth();
	year = d.getFullYear();
	day = d.getDate();
        $(this).find("td").eq(1).html(month + '/' + day + '/' + year);
    });
    
};
