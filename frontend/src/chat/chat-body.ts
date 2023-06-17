import { Communications } from "../communications";
import { ChatItem } from "./chat-item";
import { InstanceIdentifier } from "../helper/instanceIdentifier";
import { QrLoginScanner } from "./qrLoginScanner";

export class ChatBody extends HTMLElement {
  public ID: number = 1;
  private comms: Communications;

  constructor(comms: Communications) {
    super();
    this.attachShadow({ mode: "open" });
    this.comms = comms;

    customElements.define("chat-item", ChatItem);
    customElements.define("qr-scanner-item", QrLoginScanner);

    fetch("./chat-body.html")
      .then((response) => response.text())
      .then((html) => {
        const fragment = document.createRange().createContextualFragment(html);
        this.shadowRoot.appendChild(fragment);
        this.setup();
      });
  }

  setup() {
    this.setInputEventListener();
    this.render();
    this.comms.Router.register("message", this.messageHandler.bind(this));
    this.comms.Router.register(
      "qrLoginRequest",
      this.qrLoginHandler.bind(this)
    );
    this.comms.Router.register("typing", this.typingHandler.bind(this));
  }

  async typingHandler(message: any) {
    // ignore own typing
    if (message.id === InstanceIdentifier.getInstanceIdentifier()) {
      return;
    }

    this.adjustTextArea();
    this.chatInput.value = message.data;
  }

  async messageHandler(message: any) {
    const chatItem = new ChatItem();
    await chatItem.init(this.comms);
    chatItem.addContent(message);
    this.chatHistory.appendChild(chatItem);

    //scroll to bottom
    this.chatHistory.scrollTop = this.chatHistory.scrollHeight;
  }

  async qrLoginHandler(message: any) {
    const qrItem = new QrLoginScanner();
    await qrItem.init(this.comms);
    this.chatHistory.appendChild(qrItem);

    //scroll to bottom
    this.chatHistory.scrollTop = this.chatHistory.scrollHeight;
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
    this.chatInput.addEventListener(
      "keypress",
      async (event: KeyboardEvent) => {
        this.adjustTextArea();

        if (
          event.key === "Enter" &&
          !event.shiftKey &&
          this.chatInput.value.trim() == ""
        ) {
          this.chatInput.value = "";
          this.adjustTextArea();
          event.preventDefault();
          return;
        }

        if (
          event.key === "Enter" &&
          !event.shiftKey &&
          this.chatInput.value != ""
        ) {
          const message = this.chatInput.value;
          const id = this.generateId();
          this.comms.send(id, "message", message);

          // clear input after keypress
          setTimeout(() => {
            this.chatInput.value = "";
            this.adjustTextArea();
          }, 0);
        } else {
          const id = InstanceIdentifier.getInstanceIdentifier();
          this.comms.send(id, "typing", this.chatInput.value);
        }
      }
    );
  }

  adjustTextArea() {
    this.chatInput.style.height = ""; // Reset the height to recalculate the scroll height
    this.chatInput.style.height = this.chatInputContentHeight + "px";
    this.chatHistory.scrollTop = this.chatHistory.scrollHeight;
  }

  get chatHistory(): HTMLDivElement {
    return this.shadowRoot.querySelector("#chatHistory");
  }

  get chatInput(): HTMLTextAreaElement {
    return this.shadowRoot.querySelector("#chatInput");
  }

  get chatInputContentHeight() {
    var style = window.getComputedStyle(this.chatInput);
    const paddingTop = parseFloat(style.paddingTop);
    const paddingBottom = parseFloat(style.paddingBottom);
    const totalPaddingHeight = paddingTop + paddingBottom;

    return totalPaddingHeight / 2 + this.chatInput.scrollHeight;
  }

  private generateId(): string {
    // Generate a unique identifier using any desired method
    // For simplicity, let's use a random string with 10 characters
    const characters =
      "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
    let id = "";

    for (let i = 0; i < 10; i++) {
      id += characters.charAt(Math.floor(Math.random() * characters.length));
    }

    return id;
  }
}
