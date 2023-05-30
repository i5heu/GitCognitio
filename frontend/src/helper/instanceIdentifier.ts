export class InstanceIdentifier {
  private static readonly cookieName: string = "instanceIdentifier";

  static getInstanceIdentifier(): string {
    let instanceIdentifier = this.getCookie();

    if (!instanceIdentifier) {
      instanceIdentifier = this.generateInstanceIdentifier();
      this.setCookie(instanceIdentifier);
    }

    return instanceIdentifier;
  }

  private static getCookie(): string | null {
    const cookies = document.cookie.split(";");

    for (let i = 0; i < cookies.length; i++) {
      const cookie = cookies[i].trim();

      if (cookie.startsWith(this.cookieName + "=")) {
        return cookie.substring(this.cookieName.length + 1);
      }
    }

    return null;
  }

  private static setCookie(value: string, days = 365): void {
    const expirationDate = new Date();
    expirationDate.setDate(expirationDate.getDate() + days);

    const cookie = `${
      this.cookieName
    }=${value}; expires=${expirationDate.toUTCString()}; path=/`;
    document.cookie = cookie;
  }

  private static generateInstanceIdentifier(): string {
    // Generate a unique identifier using any desired method
    // For simplicity, let's use a random string with 10 characters
    const characters =
      "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
    let instanceIdentifier = "";

    for (let i = 0; i < 10; i++) {
      instanceIdentifier += characters.charAt(
        Math.floor(Math.random() * characters.length)
      );
    }

    return instanceIdentifier;
  }
}
