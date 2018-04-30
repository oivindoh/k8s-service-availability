FROM golang:1.10.1-alpine3.7 as buildstage
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH
WORKDIR $GOPATH/src/app
RUN apk add --no-cache curl git \
  && mkdir -p /go/src \
  && mkdir -p /go/bin \
  && mkdir -p /go/pkg \
  && mkdir -p /go/src/app \
  && curl https://glide.sh/get | sh

ADD ./glide.yaml $GOPATH/src/app
ADD ./main.go $GOPATH/src/app

RUN glide up -v

ENV CGO_ENABLED=0
ENV GOOS=linux
RUN go build -ldflags '-w -s' -a -installsuffix cgo -o ./k8stest .

FROM scratch
WORKDIR /app
ENV PATH=/app;$PATH
COPY --from=buildstage /go/src/app/k8stest /app
CMD ["./k8stest"]