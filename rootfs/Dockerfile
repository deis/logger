FROM quay.io/deis/base:v0.3.6

# Add logger user and group
RUN adduser --system \
	--shell /bin/bash \
	--disabled-password \
	--home /opt/logger \
	--group \
	logger

COPY . /

# Fix some permission since we'll be running as a non-root user
RUN chown -R logger:logger /opt/logger

USER logger

CMD ["/opt/logger/sbin/logger"]
EXPOSE 1514 8088
