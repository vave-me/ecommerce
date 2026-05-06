# ==================================
# Stage 1: Builder Environment
# ==================================
FROM node:20-alpine AS builder
LABEL stage=builder

# Set working directory
WORKDIR /app

# Optional: Argument to invalidate cache if needed (e.g., ARG CACHE_BUST=$(date +%s))
ARG CACHE_BUST=1

# Install build-time dependencies (Python, C++, etc. for native modules)
# These will NOT be included in the final image.
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

# Create app directory (redundant if WORKDIR creates it, but harmless)
RUN mkdir -p /app

# Copy package files first to leverage Docker cache for dependencies
COPY frontend/package*.json ./

# Install ALL dependencies (including devDependencies needed for build)
# Use npm ci if package-lock.json exists for deterministic installs
RUN if [ -f package-lock.json ]; then \
      npm ci; \
    else \
      echo "Warning: package-lock.json not found, using npm install" && \
      npm install; \
    fi

# Copy the rest of the frontend source code
COPY frontend/. ./

# Optional: Clear Next.js cache before building (might be useful sometimes)
RUN rm -rf .next/cache || true

# Build the Next.js application
RUN npm run build

# ==================================
# Stage 2: Runtime Environment
# ==================================
FROM node:20-alpine AS runtime # You can optionally name this stage
LABEL stage=runtime

WORKDIR /app

# Set production environment
ENV NODE_ENV=production
# Set runtime environment variables (can also be set via K8s)
ENV API_URL=https://api.sfx-markt.com

# Create non-root user and group
RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

# Copy only necessary package files for installing production dependencies
COPY --from=builder --chown=nextjs:nodejs /app/package*.json ./

# Install ONLY production dependencies
# Use --omit=dev with npm ci for cleaner production installs
RUN if [ -f package-lock.json ]; then \
      npm ci --omit=dev; \
    else \
      echo "Warning: package-lock.json not found, using npm install --production" && \
      npm install --production; \
    fi

# Copy built application artifacts from the builder stage
# Ensure correct ownership for the non-root user
COPY --from=builder --chown=nextjs:nodejs /app/public       ./public
COPY --from=builder --chown=nextjs:nodejs /app/.next        ./.next
COPY --from=builder --chown=nextjs:nodejs /app/next.config.mjs ./next.config.mjs
# Note: node_modules are installed fresh above, DON'T copy from builder

# Switch to the non-root user
USER nextjs

# Expose the port the app runs on
EXPOSE 3000

# Command to run the application
# Uses the start script defined in package.json
CMD ["npm", "start"]