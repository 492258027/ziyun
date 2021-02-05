FROM alpine

# 容器主目录
WORKDIR /app

# 拷贝项目
COPY auth .

# 创建配置文件目录和日志文件目录
RUN mkdir -p /app/bootstrap
RUN mkdir -p /app/logs

# 暴露端口
EXPOSE 5140 
EXPOSE 5150 

# 运行时命令
CMD ./auth > /app/logs/auth$RANDOM.$$ 2>&1
