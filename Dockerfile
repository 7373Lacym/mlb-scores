# create image from the official Go image
FROM golang:alpine
 
RUN mkdir /app
ADD . /app/
WORKDIR /app
COPY cron /var/spool/cron/crontabs/root


RUN go build -o main .
RUN chmod +x main
RUN chmod +r config.json

CMD /usr/sbin/crond -l 2 -f
