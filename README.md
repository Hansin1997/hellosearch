# Hello Search 🔍
白度搜索，心应搜索，G00gle？🤣

## 简介
> 此项目并不涉及搜索引擎的任何技术，纯粹是 ```Elasticsearch``` 的应用展示。

基于 [Hello Spider](https://github.com/Hansin1997/hellospider) 爬虫抓取的数据，
结合 ```Elasticsearch``` 实现的简单搜索应用。

### 后端
基于 ```Go``` & ```Gin``` 提供 RESTful 接口。

运行前修改 ```config.json``` 配置文件中的 ```Elasticsearch``` 地址、用户名、密码。

### 前端
基于 ```Vue``` & ```Quasar``` 实现单页面应用。

## 展示
### 主页
![主页](docs/img/home.jpeg?raw=true)

### 搜索
![搜索](docs/img/search.jpeg?raw=true)