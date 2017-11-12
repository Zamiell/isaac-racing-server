$(document).ready(function() {
    ConvertTimeProfileStamps();
    BannedUser();
});

function ConvertTimeProfileStamps() {
    var m_names = new Array("Jan", "Feb", "Mar", "Apr", "May", "June", "July", "Aug", "Sept", "Oct", "Nov", "Dec");
    var d_names = new Array("Sun", "Mon", "Tue", "Wed", "Thur", "Fri", "Sat");

    // Miserable hack to help with Safari's strict JS date restrictions
    dt = new Date($('span#date').html().replace(/\s?/, '').replace(/\s/, 'T').replace(' +0000 UTC', ''));
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
    $("span#date").html(d_names[dt.getDay()] + ", " + m_names[dt.getMonth()] + " " + curr_date + sup + ", " + curr_time);
};

function BannedUser() {
    if ($('div#banned').html() == "true") {
        var docWidth = $(document).width();
        var docHeight = $(document).height();
        var overlayDiv = "<div id=\"overlay-div\"></div>";
        $("span#span-ban").append(overlayDiv);
        $("#overlay-div").css({"opacity":"1.0", "position":"fixed","width": docWidth + "px","height":docHeight + "px","text-align":"center","z-index":"10"});
        $("#overlay-div").append("<div id=\"image-div\"></div>");
        $("#image-div").css("position","relative", "left",docWidth/4 + "px","width", docWidth/4);
        $("#image-div").append("<img src=\"/public/img/no.png\"id=\"zoomed-img\" />");
        var imgWidth = $("#image-div").width();
        var imgHeight = $("#image-height").height();
        //$("#image-div").css("position","absolute", "top","10px","left","10px");
    };

};