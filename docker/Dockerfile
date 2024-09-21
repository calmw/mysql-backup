# SPDX-License-Identifier: LGPL-3.0-only

FROM  golang AS builder
ADD . /src
WORKDIR /src
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn
RUN go mod tidy
RUN go build -o /backup .

# final stage
FROM ubuntu
RUN apt-get -y update && apt-get -y upgrade && apt-get install mysql-client zip ca-certificates wget -y

COPY --from=builder /backup ./
RUN chmod +x ./backup
RUN mkdir -p "/src/backup_data"

ENTRYPOINT ["./backup"]
