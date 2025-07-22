document.addEventListener("DOMContentLoaded", () => {
  const form = document.getElementById("promptForm");
  const promptInput = document.getElementById("prompt");
  const responseBox = document.getElementById("responseBox");
  const errorBox = document.getElementById("errorBox");

  console.log("LOADED!");

  form.addEventListener("submit", async (e) => {
    e.preventDefault();

    errorBox.textContent = "";
    responseBox.textContent = "responding...";

    responseBox.disabled = true;
    promptInput.disabled = true;

    const res = await fetch("/generate", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },

      // TODO FHA - MOVE TO SERVER !!
      body: JSON.stringify({
        model: "llama3.2",
        prompt: promptInput.value,
        stream: false,
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
    promptInput.value = "";
    promptInput.disabled = false;
    promptInput.focus();
  });
});
