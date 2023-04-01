export class MessageWebSocket {
  private socket: WebSocket | null;
  private messages: Map<
    number,
    [(message: any) => void, (error: Error) => void]
  >;
  private InputHandler: (message: any) => void;

  constructor(private url: string) {
    this.socket = null;
    this.messages = new Map();
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
        if (this.messages.has(message.id)) {
          const [resolve, reject] = this.messages.get(message.id)!;
          this.messages.delete(message.id);
          resolve(message);
        } else {
          console.log("Message not found", message);
          if (message.type == "input") {
            this.InputHandler(message.data);
          }
        }
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
