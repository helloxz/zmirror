FROM alpine:3.21
RUN mkdir -p /opt/zmirror/data
WORKDIR /opt/zmirror
COPY bin/zmirror /opt/zmirror/
COPY static /opt/zmirror/static
# 暴露文件夹和端口
EXPOSE 5080
# 启动程序
CMD ["./zmirror"]