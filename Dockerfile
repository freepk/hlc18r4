FROM ubuntu

ADD hlc18r4 /

EXPOSE 80

ENTRYPOINT ["/hlc18r4"]

