import mitt from "mitt"

export type AppEvents = {
  triggerAction: {
    action: string;
    arguments: any;
  }
}

export const emitter = mitt<AppEvents>();
