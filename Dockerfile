FROM scratch

WORKDIR /

COPY --from=golang /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY ./assets/carriers.json /carriers.json
COPY ./bin/main /app

CMD [ "/app" ]
