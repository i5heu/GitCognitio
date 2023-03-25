import { MessageWebSocket } from "../communications";

export class ChatBody extends HTMLElement {
  private socket: MessageWebSocket;

  constructor(socket: MessageWebSocket) {
    super();
    this.socket = socket;
    this.attachShadow({ mode: "open" });

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
}
