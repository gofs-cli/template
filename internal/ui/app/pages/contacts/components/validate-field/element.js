// @ts-check
import { fragment, debounce } from "gofs/request.js";

class ValidateField extends HTMLElement {
  connectedCallback() {
    const input = this.querySelector("input, select, textarea");
    const output = this.querySelector("output");
    const src = this.getAttribute("src");
    if (!(input instanceof HTMLInputElement) || !output || !src) return;

    const validate = async () => {
      const body = new URLSearchParams({ [input.name]: input.value });
      try {
        const html = await fragment(this, src, { method: "POST", body });
        if (html !== null) output.textContent = html.trim();
      } catch (e) {
        output.textContent = e instanceof Error ? e.message : String(e);
      }
    };

    input.addEventListener("input", debounce(validate, 300));
    input.addEventListener("change", validate);
  }
}

customElements.define("validate-field", ValidateField);
