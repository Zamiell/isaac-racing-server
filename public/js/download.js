$(document).ready(function () {
    GetLatestReleaseInfo();
});

function GetLatestReleaseInfo() {
    $.getJSON('https://api.github.com/repos/Zamiell/isaac-racing-client/releases/latest').done(function(json) {
        // The elements of the array correspond to the order that the files are uploaded, e.g.
        // 0 - latest.yml
        // 1 - RacingPlus-#.##.##-ia32.nsis.7z
        // 2 - RacingPlus-#.##.##-x64.nsis.7z
        // 3 - RacingPlus-WebSetup-#.##.##.exe

        // Make the Windows download button
        let windowsURL = '';
        for (const asset of json.assets) {
            const url = asset.browser_download_url;
            if (url.endsWith('.exe')) {
                windowsURL = url;
                break;
            }
        }
        if (windowsURL !== '') {
            $('#download-button-windows').attr('href', windowsURL);
        }

        // Make the macOS download button
        let macOSURL = '';
        for (const asset of json.assets) {
            const url = asset.browser_download_url;
            if (url.endsWith('.dmg')) {
                macOSURL = url;
                break;
            }
        }
        if (macOSURL !== '') {
            $('#download-button-macos').attr('href', macOSURL);
        }

        // Set the current version and release date
        $('#version').html(json.tag_name);
        var released = new Date(json.published_at);
        $('#released').html(released);
   });
}
