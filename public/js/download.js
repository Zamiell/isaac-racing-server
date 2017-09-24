$(document).ready(function () {
    GetLatestReleaseInfo();
});

function GetLatestReleaseInfo() {
    $.getJSON('https://api.github.com/repos/Zamiell/isaac-racing-client/releases/latest').done(function(json) {
        // Make the Windows download button
        var downloadURL = json.assets[3].browser_download_url;
        // The elements of the array correspond to the order that the files are uploaded
        // 0 - latest.yml
        // 1 - RacingPlus-#.##.##-ia32.nsis.7z
        // 2 - RacingPlus-#.##.##-x64.nsis.7z
        // 3 - RacingPlus-WebSetup-#.##.##.exe
        $('#download-button-windows').attr('href', downloadURL);

        // Make the OS X download button
        // TODO

        // Set the current version and release date
        $('#version').html(json.tag_name);
        var released = new Date(json.published_at);
        $('#released').html(released);
   });
}
