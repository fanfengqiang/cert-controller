# cert-controller

作为一个苦逼的屌丝，我们是没有money去买TLS证书的，怎么办呢？当然是去找免费的证书签发机构啦！
我们经常用的免费证书签发机构有Let's encrypt.但证书有效期只有短短的90天，你说天天让人换证书是不是挺麻烦的，而且生成证书再导入K8S还需要手动完成。
有没有一种方法让K8S自动的签署证书并保存为secret，并在证书快过期的时候自动更新？
在我的一番Google之下还真找到一个项目[cert-manager](https://github.com/jetstack/cert-manager),但是！！！
这个项目需要将要生成https证书的domain指向K8S的ingress，而且必须是一个公网地址（感兴趣的可以自行研究下）。而我的测试环境DNS是指向私网地址的!这不是蛋疼吗！
怎么办？自己动手丰衣足食。

## 功能介绍

- 安装cert-controller后，通过添加cert自定义资源，实现证书自动签发
- 将签发的证书自动保存到K8S的secret中
- 设置secret有效器，过期后重新签发证书，并自动更细secret

## 使用说明

1. 创建自定义资源CRD

```bash
kubectl apply -f https://raw.githubusercontent.com/fanfengqiang/cert-controller/master/deploy/cert-controller-rbac.yaml
kubectl apply -f https://raw.githubusercontent.com/fanfengqiang/cert-controller/master/deploy/cert-controller-deploy.yaml
```
2. 编写cert资源清单并应用

```bash
cat > cert.yaml << EOF
apiVersion: certcontroller.5ik8s.com/v1beta1
kind: Cert
metadata:
  name: example-cert
spec:
  secretName: cert-5ik8s.com
  domain: 5ik8s.com
  validityPeriod: 30
  type: dns_ali
  env:
    Ali_Key: "XXXXXXXXXXXXXXXXXXXXXXXX"
    Ali_Secret: "XXXXXXXXXXXXXXXXXXXXXXX"
EOF
kubectl apply -f cert.yaml
```
3. 参数定义

   | 参数                 | 含义                               |
   | :------------------- | :--------------------------------- |
   | .metadata.name       | cert资源的名字                     |
   | .spec.secretName     | 生成的secret的名字                 |
   | .spec.domain         | 生成证书的域名                     |
   | .spec.validityPeriod | secret的有效时长，单位天，范围1~89 |
   | .spec.type           | 域名托管商的标示                   |
   | .spec.env            | 域名托管商API的accesskey和secret   |

   完整域名托管商的格式，accesskey和secret格式[参见](https://github.com/fanfengqiang/cert-controller/blob/master/docs/dnsapi.md)

## 手动构建

```bash
go get github.com/fanfengqiang/cert-controller
cd $GOPATH/src/github.com/fanfengqiang/cert-controller
docker build -t cert-controller:latest .
```

## 致谢

本项目参考了如下两个项目

- [acme.sh](https://github.com/Neilpang/acme.sh)

- [sample-controller](<https://github.com/kubernetes/sample-controller>)