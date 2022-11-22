---
title: 如何贡献承影
sidebar_position: 3
---

## 如何提问题和需求？
当前提供了4个类型的问题分类：
Feature request
Ask a Question
Bug report
Development Task

## 如何提交pr?
1）新建issuse .描述问题信息。
2）基于对应的release分支拉取开发分支,分支命名参考 【分支命名规范】
3) commit 信息：[type-redmine] [module] msg，type 定义参考：【Commit type 类别】
eg:
[hotfix-31280][core] 修复bigdecimal转decimal运行失败问题
[feat-31372][rdb] RDB结果表Upsert模式支持选择更新策略
4）多次提交使用rebase 合并成一个。
5）pr 名称：[chengying-issuseid][module名称] 标题
6）对应模块的test 测试通过，并通过代码检查

### Commit type 类别
feat：表示是一个新功能（feature)
hotfix：hotfix，修补bug
docs：改动、增加文档
opt：修改代码风格及opt imports这些，不改动原有执行的代码
test：增加测试

### 分支命名规范
新功能：feat: feat_flink版本_issuseid
eg: feat_1.12_11111
bug修复： hotfix: hotfix_flink版本_issuseid
eg: hotfix_1.12_11112
注意当前chengying版本依赖flink 版本上进行开发,比如1.12_release 就是对应的flink 1.12 版本；
所以在提交分支的时候请添加上对应的版本