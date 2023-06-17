import { Html5QrcodeScanType, Html5QrcodeScanner } from "html5-qrcode";
import { Communications } from "../communications";
import { QrScanner } from "../helper/qrcodescanner";

export class QrLoginScanner extends HTMLElement {
  message: any;
  coms: Communications;
  id: string;

  constructor() {
    super();
    this.attachShadow({ mode: "open" });
  }

  async init(coms: Communications) {
    await this.loadTemplate();
    this.attachScanButton();
    this.coms = coms;
  }

  private async loadTemplate() {
    const response = await fetch("./qrLoginScanner.html");
    const html = await response.text();
    const fragment = document.createRange().createContextualFragment(html);
    (this as any).shadowRoot.appendChild(fragment);
  }

  attachScanButton() {
    console.log("attachDeleteBtnEvent");
    this.shadowRoot
      .querySelector(".qr-login-scan")
      .addEventListener("click", () => {
        function onScanSuccess(decodedText: any, decodedResult: any) {
          // handle the scanned code as you like, for example:
          console.log(`Code matched = ${decodedText}`, decodedResult);
        }

        let config = {
          fps: 10,
          qrbox: { width: 500, height: 500 },
          rememberLastUsedCamera: true,
          // Only support camera scan type.
          supportedScanTypes: [Html5QrcodeScanType.SCAN_TYPE_CAMERA],
        };

        let html5QrcodeScanner = new Html5QrcodeScanner(
          "qrscanner",
          config,
          /* verbose= */ false
        );
        html5QrcodeScanner.render(onScanSuccess, (errorMessage) =>
          console.log(errorMessage)
        );
      });
  }
}
