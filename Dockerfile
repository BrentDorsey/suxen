FROM alpine:latest

COPY ./bin /usr/local/bin/
COPY ./ui/dist /dist
COPY ./scripts/docker-entrypoint.sh /

RUN apk --no-cache add curl

EXPOSE 8080

ENTRYPOINT ["/docker-entrypoint.sh"]
CMD ["suxend"]
