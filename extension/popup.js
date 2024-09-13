document.addEventListener('DOMContentLoaded', function() {
    const bestQualityCheckbox = document.getElementById('bestQualityCheckbox');
    const qualitySelect = document.getElementById('qualitySelect');
    const saveButton = document.getElementById('saveSettings');

    bestQualityCheckbox.disabled = false;
    chrome.storage.sync.get(['bestQuality', 'selectedQuality'], function(data) {
        if (data.bestQuality === undefined) {
            chrome.storage.sync.set({ bestQuality: true });
            bestQualityCheckbox.checked = true;
            qualitySelect.disabled = true;
        } else {
            bestQualityCheckbox.checked = data.bestQuality;
            qualitySelect.disabled = data.bestQuality;
            if (data.selectedQuality) {
                qualitySelect.value = data.selectedQuality;
            }
        }
    });

    bestQualityCheckbox.addEventListener('change', function() {
        qualitySelect.disabled = !!bestQualityCheckbox.checked;
    });

    saveButton.addEventListener('click', function() {
        const bestQuality = bestQualityCheckbox.checked;
        const selectedQuality = qualitySelect.value;

        chrome.storage.sync.set({
            bestQuality: bestQuality,
            selectedQuality: selectedQuality
        }, function() {
            alert('Settings saved!');
        });
    });
});

