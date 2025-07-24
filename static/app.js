document.addEventListener("DOMContentLoaded", () => {
  let count = 0;
  const form = document.getElementById("promptForm");
  const promptInput = document.getElementById("prompt");
  const errorBox = document.getElementById("errorBox");
  const conversation = document.querySelector("div.conversation-container");
  const conversationWindow = document.querySelector("div.conversation-window");
  const scrollFromTopOffset = 12; //12px

  function pasteTextInPrompt(text) {
    promptInput.value = text;
    promptInput.select();
  }

  function addToConversation(type, text) {
    const div = document.createElement("div");
    const isPrompt = type === "prompt";
    const classType = isPrompt ? "prompt" : "reply";
    div.classList.add("conversation", classType);
    div.textContent = text;
    conversation.append(div);

    if (isPrompt) {
      count += 1;
      div.id = "prompt-" + count;
      const newDiv = document.getElementById(div.id);
      if (newDiv && conversationWindow) {
        conversationWindow.scrollTop =
          newDiv.offsetTop - conversationWindow.offsetTop - scrollFromTopOffset;
      }
      div.onclick = () => pasteTextInPrompt(div.textContent);
    }
  }

  form.addEventListener("submit", async (e) => {
    e.preventDefault();

    if (promptInput.value.trim().length < 1) return;
    errorBox.textContent = "";

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
      errorBox.textContent = "ERROR: " + text;
      return;
    }
    addToConversation("prompt", promptInput.value);
    promptInput.value = "";

    const json = await res.json();
    addToConversation("reply", json.response);
    promptInput.disabled = false;
    promptInput.focus();
  });

  promptInput.focus();
});
