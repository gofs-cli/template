// @ts-check
const inflight = new WeakMap();

/**
 * Fetch a server-rendered HTML fragment on behalf of `el`.
 *
 * Aborts the previous request made for `el`, sets aria-busy on it while the
 * request is pending, sends Accept: text/html so the handler returns a
 * fragment rather than a page, and adds the CSRF token.
 *
 * Resolves to the fragment, or null if this request was superseded by a
 * newer one for the same element. A 422 (Unprocessable Content) response is
 * a fragment too, by convention a form rendered again with its errors, so it
 * is returned like a 2xx. Any other non-2xx response throws an Error whose
 * message is the server's response text.
 *
 * @param {Element} el
 * @param {string | URL} url
 * @param {RequestInit} [init]
 * @returns {Promise<string | null>}
 */
export async function fragment(el, url, init = {}) {
  inflight.get(el)?.abort();
  const ctrl = new AbortController();
  inflight.set(el, ctrl);
  el.setAttribute("aria-busy", "true");
  try {
    const res = await fetch(url, {
      ...init,
      headers: { Accept: "text/html", "X-CSRF-Token": csrfToken(), ...init.headers },
      signal: ctrl.signal,
    });
    if (!res.ok && res.status !== 422) throw new Error((await res.text()) || res.statusText);
    return await res.text();
  } catch (e) {
    if (e instanceof DOMException && e.name === "AbortError") return null;
    throw e;
  } finally {
    if (inflight.get(el) === ctrl) {
      inflight.delete(el);
      el.removeAttribute("aria-busy");
    }
  }
}

/**
 * Submit a form the way the browser would, but ask for a fragment.
 * `submitter` is the button that submitted the form, if any, so that its
 * name and value are included (as the browser does).
 *
 * @param {Element} el
 * @param {HTMLFormElement} form
 * @param {HTMLElement | null} [submitter]
 * @returns {Promise<string | null>}
 */
export function submit(el, form, submitter) {
  // Read the attributes rather than form.action / form.method: a control named
  // "action" or "method" (common for buttons) shadows those properties.
  const button = submitter instanceof HTMLButtonElement ? submitter : undefined;
  const method = (button?.getAttribute("formmethod") ?? form.getAttribute("method") ?? "GET").toUpperCase();
  const url = new URL(button?.getAttribute("formaction") ?? form.getAttribute("action") ?? "", location.href);
  const data = new FormData(form, button);
  if (method === "GET") {
    url.search = new URLSearchParams(/** @type {any} */(data)).toString();
    return fragment(el, url);
  }
  // Post url-encoded, as a plain form does, so handlers read both the same
  // way (r.PostFormValue / r.ParseForm). Forms that upload files declare
  // enctype="multipart/form-data" and are sent as FormData.
  const multipart = form.getAttribute("enctype") === "multipart/form-data";
  return fragment(el, url, { method, body: multipart ? data : new URLSearchParams(/** @type {any} */(data)) });
}

/** Reads the CSRF token the layout renders into <meta name="csrf-token">. */
function csrfToken() {
  return document.querySelector('meta[name="csrf-token"]')?.getAttribute("content") ?? "";
}

/**
 * Delay calls to `fn` until `ms` have passed without another call.
 * @template {(...args: any[]) => void} F
 * @param {F} fn
 * @param {number} ms
 * @returns {F}
 */
export function debounce(fn, ms) {
  /** @type {ReturnType<typeof setTimeout> | undefined} */
  let t;
  return /** @type {F} */ (
    (...args) => {
      clearTimeout(t);
      t = setTimeout(() => fn(...args), ms);
    }
  );
}

/**
 * Show an error message in `el`'s own output area: the first element matching
 * `selector`, or a <p class="error"> appended to `el`.
 * @param {Element} el
 * @param {unknown} err
 * @param {string} [selector]
 */
export function showError(el, err, selector = ".error") {
  const message = err instanceof Error ? err.message : String(err);
  let out = el.querySelector(selector);
  if (!out) {
    out = document.createElement("p");
    out.className = "error";
    el.append(out);
  }
  out.textContent = message; // textContent: never build HTML from data on the client
}
