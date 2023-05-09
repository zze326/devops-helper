FROM registry-azj-registry.cn-shanghai.cr.aliyuncs.com/ops/alpine-with-certs:v1.0
WORKDIR /opt/app
ADD devops-helper.tar.gz /opt/app

EXPOSE 8888
CMD ["/opt/app/run.sh"]
