import { render } from "./render";

export class MessageWebSocket {
  private socket: WebSocket | null;
  private messages: Map<
    number,
    [(message: any) => void, (error: Error) => void]
  >;
  private InputHandler: (message: any) => void;

  private renderInstance: render;

  constructor(private url: string, renderInstance: render) {
    this.socket = null;
    this.messages = new Map();
    this.renderInstance = renderInstance;
  }

  public connect(): Promise<void> {
    return new Promise((resolve, reject) => {
      this.socket = new WebSocket(this.url);

      this.socket.addEventListener("open", () => {
        console.log("WebSocket connection established");
        resolve();
      });

      this.socket.addEventListener("error", (event) => {
        reject(event);
      });

      this.socket.addEventListener("close", () => {
        console.log("WebSocket connection closed");
      });

      this.socket.addEventListener("message", (event) => {
        const message = JSON.parse(event.data);
        router.route(this.renderInstance, message.type, message.data);
      });
    });
  }

  public send(type: string, data: any): Promise<any> {
    if (!this.socket || this.socket.readyState !== WebSocket.OPEN) {
      return Promise.reject(new Error("WebSocket connection not open"));
    }
    const id = new Date().getTime();
    const message = { id, type, data };
    const promise = new Promise<any>((resolve, reject) => {
      this.messages.set(id, [resolve, reject]);
    });
    this.socket.send(JSON.stringify(message));
    return promise;
  }

  public isConnected(): boolean {
    return !!this.socket && this.socket.readyState === WebSocket.OPEN;
  }

  public disconnect(): void {
    this.socket?.close();
  }

  public setInputHandler(func: (message: any) => void): void {
    this.InputHandler = func;
  }
}

class router {
  public static route(
    renderInstance: render,
    route: string,
    message: any
  ): void {
    console.log("route", route, message);

    switch (route) {
      case "message":
        console.log("Message", message);
        renderInstance.renderChatItem(message);

        break;
      default:
        console.log("Unknown route", route);
    }
  }
}
