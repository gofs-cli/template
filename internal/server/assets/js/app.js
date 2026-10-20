// @ts-check
const loaded = new Set();

/** Import the definition of every custom element under root that is not loaded yet.
 * @param {ParentNode} root */
function loadComponents(root = document) {
  for (const el of root.querySelectorAll("*")) {
    const tag = el.localName;
    if (!tag.includes("-") || loaded.has(tag) || customElements.get(tag)) continue;
    loaded.add(tag);
    import(`/assets/components/${tag}/element.js`).catch((e) => console.error(`component <${tag}>:`, e));
  }
}

loadComponents();

// Fragments swapped in by components may contain new custom elements.
new MutationObserver((records) => {
  for (const r of records) {
    for (const n of r.addedNodes) if (n instanceof Element) loadComponents(n.parentNode ?? n);
  }
}).observe(document.body, { childList: true, subtree: true });

if (!("command" in HTMLButtonElement.prototype)) {
  document.addEventListener("click", (e) => {
    const button = e.target instanceof Element ? e.target.closest("button[commandfor]") : null;
    if (!button) return;
    const target = document.getElementById(button.getAttribute("commandfor") ?? "");
    if (!(target instanceof HTMLDialogElement)) return;
    switch (button.getAttribute("command")) {
      case "show-modal":
        target.showModal();
        break;
      case "close":
        target.close();
        break;
    }
  });
}
