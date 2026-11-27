(function () {
  document.querySelectorAll("[data-copy]").forEach(function (button) {
    button.addEventListener("click", function () {
      var node = document.getElementById(button.getAttribute("data-copy"));
      if (!node) return;
      var text = node.innerText.replace(/\u00a0/g, " ");
      var finish = function () {
        button.textContent = button.getAttribute("data-copied") || "OK";
        window.setTimeout(function () {
          button.textContent = button.getAttribute("data-label") || button.textContent;
        }, 1600);
      };
      var fallback = function () {
        var area = document.createElement("textarea");
        area.value = text;
        document.body.appendChild(area);
        area.select();
        var ok = false;
        try { ok = document.execCommand("copy"); } catch (e) {}
        area.remove();
        if (ok) finish();
      };
      if (navigator.clipboard && navigator.clipboard.writeText) {
        navigator.clipboard.writeText(text).then(finish).catch(fallback);
        return;
      }
      fallback();
    });
  });
})();
