export class Communications {
  private socket: WebSocket | null;
  public Router: router;
  private messages: Map<
    string,
    [(message: any) => void, (error: Error) => void]
  >;
  private reconnectAttempts: number;
  private maxReconnectAttempts: number;
  private reconnectInterval: number;

  constructor(
    private url: string,
    maxReconnectAttempts = 300,
    reconnectInterval = 2000
  ) {
    this.socket = null;
    this.messages = new Map();
    this.Router = new router();
    this.reconnectAttempts = 0;
    this.maxReconnectAttempts = maxReconnectAttempts;
    this.reconnectInterval = reconnectInterval;
  }

  public connect(): Promise<void> {
    return new Promise((resolve, reject) => {
      this.socket = new WebSocket(this.url);

      this.socket.addEventListener("open", () => {
        console.log("WebSocket connection established");
        this.reconnectAttempts = 0; // Reset reconnection attempts
        resolve();
      });

      this.socket.addEventListener("error", (event) => {
        reject(event);
      });

      this.socket.addEventListener("close", () => {
        console.log("WebSocket connection closed");
        if (this.reconnectAttempts < this.maxReconnectAttempts) {
          setTimeout(() => {
            console.log(
              `WebSocket attempting to reconnect. Attempt number ${
                this.reconnectAttempts + 1
              }`
            );
            this.reconnectAttempts++;
            this.connect(); // Try to reconnect
          }, this.reconnectInterval);
        } else {
          console.log(
            "WebSocket connection could not be re-established after maximum attempts"
          );
        }
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
    if (this.socket) {
      this.socket.close();
      this.socket = null;
    }
  }
}

type RouteFunction = (message: any) => void;

class router {
  private routes: Record<string, RouteFunction> = {};
  private threads: Record<string, RouteFunction> = {};

  public route(route: string, message: any): void {
    console.log("route", route, message);

    if (route.startsWith("thread")) {
      console.log("thread", message);
      if (this.threads[message.id]) {
        this.threads[message.id](message);
      } else {
        console.log("no thread", message);
      }
      return;
    }

    if (this.routes[route]) {
      console.log("route", route, message);
      this.routes[route](message);
    }
  }

  public register(route: string, callback: RouteFunction): void {
    this.routes[route] = callback;
  }

  public registerThread(id: string, callback: RouteFunction): void {
    this.threads[id] = callback;
  }
}
