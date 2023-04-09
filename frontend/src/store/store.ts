type Storage = {
  name: string;
  items: StorageItem[];
  callbacks: ((value: Storage) => void)[];
};

export class Store {
  private data: { [key: string]: Storage } = {};
  constructor() {}

  get(key: string): Storage {
    return this.data[key];
  }

  set(key: string, value: Storage) {
    this.data[key] = value;
    this.data[key].callbacks.forEach((callback) => callback(value));
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
