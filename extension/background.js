chrome.contextMenus.create({
    title: "Download YouTube video to server",
    contexts: ["video"],
    documentUrlPatterns: ["*://*.youtube.com/*"],
    id: "youtubeDownloadMenu"
});

chrome.contextMenus.create({
    title: "Save video to server",
    contexts: ["video"],
    documentUrlPatterns: ["*://*/*"],
    targetUrlPatterns: ["*://*/*"],
    id: "generalDownloadMenu"
});

chrome.contextMenus.onClicked.addListener((info, tab) => {
    if (info.menuItemId === "youtubeDownloadMenu") {
        downloadVideo(tab.url, "youtube");
    } else if (info.menuItemId === "generalDownloadMenu") {
        downloadVideo(info.srcUrl, "general");
    }
});

function downloadVideo(videoUrl, type) {
    chrome.storage.sync.get(['bestQuality', 'selectedQuality'], function(data) {
        let quality = 'best';

        if (!data.bestQuality && data.selectedQuality) {
            quality = data.selectedQuality;
        }

    fetch('http://localhost:8080/videos/download', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ video_url: videoUrl, type: type, quality: quality}),
    })
        .then(response => response.json())
        .then(data => {
            console.log('Video saved successfully:', data);
        })
        .catch(error => {
            console.error('Error saving video:', error);
        });
    });
}