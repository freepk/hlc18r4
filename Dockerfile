FROM ubuntu

ADD hlc18r4 /

EXPOSE 80

ENTRYPOINT ["GOGC=off /hlc18r4"]

