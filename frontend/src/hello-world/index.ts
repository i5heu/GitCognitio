export class Hello extends HTMLElement {
  constructor() {
    super();
    this.attachShadow({ mode: "open" });

    console.log("Hello world!!!!!");

    fetch("./hello-world.html")
      .then((response) => response.text())
      .then((html) => {
        const template = document.createElement("div");
        template.innerHTML = html;
        console.log("--->, ", html);
        (this as any).shadowRoot.appendChild(template);
      });
  }
}
