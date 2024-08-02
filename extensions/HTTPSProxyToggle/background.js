const OFF_TEXT = 'OFF';
const ON_TEXT = 'ON';

const PROXY_CONFIG = {
    mode: "fixed_servers",
    rules: {
        singleProxy: {
            scheme: "https",
            host: "127.0.0.1",
            port: 8080
        }
    }
};

chrome.runtime.onInstalled.addListener(() => {
    chrome.storage.local.set({ proxyEnabled: false }, () => {
        updateProxy(false);
    });
});

chrome.action.onClicked.addListener(() => {
    chrome.storage.local.get("proxyEnabled", (data) => {
        const newStatus = !data.proxyEnabled;
        chrome.storage.local.set({ proxyEnabled: newStatus }, () => {
            updateProxy(newStatus);
        });
    });
});

function updateProxy(enabled) {
    if (enabled) {
        chrome.proxy.settings.set({ value: PROXY_CONFIG, scope: "regular" }, () => {
            chrome.action.setBadgeText({ text: ON_TEXT });
        });
    } else {
        chrome.proxy.settings.clear({ scope: "regular" }, () => {
            chrome.action.setBadgeText({ text: OFF_TEXT });
        });
    }
}

chrome.storage.local.get("proxyEnabled", (data) => {
    updateProxy(data.proxyEnabled);
});
