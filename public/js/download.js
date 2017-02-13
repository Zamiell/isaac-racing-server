$(document).ready(function () {
    GetLatestReleaseInfo();
});

function GetLatestReleaseInfo() {
    $.getJSON('https://api.github.com/repos/Zamiell/isaac-racing-client/releases/latest').done(function(json) {
        // Make the Windows download button
        // (the 0th element of the array is always the "latest.yml" file)
        var downloadURL = json.assets[1].browser_download_url;
        $('#download-button-windows').attr('href', downloadURL);

        // Make the OS X download button
        // TODO

        // Set the current version and release date
        $('#version').html(json.tag_name);
        var released = new Date(json.published_at);
        $('#released').html(released);
   });
}
