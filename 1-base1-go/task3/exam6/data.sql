-- 执行本脚本前，需先执行 exam5 的 CodeFirst 操作

-- 插入用户数据
INSERT INTO users (username, password) VALUES
('张三', 'password1'),
('李四', 'password2'),
('王五', 'password3');

-- 插入帖子数据
INSERT INTO posts (title, content, user_id) VALUES
('如何学习 Go 语言', '学习 Go 语言需要掌握基础语法、并发编程和标准库。', 1),
('数据库设计技巧', '设计数据库时需要注意规范化和性能优化。', 2),
('前端开发入门', '前端开发需要熟悉 HTML、CSS 和 JavaScript。', 3);

-- 插入评论数据
INSERT INTO comments (content, post_id) VALUES
('Go 语言确实很强大！', 1),
('我最近也在学习数据库设计。', 2),
('前端开发很有意思，但也很复杂。', 3),
('感谢分享，非常有帮助！', 1);