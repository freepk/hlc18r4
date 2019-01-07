FROM ubuntu

ADD hlc18r4 /

EXPOSE 80

ENV GOGC off

ENTRYPOINT ["/hlc18r4"]

