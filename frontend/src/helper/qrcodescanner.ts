import { Html5Qrcode } from "html5-qrcode";

export class QrScanner {
  private scanner: Html5Qrcode;

  constructor() {
    this.scanner = new Html5Qrcode("qrscanner");
  }

  public startScan(): Promise<string> {
    return new Promise((resolve, reject) => {
      this.scanner
        .start(
          { facingMode: "environment" }, // Prefer rear camera
          {
            fps: 10, // Optional frames per second
            qrbox: 250, // Optional, if you want a QR box
          },
          (qrCodeMessage: string) => {
            this.scanner
              .stop()
              .then((_) => resolve(qrCodeMessage))
              .catch((error: any) => reject(error));
          },
          (errorMessage: any) => {
            reject(errorMessage);
          }
        )
        .catch((error: any) => reject(error));
    });
  }
}
