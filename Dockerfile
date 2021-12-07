FROM debian

ENV NEW_RELIC_LICENSE_KEY=""
ENV NEW_RELIC_APP_NAME=""
ENV WORKLOAD_NAME=""

ADD start.sh /
ADD infra-lite /
RUN apt update && apt install -y ca-certificates

CMD ["/start.sh"]
