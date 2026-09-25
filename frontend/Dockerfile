# Stage 1: Сборка статических файлов
FROM node:22-alpine AS builder
WORKDIR /src
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

# Stage 2: Runtime (Nginx)
FROM nginx:1.27-alpine AS runtime
COPY --from=builder /src/dist /usr/share/nginx/html
EXPOSE 80
