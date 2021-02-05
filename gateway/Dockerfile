FROM alpine

# 容器主目录
WORKDIR /app

# 拷贝项目
COPY gateway .

# 创建配置文件目录和日志文件目录
RUN mkdir -p /app/bootstrap
RUN mkdir -p /app/logs

# 暴露端口
EXPOSE 5040 
EXPOSE 5050 
EXPOSE 6000 

# 运行时命令
CMD ./gateway > /app/logs/gateway$RANDOM.$$ 2>&1
