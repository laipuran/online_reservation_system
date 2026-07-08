import { setupWorker } from "msw/browser";
import { handlers } from "./handlers";
import { seed } from "./db";

seed();

export const worker = setupWorker(...handlers);
