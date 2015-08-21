document.addEvent('domready', function () {
    setInterval(refreshScreenshots, 3000);
});

function refreshScreenshots () {
    var imgs = $$('.screenshot');
    var deviceNames = Array.from(imgs).filter(function (img) {
        return !img.hasClass('offline');
    }).map(function (img) {
        return img.getProperty('alt');
    });
    console.log(deviceNames);
    new Request.JSON({
        url: '/update',
        onSuccess: handleScreenshotUpdate
    }).post({
        devices: JSON.encode(deviceNames)
    });
}

function handleScreenshotUpdate (screenshots) {
    Object.each(screenshots, function (screenshot, deviceID) {
        screenshot = 'devices/' + deviceID + '/screenshots/' + screenshot;
        $$('img[alt="' + deviceID + '"]').setProperty('src', screenshot);
    });
}
