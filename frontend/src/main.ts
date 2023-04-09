import { LoginPage } from "./login-page";
import { ChatBody } from "./chat/chat-body";
import "./main.scss";

import { MessageWebSocket } from "./communications";
import { Store } from "./store/store";

//get ip address of websocket server
const socketManager = new MessageWebSocket(
  "ws://" + window.location.hostname + ":8081/"
);

//create storage
const store = new Store();

customElements.define("login-page", LoginPage);
customElements.define("chat-body", ChatBody);

const renderTarget = document.getElementById("root");
function render() {
  renderTarget.innerHTML = "";

  // if (localStorage.getItem("token")) {
  //   renderTarget.appendChild(new ChatBody());
  // } else {
  //   renderTarget.appendChild(new LoginPage());
  // }

  renderTarget.appendChild(new ChatBody(socketManager, store));
}

render();
