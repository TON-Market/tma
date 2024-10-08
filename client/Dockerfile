ARG NODE_VERSION=21.6.1-alpine

FROM node:${NODE_VERSION} as builder

WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build


FROM alpine:latest

COPY --from=builder /app/dist /app/dist
WORKDIR /app
CMD ["/bin/sh"]