# 短链接跳转系统

## 构建 & 运行

运行时需要讲`config.toml`与二进制包放在同一目录下并按照需要进行配置

```bash
# 运行
$ ./dev/run.sh
# 发布
$ ./dev/build.sh
$ ./release/linux/short-url
```

## 设计要点

- 通过 [Memo](https://github.com/yourtion/go-utils/blob/master/memo) 模块防止并发时对数据库并发请求
- 通过 [Cache](https://github.com/yourtion/go-utils/blob/master/cache) 对跳转信息进行缓存
- 通过 [statics](api/redirect/statistics.go) 先将统计信息存入内存变量再定时同步到 Redis 与 MySQL 降低统计压力
- 每日持久化Redis数据到数据库

性能：在阿里云 4C8G 机器上，可以达到 5w QPS

## 功能

- [x] 短链接跳转
- [x] 可选PV/UV记录
- [x] 可选AccessLog记录
- [ ] 跳转带上时间戳
- [ ] 指定短链接域名
- [ ] 修改重定向地址
- [x] 通过API生成短链接
- [ ] API生成短链接域名白名单
- [ ] 活动管理与密钥生成
- [ ] 聚合PV/UV到对应的活动（按活动ID算UV）
- [ ] 活动有效期配置与活动关闭
- [ ] 根据活动生成短链接
- [ ] 实时访问记录
- [ ] 访问信息图表功能 
