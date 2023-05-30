import { Communications } from "../communications";
import { ChatItem } from "./chat-item";
import { InstanceIdentifier } from "../helper/instanceIdentifier";

export class ChatBody extends HTMLElement {
  public ID: number = 1;
  private comms: Communications;

  constructor(comms: Communications) {
    super();
    this.attachShadow({ mode: "open" });
    this.comms = comms;

    customElements.define("chat-item", ChatItem);

    fetch("./chat-body.html")
      .then((response) => response.text())
      .then((html) => {
        const fragment = document.createRange().createContextualFragment(html);
        this.shadowRoot.appendChild(fragment);
        this.setup();
      });
  }

  setup() {
    this.sendKeyStrokeListener();
    this.render();
    this.comms.Router.register("message", this.messageHandler.bind(this));
    this.comms.Router.register("typing", this.typingHandler.bind(this));
  }

  async typingHandler(message: any) {
    if (message.id === InstanceIdentifier.getInstanceIdentifier()) {
      return;
    }

    (this.shadowRoot.querySelector("#chatInput") as HTMLInputElement).value =
      message.data;
  }

  async messageHandler(message: any) {
    console.log("this", this.shadowRoot);
    console.log("messageHandler", message);
    const chatHistory = this.shadowRoot.querySelector("#chatHistory");
    const chatItem = new ChatItem();
    await chatItem.init();
    chatItem.addContent(message);
    chatHistory.appendChild(chatItem);

    //scroll to bottom
    chatHistory.scrollTop = chatHistory.scrollHeight;
  }

  async sendKeyStrokeListener() {
    console.log("this.fragment", this.shadowRoot);
    this.setInputEventListener();
  }

  render() {
    console.log("render");
    const chatHistory = this.shadowRoot.querySelector("#chatHistory");
    chatHistory.innerHTML = "";
  }

  destroy() {
    console.log("destroy");
  }

  setInputEventListener() {
    this.shadowRoot
      .querySelector("#chatInput")
      .addEventListener("keyup", async (event: KeyboardEvent) => {
        if (event.key === "Enter") {
          const message = (
            this.shadowRoot.querySelector("#chatInput") as HTMLInputElement
          ).value;

          this.comms.send("1", "message", message);
        } else {
          const id = InstanceIdentifier.getInstanceIdentifier();
          this.comms.send(
            id,
            "typing",
            (this.shadowRoot.querySelector("#chatInput") as HTMLInputElement)
              .value
          );
        }
      });
  }
}
