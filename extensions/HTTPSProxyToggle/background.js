const OFF_TEXT = "OFF";
const ON_TEXT = "ON";

const PROXY_CONFIG = {
    mode: "fixed_servers",
    rules: {
        proxyForHttps: {
            scheme: "http",
            host: "127.0.0.1",
            port: 8080
        }
    }
};

function initializeProxyState() {
    chrome.storage.local.get("proxyEnabled", (data) => {
        updateProxy(data.proxyEnabled);
    });
}

chrome.runtime.onInstalled.addListener(initializeProxyState);
chrome.runtime.onStartup.addListener(initializeProxyState);

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
        return;
    }
    chrome.proxy.settings.clear({ scope: "regular" }, () => {
        chrome.action.setBadgeText({ text: OFF_TEXT });
    });
}


