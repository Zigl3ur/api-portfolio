import { Hono } from "hono";
import { LastFMHandler } from "./handlers/music";
import { cors } from "hono/cors";

type Bindings = {
  LASTFM_API_KEY: string;
};

const app = new Hono<{ Bindings: Bindings }>();

app.use(
  "/api/*",
  cors({
    origin: "https://eden.douru.fr",
    allowMethods: ["GET"],
    allowHeaders: ["Content-Type"],
  })
);

app.get("/api/music", async (c) => {
  const LASTFM_API_KEY = c.env.LASTFM_API_KEY;

  const json = await LastFMHandler(LASTFM_API_KEY);

  return c.json(json, 200, {
    "Content-Type": "application/json",
  });
});

export default app;
