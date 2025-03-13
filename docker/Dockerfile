# Stage 1: Build the Go application
FROM golang:1.21-bullseye AS go-builder
WORKDIR /app
COPY .. .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o blogflex .

# Stage 2: Build the frontend assets
FROM node:18-bullseye AS node-builder
WORKDIR /app
COPY ../package.json ../package-lock.json ./ 
RUN npm install
COPY ../tailwind.config.js ../postcss.config.js ./ 
COPY ../static ./static
RUN npm run build

# Stage 3: Final image
FROM debian:bullseye-slim
RUN apt-get update && apt-get install -y curl
WORKDIR /app
COPY --from=go-builder /app/blogflex .
COPY --from=node-builder /app/static ./static

COPY public ./public
COPY ../views ./views
# COPY ../.env ./.env  
EXPOSE 8080
CMD ["./blogflex"]
