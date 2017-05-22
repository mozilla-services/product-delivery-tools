FROM golang:1.8

ENV PROJECT=github.com/mozilla-services/product-delivery-tools

COPY version.json /app/version.json
COPY . /go/src/$PROJECT

RUN go install $PROJECT/post_upload && \
    go install $PROJECT/bucketlister
