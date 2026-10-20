// @ts-check
import { submit, showError } from "gofs/request.js";

class RowDelete extends HTMLElement {
  connectedCallback() {
    this.addEventListener("submit", async (e) => {
      if (!(e.target instanceof HTMLFormElement)) return;
      e.preventDefault();
      if (!(await this.#confirmed())) return;
      try {
        const html = await submit(this, e.target, e.submitter);
        if (html !== null) this.closest("tr")?.remove();
      } catch (err) {
        showError(this, err);
      }
    });
  }

  async #confirmed() {
    const dialog = document.getElementById(this.getAttribute("confirm") ?? "");
    if (!(dialog instanceof HTMLDialogElement)) return true;
    dialog.returnValue = "";
    dialog.showModal();
    await new Promise((resolve) => dialog.addEventListener("close", resolve, { once: true }));
    return dialog.returnValue === "confirm";
  }
}

customElements.define("row-delete", RowDelete);
