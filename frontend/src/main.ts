import { LoginPage } from "./login-page";
import { ChatBody } from "./chat";
import "./main.scss";

import { MessageWebSocket } from "./communications";

customElements.define("login-page", LoginPage);
customElements.define("chat-body", ChatBody);

const renderTarget = document.getElementById("root");
function render() {
  renderTarget.innerHTML = "";

  if (localStorage.getItem("token")) {
    renderTarget.appendChild(new ChatBody());
  } else {
    renderTarget.appendChild(new LoginPage());
  }
}

render();

const socket = new MessageWebSocket("ws://localhost:8081/");

async function bob() {
  await socket.connect();
  let start = Date.now();
  const vov = await socket.send("Hello, world!");
  console.log("Time", Date.now() - start, "VOV", vov);
}

bob();
