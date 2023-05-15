import { ChatBody } from "./chat/chat-body";
import { ChatItem } from "./chat/chat-item";
import { MessageWebSocket } from "./communications";

export class render {
  cb: ChatBody[] = [];
  ci: ChatItem[] = [];

  private socket: MessageWebSocket;

  constructor() {
    const renderTarget = document.getElementById("root");

    const vov = renderTarget.appendChild(new ChatBody());
    this.cb[vov.ID] = vov;
  }

  renderChatItem(message: any) {
    console.log("renderChatItem", message);
  }
}
