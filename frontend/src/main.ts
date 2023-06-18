import { ChatBody } from "./chat/chat-body";
import "./main.scss";

import { Communications } from "./communications";

const comms = new Communications("wss://ws.miau.email");

customElements.define("chat-body", ChatBody);

const renderTarget = document.getElementById("root");
renderTarget.appendChild(new ChatBody(comms));

comms.connect().then(() => {
  console.log("connected");
  comms.send("1", "login", "bob");
  comms.send("1", "login", "bob");
  comms.send("1", "login", "bob");
});

document.querySelector("#close-qrscanner").addEventListener("click", () => {
  (document.querySelector("#modal-qrscanner") as HTMLDivElement).style.display =
    "none";
});
