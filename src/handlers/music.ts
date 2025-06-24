import type { LastFMData } from "../types/types";

export async function LastFMHandler(LASTFM_API_KEY: string) {
  const response = await fetch(
    `https://ws.audioscrobbler.com/2.0/?method=user.getrecenttracks&user=zigl3ur&api_key=${LASTFM_API_KEY}&format=json`
  );

  if (!response.ok) return { isListening: false };

  const data: LastFMData = (await response.json()) as LastFMData;

  const tracks = data.recenttracks.track;
  const isListening =
    tracks.length > 0 && tracks[0]["@attr"]?.nowplaying === "true";

  if (isListening)
    return {
      isListening: isListening,
      track: {
        artist: tracks[0].artist["#text"],
        album: tracks[0].album["#text"],
        name: tracks[0].name,
        image:
          tracks[0].image.find((img) => img.size === "large")?.["#text"] || "",
        url: tracks[0].url,
      },
    };

  return { isListening: isListening };
}
