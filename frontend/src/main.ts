import { LoginPage } from "./login-page";
import { ChatBody } from "./chat/chat-body";
import "./main.scss";

import { MessageWebSocket } from "./communications";
import { render } from "./render";

customElements.define("login-page", LoginPage);
customElements.define("chat-body", ChatBody);

const renderInstance = new render();

//get ip address of websocket server
const socketManager = new MessageWebSocket(
  "ws://" + window.location.hostname + ":8081/",
  renderInstance
);

socketManager.connect();

(window as any).bob = socketManager.send;
