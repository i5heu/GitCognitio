import { MessageWebSocket } from "../communications";
import { Store } from "../store/store";
import { ChatItem } from "./chat-item";

export class ChatBody extends HTMLElement {
  public ID: number = 1;

  constructor() {
    super();
    this.attachShadow({ mode: "open" });

    customElements.define("chat-item", ChatItem);

    fetch("./chat-body.html")
      .then((response) => response.text())
      .then((html) => {
        const fragment = document.createRange().createContextualFragment(html);
        (this as any).shadowRoot.appendChild(fragment);
        this.setup();
      });
  }

  async setup() {
    this.sendKeyStrokeListener();
    this.render();
  }

  async sendKeyStrokeListener() {
    console.log("this.fragment", this.shadowRoot);

    this.shadowRoot
      .querySelector("#chatInput")
      .addEventListener("keyup", async (event: KeyboardEvent) => {
        let start = Date.now();
        //get key pressed
        this.render();
      });

    new interaction(this.shadowRoot).setInputEventListener();
  }

  render() {
    console.log("render");
    const chatHistory = this.shadowRoot.querySelector("#chatHistory");
    chatHistory.innerHTML = "";
  }
}

class interaction {
  shadowRoot: ShadowRoot;
  store: Store;
  constructor(shadowRoot: ShadowRoot) {
    this.shadowRoot = shadowRoot;
  }

  setInputEventListener() {
    this.shadowRoot
      .querySelector("#chatInput")
      .addEventListener("keyup", async (event: KeyboardEvent) => {
        if (event.key === "Enter") {
          const message = (
            this.shadowRoot.querySelector("#chatInput") as HTMLInputElement
          ).value;
          this.store.addStorageItem("userInput", message);
        }
      });
  }
}
