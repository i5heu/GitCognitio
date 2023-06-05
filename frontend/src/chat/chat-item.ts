import * as MarkdownIt from "markdown-it";

export class ChatItem extends HTMLElement {
  deleteCall: () => void;
  constructor(deleteCall = () => {}) {
    super();
    this.attachShadow({ mode: "open" });
    this.deleteCall = deleteCall;
    console.log("ChatItem");
  }

  async init() {
    await this.loadTemplate();
    await this.attachDeleteBtnEvent();
  }

  private async loadTemplate() {
    const response = await fetch("./chat-item.html");
    const html = await response.text();
    const fragment = document.createRange().createContextualFragment(html);
    (this as any).shadowRoot.appendChild(fragment);
  }

  addContent(message: any) {
    console.log("addContent", message);
    const md = new MarkdownIt();

    // Convert Markdown to HTML
    const html = md.render(message.data);
    this.shadowRoot.querySelector(".content").innerHTML = html;
  }

  attachDeleteBtnEvent() {
    this.shadowRoot
      .querySelector(".delete-btn")
      .addEventListener("click", this.deleteCall);
  }
}
