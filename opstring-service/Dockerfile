FROM alpine

# 容器主目录
WORKDIR /app

# 拷贝项目
COPY opstring .

# 创建配置文件目录和日志文件目录
RUN mkdir -p /app/bootstrap
RUN mkdir -p /app/logs

# 暴露端口
EXPOSE 5240 
EXPOSE 5250 

# 运行时命令
CMD ./opstring > /app/logs/opstring$RANDOM.$$ 2>&1
