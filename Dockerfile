FROM ubuntu

ADD hlc18r4 .

EXPOSE 80

# ENV GOGC off

CMD ulimit -n 8192 && ulimit -n && ./hlc18r4

