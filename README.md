# gopressblog
a blog 
## 功能

- 首页
    - 文章列表
    - 模糊搜索文章标题
    - 时间正、反序
    - 分页
    - 今日活跃排行
    
- 文章
    - 发表文章
        - 获得积分
    - 查看文章详情
    - 修改文章
    - 评论文章
        - 回复评论
            - 获得积分
        - 回复时@某人
- 账号
    - 积分信息
    - 头像信息
        - 修改头像
    
- 消息
  - 查看别人@自己的消息
  - 查看系统提示
      - 全部标为已读

- 退出

- migrations
    - 同步数据库结构

- console
    - 提供计划任务
        - 清空每日获得积分
        - 奖励前10名活跃用户
        - 创建 es 索引
        - 同步 es 数据


## 配置文件

- config/config.yaml
    - database: 数据库信息
    - score: 操作积分奖励
    - elastic： elastic配置
    
## 运行
- 1.glide install
- 2.go run main.go (测试）

