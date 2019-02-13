FROM centos

ADD hlc18r4 /

EXPOSE 80

ENV GOGC off

CMD ["/hlc18r4"]
