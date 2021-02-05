#编译
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gateway cmd/gateway/main.go

#构建镜像
docker build -t gateway .

#运行
docker run --name=gateway -v /data/ziyun/gateway/bootstrap/:/app/bootstrap -v /data/ziyun/gateway/logs:/app/logs -p 5040:5040 -p 5050:5050 -p 6000:6000 -d gateway

#测试
curl 192.168.73.3:5050/health

#推送到阿里云
docker tag gateway registry.cn-beijing.aliyuncs.com/492258027/ziyun/gateway:v5.0
docker push registry.cn-beijing.aliyuncs.com/492258027/ziyun/gateway:v5.0

#创建configmap
kubectl create configmap gateway --from-file=./bootstrap/bootstrap.yaml

#k8s中运行
kubectl apply -f gateway.yaml

#配置文件注意项
当docker环境运行时，bootstrap.yml中http和rpc的host置空，会自动找到docker容器的地址
