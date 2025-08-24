FROM oven/bun:alpine AS builder

WORKDIR /app
COPY src ./
COPY package.json bun.lock ./

RUN bun install --production
RUN bun build index.ts --outdir=dist --target=bun

FROM oven/bun:alpine AS runner

WORKDIR /app
COPY --from=builder /app/dist .

EXPOSE 3000
ENV NODE_ENV=production
CMD ["bun", "index.js"]