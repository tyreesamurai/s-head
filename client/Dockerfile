# --- build ---
FROM node:20-bookworm AS build
WORKDIR /app
ENV NEXT_TELEMETRY_DISABLED=1
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

# --- run (standalone Node) ---
FROM node:20-bookworm
ENV NODE_ENV=production
WORKDIR /app
COPY --from=build /app/.next/standalone ./
COPY --from=build /app/public ./public
COPY --from=build /app/.next/static ./.next/static

USER node
ENV PORT=3000
EXPOSE 3000
CMD ["node", "server.js"]
