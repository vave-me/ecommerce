FROM node:20-alpine AS builder
LABEL stage=builder

WORKDIR /app

ARG CACHE_BUST=1

RUN apk add --no-cache \
    autoconf \
    automake \
    libtool \
    nasm \
    make \
    pkgconfig \
    g++ \
    python3 \
    build-base

RUN mkdir -p /app

# Copy package.json and package-lock.json before installing dependencies
COPY frontend/package*.json ./

RUN if [ -f package-lock.json ]; then \
      npm ci --legacy-peer-deps; \
    else \
      npm install --legacy-peer-deps; \
    fi

COPY frontend/. ./

RUN rm -rf .next/cache || true

RUN npm run build

FROM node:20-alpine
LABEL stage=runtime

WORKDIR /app

ENV NODE_ENV=production
ENV NEXT_PUBLIC_API_BASE_URL=http://192.168.178.84:8080/api
ENV API_URL=http://192.168.178.84:8080

RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

COPY --from=builder --chown=nextjs:nodejs /app/public       ./public
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone  ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static  ./.next/static
COPY --from=builder --chown=nextjs:nodejs /app/next.config.mjs ./next.config.mjs

USER nextjs

EXPOSE 3000
CMD ["node", "server.js"]