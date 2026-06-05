# =========================================================
# AliStroy frontend (Nuxt) — multi-stage Dockerfile
# =========================================================

# ---- build stage ----
FROM node:22-alpine AS builder

WORKDIR /app

# Build-time public env vars (baked into the SSR build for hydration)
ARG NUXT_PUBLIC_API_BASE=/api/v1
ARG NUXT_PUBLIC_SITE_URL=http://localhost:3000
ARG NUXT_PUBLIC_UPLOADS_BASE=/uploads
ARG NUXT_PUBLIC_DEFAULT_LOCALE=tg

ENV NUXT_PUBLIC_API_BASE=${NUXT_PUBLIC_API_BASE} \
    NUXT_PUBLIC_SITE_URL=${NUXT_PUBLIC_SITE_URL} \
    NUXT_PUBLIC_UPLOADS_BASE=${NUXT_PUBLIC_UPLOADS_BASE} \
    NUXT_PUBLIC_DEFAULT_LOCALE=${NUXT_PUBLIC_DEFAULT_LOCALE}

COPY package.json package-lock.json* ./
RUN npm install --no-audit --no-fund

COPY . .
RUN npm run build

# ---- runtime stage ----
FROM node:22-alpine AS runtime

WORKDIR /app

ENV NODE_ENV=production \
    NUXT_HOST=0.0.0.0 \
    NUXT_PORT=3000 \
    PORT=3000

COPY --from=builder /app/.output ./.output
COPY --from=builder /app/package.json ./package.json

RUN addgroup -S app && adduser -S app -G app && chown -R app:app /app
USER app

EXPOSE 3000

CMD ["node", ".output/server/index.mjs"]
