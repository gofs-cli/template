// @ts-check
import { fragment, submit, showError } from "gofs/request.js";

class EditInline extends HTMLElement {
  connectedCallback() {
    this.addEventListener("click", async (e) => {
      const link = e.target instanceof Element ? e.target.closest("a") : null;
      if (!link) return;
      e.preventDefault();
      this.#swapRow(fragment(this, link.href));
    });
    this.addEventListener("submit", (e) => {
      if (!(e.target instanceof HTMLFormElement)) return;
      e.preventDefault();
      this.#swapRow(submit(this, e.target, e.submitter));
    });
  }

  /** @param {Promise<string | null>} request */
  async #swapRow(request) {
    const row = this.closest("tr");
    if (!row) return;
    try {
      const html = await request;
      if (html !== null) row.outerHTML = html;
    } catch (err) {
      showError(this, err);
    }
  }
}

customElements.define("edit-inline", EditInline);
