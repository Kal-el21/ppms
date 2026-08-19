FROM node:24-alpine AS builder

WORKDIR /app

COPY package*.json ./
# Install ALL dependencies including devDependencies for building
RUN npm ci

COPY . .
RUN npm run build

# Production stage with Nginx
FROM nginx:alpine

# Copy built assets
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy nginx config
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]