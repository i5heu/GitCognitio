import { LoginPage } from "./login-page";
import { ChatBody } from "./chat/chat-body";
import "./main.scss";

import { MessageWebSocket } from "./communications";
import { render } from "./render";

customElements.define("login-page", LoginPage);
customElements.define("chat-body", ChatBody);

const renderTarget = document.getElementById("root");
const vov = renderTarget.appendChild(new ChatBody());

//get ip address of websocket server
const socketManager = new MessageWebSocket(
  "ws://" + window.location.hostname + ":8081/",
  vov
);

socketManager.connect().then(() => {
  console.log("connected");
  socketManager.send("login", { username: "bob" });
  socketManager.send("login", { username: "bob" });
  socketManager.send("login", { username: "bob" });
});
