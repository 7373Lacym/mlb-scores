# create image from the official Go image
FROM golang:alpine
 
RUN mkdir /app
ADD . /app/
WORKDIR /app
COPY cron /var/spool/cron/crontabs/root


RUN go build -o main .
RUN chmod +x main
run chmod +r config

CMD /usr/sbin/crond -l 2 -f
