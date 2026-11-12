# ---- Stage 1: build Angular frontend (project root) ----
FROM node:22-alpine AS frontend-build
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build -- --configuration production

# ---- Stage 2: build Go backend ----
FROM golang:1.27-alpine AS backend-build
WORKDIR /app
COPY backend/go.mod backend/go.sum* ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /qualiopi .

# ---- Stage 3: runtime ----
FROM alpine:3.20
WORKDIR /app
RUN adduser -D app
COPY --from=backend-build /qualiopi ./qualiopi
COPY --from=frontend-build /app/dist/Qualiopi/browser ./public
USER app

ENV PORT=4200
EXPOSE 4200
CMD ["./qualiopi"]