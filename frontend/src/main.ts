import { ChatBody } from "./chat/chat-body";
import "./main.scss";

import { Communications } from "./communications";
import { Ping } from "./ping";

const comms = new Communications("wss://ws.miau.email");
new Ping(comms).startPing(500);

customElements.define("chat-body", ChatBody);

const renderTarget = document.getElementById("root");
renderTarget.appendChild(new ChatBody(comms));

document.querySelector("#close-qrscanner").addEventListener("click", () => {
  (document.querySelector("#modal-qrscanner") as HTMLDivElement).style.display =
    "none";
});
