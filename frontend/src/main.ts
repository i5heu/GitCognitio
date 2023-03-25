import { LoginPage } from "./login-page";
import { ChatBody } from "./chat";
import "./main.scss";

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
