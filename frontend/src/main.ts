import { ChatBody } from "./chat/chat-body";
import "./main.scss";

import { Communications } from "./communications";

const comms = new Communications("ws://" + window.location.hostname + ":8081/");

customElements.define("chat-body", ChatBody);

const renderTarget = document.getElementById("root");
renderTarget.appendChild(new ChatBody(comms));

comms.connect().then(() => {
  console.log("connected");
  comms.send("1", "login", "bob");
  comms.send("1", "login", "bob");
  comms.send("1", "login", "bob");
});
