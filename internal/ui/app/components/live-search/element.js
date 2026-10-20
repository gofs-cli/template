// @ts-check
import { submit, debounce, showError } from "gofs/request.js";

class LiveSearch extends HTMLElement {
  connectedCallback() {
    const form = this.querySelector("form");
    const target = document.getElementById(this.getAttribute("target") ?? "");
    if (!form || !target) return;

    const search = async () => {
      try {
        const html = await submit(this, form);
        if (html === null) return; // superseded by a newer search
        target.innerHTML = html; // a same-origin fragment rendered and escaped by the server
        const params = new URLSearchParams(/** @type {any} */(new FormData(form)));
        for (const [k, v] of [...params]) if (v === "") params.delete(k);
        const query = params.toString();
        history.replaceState(null, "", query ? `?${query}` : location.pathname);
      } catch (e) {
        showError(target, e);
      }
    };

    form.addEventListener("input", debounce(search, 200));
    form.addEventListener("change", (e) => {
      if (e.target instanceof HTMLSelectElement) search();
    });
    form.addEventListener("submit", (e) => {
      e.preventDefault();
      search();
    });
  }
}

customElements.define("live-search", LiveSearch);
