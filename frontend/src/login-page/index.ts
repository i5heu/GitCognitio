export class LoginPage extends HTMLElement {
  constructor() {
    super();
    this.attachShadow({ mode: "open" });

    fetch("./login-page.html")
      .then((response) => response.text())
      .then((html) => {
        const template = document.createElement("div");
        template.innerHTML = html;
        (this as any).shadowRoot.appendChild(template);
      });
  }
}
