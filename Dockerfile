FROM golang:latest

LABEL maintainer="noura.r.saad@gmail.com"
WORKDIR /data
# Copy only the necessary files for building
COPY . .
RUN go mod download
#
# Install make
USER root
#RUN apk --no-cache add make
RUN make


USER nobody
CMD  ["./main"]