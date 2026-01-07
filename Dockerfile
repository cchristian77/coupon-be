FROM alpine:latest

RUN mkdir /app

WORKDIR /app

COPY /env.json /app/env.json

COPY base-api-build /app

CMD [ "/app/base" ]