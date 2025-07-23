document.addEventListener("DOMContentLoaded", () => {
  const form = document.getElementById("promptForm");
  const promptInput = document.getElementById("prompt");
  const responseBox = document.getElementById("responseBox");
  const errorBox = document.getElementById("errorBox");
  const queryText = document.getElementById("query");

  console.log("LOADED!");

  form.addEventListener("submit", async (e) => {
    e.preventDefault();

    if (promptInput.value.trim().length < 1) return;

    errorBox.textContent = "";
    responseBox.textContent = "responding...";

    responseBox.disabled = true;
    promptInput.disabled = true;

    const res = await fetch("/generate", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },

      body: JSON.stringify({
        prompt: promptInput.value,
      }),
    });

    if (!res.ok) {
      const text = await res.text();
      errorBox.textContent = text;
      return;
    }

    const json = await res.json();
    responseBox.disabled = false;
    responseBox.textContent = json.response;
    queryText.innerText = promptInput.value;
    promptInput.value = "";
    promptInput.disabled = false;
    promptInput.focus();
  });
});
