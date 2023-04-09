export class ChatItem extends HTMLElement {
  constructor() {
    super();
    this.attachShadow({ mode: "open" });

    console.log("ChatItem");

    fetch("./chat-item.html")
      .then((response) => response.text())
      .then((html) => {
        const fragment = document.createRange().createContextualFragment(html);
        (this as any).shadowRoot.appendChild(fragment);
      });
  }
}
