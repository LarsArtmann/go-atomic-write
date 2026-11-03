(function () {
  const resetDelay = 2000;

  if (document.documentElement.dataset.copyBound) {
    return;
  }
  document.documentElement.dataset.copyBound = "true";

  function announce(message) {
    const region = document.querySelector("#copy-status");
    if (region) {
      region.textContent = message;
    }
  }

  function restore(button) {
    button.textContent = button.dataset.copyLabel || "Copy";
    button.classList.remove("text-success");
  }

  document.addEventListener("click", (event) => {
    const { target } = event;
    if (!(target instanceof Element)) {
      return;
    }
    const button = target.closest("[data-copy]");
    if (!button) {
      return;
    }

    if (!button.dataset.copyLabel) {
      button.dataset.copyLabel = button.textContent || "Copy";
    }

    const code = button.dataset.code || "";

    navigator.clipboard.writeText(code).then(
      () => {
        button.textContent = "Copied";
        button.classList.add("text-success");
        announce("Copied to clipboard");
        setTimeout(() => {
          restore(button);
        }, resetDelay);
      },
      () => {
        button.textContent = "Failed";
        announce("Copy failed. Select the code and copy it manually.");
        setTimeout(() => {
          restore(button);
        }, resetDelay);
      },
    );
  });
})();
