document.querySelectorAll("[data-copy]").forEach((button) => {
  const originalLabel = button.textContent || "Copy";
  const resetDelay = 2000;

  button.addEventListener("click", async () => {
    const code = button.getAttribute("data-code") || "";
    try {
      await navigator.clipboard.writeText(code);
      button.textContent = "Copied";
      button.classList.add("text-success");
      setTimeout(() => {
        button.textContent = originalLabel;
        button.classList.remove("text-success");
      }, resetDelay);
    } catch {
      button.textContent = "Failed";
      setTimeout(() => {
        button.textContent = originalLabel;
      }, resetDelay);
    }
  });
});
