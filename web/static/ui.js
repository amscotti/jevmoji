/* Small conveniences around the model's read-only reactions: try a message
 * and copy the resulting emojis. HTMX handles all requests and rendering. */

(() => {
  const message = document.getElementById("message");
  const examples = document.getElementById("examples");
  const results = document.getElementById("results");
  if (!message || !examples || !results) return;

  const updateExamples = () => {
    for (const button of examples.querySelectorAll("[data-example-message]")) {
      button.classList.toggle(
        "is-current",
        button.dataset.exampleMessage === message.value.trim(),
      );
    }
  };

  examples.hidden = false;
  examples.addEventListener("click", (event) => {
    const button = event.target.closest("[data-example-message]");
    if (!button) return;
    message.value = button.dataset.exampleMessage;
    message.dispatchEvent(new Event("input", { bubbles: true }));
  });
  message.addEventListener("input", updateExamples);
  updateExamples();

  const enhanceResults = () => {
    for (const button of results.querySelectorAll("[data-copy-emojis]")) {
      button.hidden = false;
    }
    results.setAttribute(
      "aria-busy",
      String(results.classList.contains("htmx-request")),
    );
  };
  new MutationObserver(enhanceResults).observe(results, {
    childList: true,
    attributes: true,
    attributeFilter: ["class"],
  });
  enhanceResults();

  results.addEventListener("click", async (event) => {
    const button = event.target.closest("[data-copy-emojis]");
    if (!button || button.disabled) return;
    const label = button.querySelector("[data-copy-label]");
    const status = button.parentElement.querySelector("[data-copy-status]");
    button.disabled = true;
    try {
      await navigator.clipboard.writeText(button.dataset.copyEmojis.trim());
      label.textContent = "Copied!";
      status.textContent = "Emojis copied to your clipboard.";
    } catch {
      label.textContent = "Couldn't copy";
      status.textContent = "Clipboard access is unavailable. Select the emojis and copy them.";
    } finally {
      button.disabled = false;
      setTimeout(() => {
        if (button.isConnected) label.textContent = "Copy emojis";
      }, 2000);
    }
  });
})();
