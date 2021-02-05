#编译
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o opstring cmd/opstring/main.go

#构建镜像
docker build -t opstring .

#运行
docker run --name=opstring  -v /data/ziyun/opstring-service/bootstrap/:/app/bootstrap -v /data/ziyun/opstring-service/logs:/app/logs  -p 5240:5240 -p 5250:5250 -d opstring

#测试
curl 192.168.73.3:5250/health

#推送到阿里云
docker tag opstring registry.cn-beijing.aliyuncs.com/492258027/ziyun/opstring:v8.0
docker push registry.cn-beijing.aliyuncs.com/492258027/ziyun/opstring:v8.0

#创建configmap
kubectl create configmap opstring --from-file=./bootstrap/bootstrap.yaml

#k8s中运行
kubectl apply -f opstring.yaml

#配置文件注意项
当docker环境运行时，bootstrap.yml中http和rpc的host置空，会自动找到docker容器的地址
