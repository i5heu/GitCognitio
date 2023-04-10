import { MessageWebSocket } from "../communications";
import { Store } from "../store/store";
import { ChatItem } from "./chat-item";

export class ChatBody extends HTMLElement {
  private socket: MessageWebSocket;

  constructor(socket: MessageWebSocket, store: Store) {
    super();
    this.socket = socket;
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

    this.shadowRoot
      .getElementById("chatHistory")
      .appendChild(new ChatItem(() => this.render()));

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
      });
  }

  render() {
    this.shadowRoot
      .getElementById("chatHistory")
      .appendChild(new ChatItem(() => this.render()));
  }
}
