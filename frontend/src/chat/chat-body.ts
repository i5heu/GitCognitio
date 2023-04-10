import { MessageWebSocket } from "../communications";
import { Store } from "../store/store";
import { ChatItem } from "./chat-item";

export class ChatBody extends HTMLElement {
  private socket: MessageWebSocket;
  private store: Store;

  constructor(socket: MessageWebSocket, store: Store) {
    super();
    this.socket = socket;
    this.store = store;
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
    await this.socket.connect();
    this.socket.setInputHandler((message: any) => {
      console.log("Input Handler", message);
      (this.shadowRoot.querySelector("#chatInput") as HTMLInputElement).value +=
        message;
    });
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
        const vov = await this.socket.send("input", event.key);
        console.log("Time", Date.now() - start, "VOV", vov);
        this.render();
      });

    new interaction(this.shadowRoot, this.store).setInputEventListener();
  }

  render() {
    console.log("render");
    const messages = this.store.get("userInput").items;
    const chatHistory = this.shadowRoot.querySelector("#chatHistory");
    chatHistory.innerHTML = "";
    messages.forEach((message: any) => {
      const chatItem = new ChatItem(message);
      chatHistory.appendChild(chatItem);
    });
  }
}

class interaction {
  shadowRoot: ShadowRoot;
  store: Store;
  constructor(shadowRoot: ShadowRoot, store: Store) {
    this.shadowRoot = shadowRoot;
    this.store = store;
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
