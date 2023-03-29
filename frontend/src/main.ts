import { LoginPage } from "./login-page";
import { ChatBody } from "./chat";
import "./main.scss";

import { MessageWebSocket } from "./communications";

//get ip address of websocket server
const socket = new MessageWebSocket(
  "ws://" + window.location.hostname + ":8081/"
);

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

  renderTarget.appendChild(new ChatBody(socket));
}

render();
