opstring: 字符串服务，用来熟悉go kit框架

auth： 授权鉴权服务, 基于go kit框架开发, Atoken生成采用jwt方式

gateway: 利用httputil中反向代理包开发的网关服务，支持http服务。

im：基于websockt的IM服务

study：日常学习总结

注意：
1）创建容器开发环境时，redis，rabbitmq, swagger需要在宿主机的/data目录下创建配置文件和数据文件映射目录.
也可以解压tools目录下的zip包到宿主机的/data目录，然后执行如下命令：
docker-compose -f docker-compose_xx.yml up -d
docker-compose -f docker-compose_xx.yml down

2）在k8s上运行时，需要在集群中安装consul。tools/kubernetes/consul.yaml文件是开发版consul.

3）auth通过truss生成，帮助参见 https://github.com/metaverse/truss
