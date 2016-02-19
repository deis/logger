FROM alpine:3.3

# install common packages
RUN apk add --update-cache \
	libcap \
	&& rm -rf /var/cache/apk/*

# add logger user
RUN addgroup -S logger && \
	adduser -S -G logger -H -h /opt/logger -D logger

ADD . /

# Fix some permission since we'll be running as a non-root user.  This includes enabling
# logger to always bind to low/privileged ports, even when running as a non-root user.
RUN chown -R logger:logger /opt/logger && \
	setcap 'cap_net_bind_service=+ep' /opt/logger/sbin/logger

USER logger

CMD ["/opt/logger/sbin/logger"]
EXPOSE 514 8088

ENV DEIS_RELEASE v2-beta
