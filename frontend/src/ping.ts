import { Communications } from "./communications";

export class Ping {
  private comms: Communications;
  private intervalId: number | null;
  private latencyDiv: HTMLDivElement;

  constructor(comms: Communications) {
    this.comms = comms;
    this.intervalId = null;
    this.latencyDiv = document.getElementById("ping") as HTMLDivElement;

    // Registering ping response in router
    this.comms.Router.register("pong", this.calculateLatency.bind(this));
  }

  public startPing(interval: number): void {
    this.intervalId = window.setInterval(() => {
      const pingTime = Date.now();
      this.comms.send("ping", "ping", JSON.stringify({ time: pingTime }));
      console.log("ping");
    }, interval);
  }

  public stopPing(): void {
    if (this.intervalId !== null) {
      window.clearInterval(this.intervalId);
      this.intervalId = null;
    }
  }

  private calculateLatency(message: any): void {
    const pongTime = Date.now();
    const pingTime = JSON.parse(message.data).time;
    const latency = pongTime - pingTime;

    // Writing latency to div
    this.latencyDiv.innerText = `Latency: ${latency}ms`;
  }
}
