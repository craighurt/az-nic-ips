FROM alpine:3.5

RUN apk add --update bash ca-certificates jq curl && rm -Rf /tmp/* /var/lib/cache/apk/*

ADD bin/azip /usr/bin/
ADD init.sh /usr/bin/
RUN chmod +x /usr/bin/azip
RUN chmod +x /usr/bin/init.sh

CMD [ "/usr/bin/init.sh" ]
