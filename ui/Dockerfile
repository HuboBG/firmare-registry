# ============================
# 1. Build Stage
# ============================
FROM node:20-alpine AS builder

WORKDIR /app

# Copy package files
COPY package.json yarn.lock ./

# Install deps (frozen lockfile for deterministic builds)
RUN yarn install --frozen-lockfile

# Copy full project
COPY . .

# Build production bundle
RUN yarn build


# ============================
# 2. Production Stage (Nginx)
# ============================
FROM nginx:1.27-alpine AS production

# Remove default nginx website
RUN rm -rf /usr/share/nginx/html/*

# Copy built UI from builder stage
COPY --from=builder /app/dist /usr/share/nginx/html

# Optional: custom nginx config (good for SPA routing)
# COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]