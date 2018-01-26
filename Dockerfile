FROM alpine:3.5

RUN apk add --update ca-certificates && rm -Rf /tmp/* /var/lib/cache/apk/*

ADD bin/azip /usr/bin/
RUN chmod +x /usr/bin/azip

CMD [ "azip" ]
