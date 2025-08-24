import { Hono } from "hono";
import { LastFMHandler } from "./handlers/music";
import { cors } from "hono/cors";
import { CacheType } from "./types/types";
import { env } from "hono/adapter";

const app = new Hono();

app.use(
  "/api/*",
  cors({
    origin: "https://eden.douru.fr",
    allowMethods: ["GET"],
    allowHeaders: ["Content-Type"],
  })
);

let cache: CacheType = { data: { isListening: false }, timestamp: 0 };

app.get("/api/music", async (c) => {
  const { LASTFM_API_KEY } = env<{ LASTFM_API_KEY: string }>(c);
  const timeStart = Date.now();

  if (cache && timeStart - cache.timestamp < 30000) {
    await new Promise((resolve) => setTimeout(resolve, 1000));

    return c.json(cache.data, 200, {
      "Content-Type": "application/json",
    });
  }

  const json = await LastFMHandler(LASTFM_API_KEY);

  cache = {
    data: json,
    timestamp: timeStart,
  };

  return c.json(json, 200, {
    "Content-Type": "application/json",
  });
});

export default app;
