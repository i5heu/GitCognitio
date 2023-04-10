import { StorageItem } from "../store/store";

export class ChatItem extends HTMLElement {
  deleteCall: () => void;
  storageItem: StorageItem;
  constructor(storageItem: StorageItem, deleteCall = () => {}) {
    super();
    this.attachShadow({ mode: "open" });
    this.deleteCall = deleteCall;
    this.storageItem = storageItem;

    console.log("ChatItem");
    this.init();
  }

  async init() {
    await this.loadTemplate();
    await this.attachDeleteBtnEvent();
    this.addContent(
      this.storageItem.content[this.storageItem.content.length - 1].content
    );
  }

  async loadTemplate() {
    const response = await fetch("./chat-item.html");
    const html = await response.text();
    const fragment = document.createRange().createContextualFragment(html);
    (this as any).shadowRoot.appendChild(fragment);
  }

  addContent(content: string) {
    this.shadowRoot.querySelector(".content").innerHTML = content;
  }

  attachDeleteBtnEvent() {
    this.shadowRoot
      .querySelector(".delete-btn")
      .addEventListener("click", this.deleteCall);
  }
}
