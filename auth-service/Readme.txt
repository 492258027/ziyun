#编译
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o auth cmd/auth/main.go

#构建镜像
docker build -t auth .

#运行
docker run --name=auth -v /data/ziyun/auth-service/bootstrap/:/app/bootstrap -v /data/ziyun/auth-service/logs:/app/logs  -p 5140:5140 -p 5150:5150 -d auth

#测试
curl 192.168.73.3:5150/health

#推送到阿里云
docker tag auth registry.cn-beijing.aliyuncs.com/492258027/ziyun/auth:v8.0
docker push registry.cn-beijing.aliyuncs.com/492258027/ziyun/auth:v8.0

#创建configmap
kubectl create configmap auth --from-file=./bootstrap/bootstrap.yaml

#k8s中运行
kubectl apply -f auth.yaml

#配置文件注意项
当docker环境运行时，bootstrap.yml中http和rpc的host置空，会自动找到docker容器的地址