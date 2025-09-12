document.getElementById("screenshotBtn").addEventListener("click", () => {
  chrome.tabs.captureVisibleTab(null, { format: "png" }, (dataUrl) => {
    if (chrome.runtime.lastError) {
      alert("截图失败: " + chrome.runtime.lastError.message);
      return;
    }

    // 在新窗口里显示截图
    const imgWindow = window.open();
    imgWindow.document.write(`<img src="${dataUrl}" style="max-width:100%"/>`);
  });
});
