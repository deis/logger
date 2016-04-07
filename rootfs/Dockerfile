FROM alpine:3.3

# add logger user
RUN mkdir -p /opt/logger && \
		addgroup -S logger && \
		adduser -S -G logger -H -h /opt/logger -D logger

ADD . /

# Fix some permission since we'll be running as a non-root user.
RUN chown -R logger:logger /opt/logger

USER logger

CMD ["/opt/logger/sbin/logger"]
EXPOSE 1514 8088

ENV DEIS_RELEASE v2-beta
