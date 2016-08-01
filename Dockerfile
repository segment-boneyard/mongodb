FROM golang:1.6

RUN go get -u github.com/segment-sources/mongodb

# Include a small utility to allow users to use json instead of yaml
# docker run -rm --entrypoint bash "yaml2json schema.yml > schema.json && mongodb ..."
RUN go get -u github.com/bronze1man/yaml2json

# Additionally add a cron-like runner to run on an
# interval.
RUN go get -u github.com/segmentio/go-every

ENTRYPOINT ["mongodb"]