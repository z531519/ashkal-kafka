FROM node:18-alpine

RUN apk --no-cache add \
      bash \
      g++ \
      ca-certificates \
      lz4-dev \
      musl-dev \
      cyrus-sasl-dev \
      openssl-dev \
      make \
      python3

RUN apk add --no-cache --virtual .build-deps gcc zlib-dev libc-dev bsd-compat-headers py-setuptools bash

USER node
WORKDIR /app
ENV NODE_ENV production

COPY --chown=node package*.json tsconfig.json ./
RUN npm ci

COPY --chown=node config ./config
COPY --chown=node src ./src

ENV HOME=/app


EXPOSE 8080

CMD ["npm", "start"]