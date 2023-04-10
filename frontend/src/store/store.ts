type Storage = {
  name: string;
  items: StorageItem[];
  callbacks: ((value: Storage) => void)[];
};

export class Store {
  private data: { [key: string]: Storage } = {};
  constructor() {}

  get(key: string): Storage {
    if (!this.data[key]) {
      this.data[key] = {
        name: key,
        items: [],
        callbacks: [],
      };
    }
    return this.data[key];
  }

  set(key: string, value: Storage) {
    this.data[key] = value;
    this.data[key].callbacks.forEach((callback) => callback(value));
  }

  addStorageItem(key: string, value: string) {
    console.log("addStorageItem", key, value);
    let storage = this.data[key];
    if (!storage) {
      this.data[key] = {
        name: key,
        items: [],
        callbacks: [],
      };
      storage = this.data[key];
    }
    const item = new StorageItem(value);
    storage.items.push(item);
    storage.callbacks.forEach((callback) => callback(storage));
  }

  subscribe(key: string, callback: (value: Storage) => void) {
    this.data[key].callbacks.push(callback);
  }
}

type ItemContent = {
  created: Date;
  content: string;
};

export class StorageItem {
  created: Date;
  content: ItemContent[];
  callbacks: ((value: StorageItem) => void)[];

  constructor(content: string) {
    this.created = new Date();
    this.content = [{ created: this.created, content }];
    this.callbacks = [];
  }
}
