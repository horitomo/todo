# ベースとなるDockerイメージ指定
FROM golang:latest
# コンテナ内に作業ディレクトリを作成
RUN mkdir /go/src/todo
# コンテナログイン時のディレクトリ指定
WORKDIR /go/src/todo
RUN apt-get update
RUN cd /go/src/todo
RUN go get github.com/gin-gonic/gin
RUN go get github.com/jinzhu/gorm
RUN go get github.com/mattn/go-sqlite3
# SQLite
RUN apt-get install sqlite3 libsqlite3-dev -y
# ホストのファイルをコンテナの作業ディレクトリに移行
ADD . /go/src/todo
