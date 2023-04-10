export class ChatItem extends HTMLElement {
  deleteCall: () => void;
  constructor(deleteCall = () => {}) {
    super();
    this.attachShadow({ mode: "open" });
    this.deleteCall = deleteCall;

    console.log("ChatItem");
    this.init();
  }

  async init() {
    await this.loadTemplate();
    await this.attachDeleteBtnEvent();
  }

  async loadTemplate() {
    const response = await fetch("./chat-item.html");
    const html = await response.text();
    const fragment = document.createRange().createContextualFragment(html);
    (this as any).shadowRoot.appendChild(fragment);
  }

  attachDeleteBtnEvent() {
    this.shadowRoot
      .querySelector(".delete-btn")
      .addEventListener("click", this.deleteCall);
  }
}
