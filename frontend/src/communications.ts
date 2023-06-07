export class Communications {
  private socket: WebSocket | null;
  public Router: router;
  private messages: Map<
    string,
    [(message: any) => void, (error: Error) => void]
  >;

  constructor(private url: string) {
    this.socket = null;
    this.messages = new Map();
    this.Router = new router();
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
        console.log("message", message);
        this.Router.route(message.type, message);
      });
    });
  }

  public send(
    id: string,
    type: string,
    data: string,
    path: string = null
  ): Promise<any> {
    if (!this.isConnected()) {
      return Promise.reject(new Error("WebSocket connection not open"));
    }
    const message = { id, type, data, path };
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
}

type RouteFunction = (message: any) => void;

class router {
  private routes: Record<string, RouteFunction> = {};

  public route(route: string, message: any): void {
    console.log("route", route, message);

    if (this.routes[route]) {
      console.log("route", route, message);
      this.routes[route](message);
    }
  }

  public register(route: string, callback: RouteFunction): void {
    this.routes[route] = callback;
  }
}
