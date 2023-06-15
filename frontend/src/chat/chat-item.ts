import * as MarkdownIt from "markdown-it";
import { Communications } from "../communications";

export class ChatItem extends HTMLElement {
  message: any;
  coms: Communications;
  id: string;

  constructor() {
    super();
    this.attachShadow({ mode: "open" });
    console.log("ChatItem");
  }

  async init(coms: Communications) {
    await this.loadTemplate();
    this.attachDeleteBtnEvent();
    this.coms = coms;
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

    // Set the ID
    this.id = message.id;
    this.message = message;

    this.coms.Router.registerThread(this.id, (message) =>
      this.internalRouter(message)
    );

    // Convert Markdown to HTML
    const html = md.render(message.data);
    this.shadowRoot.querySelector(".content").innerHTML = html;
  }

  internalRouter(message: any) {
    console.log("internalRouter", message);
    if (message.type === "thread-delete") {
      this.delete();
    }
  }

  delete() {
    console.log("delete", this.id);
    this.remove();
  }

  sendDeleteToServer() {
    console.log("sendDelteToServer", this.id);
    this.coms.send(this.id, "delete", "", this.message.path);
  }

  attachDeleteBtnEvent() {
    console.log("attachDeleteBtnEvent");
    this.shadowRoot
      .querySelector(".delete-btn")
      .addEventListener("click", () => this.sendDeleteToServer());
  }
}
