# AOspace-Upgrade

[English](./README.md) | 简体中文

## 简介

按需启动，主要负责aospace-all-in-one的容器升级

## 构建

### 环境准备

- docker (>=18.09)
- git
- golang 1.18 +

### 源码下载

```shell
git clone git@github.com:ao-space/space-upgrade.git
```

### 容器镜像构建

进入模块根目录，执行命令

```shell
docker build -t local/space-upgrade:{tag} . 
````

其中 tag 参数可以根据实际情况修改，和服务器整体运行的 docker-compose.yml 保持一致即可。

## 运行

一般不会常驻运行，只有在进行升级时启动，升级完成后退出

如果要手动运行，请执行

```shell
docker exec -it aospace-all-in-one docker-compose -f /etc/ao-space/aospace-upgrade.yml up -d
```

## 贡献指南

我们非常欢迎对本项目进行贡献。以下是一些指导原则和建议，希望能够帮助您参与到项目中来。

[贡献指南](https://github.com/ao-space/ao.space/blob/dev/docs/cn/contribution-guidelines.md)

## 联系我们

- 邮箱：<developer@ao.space>
- [官方网站](https://ao.space)
- [讨论组](https://slack.ao.space)

## 感谢您的贡献

最后，感谢您对本项目的贡献。我们欢迎各种形式的贡献，包括但不限于代码贡献、问题报告、功能请求、文档编写等。我们相信在您的帮助下，本项目会变得更加完善和强大。
