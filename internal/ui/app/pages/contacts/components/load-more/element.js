// @ts-check
import { fragment, showError } from "gofs/request.js";

class LoadMore extends HTMLElement {
  connectedCallback() {
    this.addEventListener("click", async (e) => {
      const link = e.target instanceof Element ? e.target.closest("a") : null;
      const target = document.getElementById(this.getAttribute("target") ?? "");
      if (!link || !target) return;
      e.preventDefault();
      try {
        const html = await fragment(this, link.href);
        if (html === null) return;
        const tpl = document.createElement("template");
        tpl.innerHTML = html;
        target.append(...tpl.content.querySelectorAll("tr"));
        this.replaceChildren(tpl.content);
      } catch (err) {
        showError(this, err);
      }
    });
  }
}

customElements.define("load-more", LoadMore);
