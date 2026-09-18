(function () {
  var resetDelay = 2000;

  if (document.documentElement.dataset.copyBound) return;
  document.documentElement.dataset.copyBound = "true";

  function announce(message) {
    var region = document.getElementById("copy-status");
    if (region) region.textContent = message;
  }

  function restore(button) {
    button.textContent = button.dataset.copyLabel || "Copy";
    button.classList.remove("text-success");
  }

  document.addEventListener("click", function (event) {
    var target = event.target;
    if (!(target instanceof Element)) return;
    var button = target.closest("[data-copy]");
    if (!button) return;

    if (!button.dataset.copyLabel) {
      button.dataset.copyLabel = button.textContent || "Copy";
    }

    var code = button.getAttribute("data-code") || "";

    navigator.clipboard.writeText(code).then(
      function () {
        button.textContent = "Copied";
        button.classList.add("text-success");
        announce("Copied to clipboard");
        setTimeout(function () {
          restore(button);
        }, resetDelay);
      },
      function () {
        button.textContent = "Failed";
        announce("Copy failed. Select the code and copy it manually.");
        setTimeout(function () {
          restore(button);
        }, resetDelay);
      },
    );
  });
})();
