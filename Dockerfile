FROM ubuntu:14.04

ADD bin/megaboom /bin/megaboom

EXPOSE 8080

CMD /bin/megaboom
