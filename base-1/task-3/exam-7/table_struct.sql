-- 执行本脚本前，请先按顺序执行：exam5的CodeFirst操作 ==》base-1\task-3\exam-6\data.sql

-- 为 users 表添加 post_count 字段
ALTER TABLE users 
ADD COLUMN post_count INT NOT NULL DEFAULT 0 COMMENT '用户的文章数量统计';

-- 为 posts 表添加 comment_status 字段
ALTER TABLE posts 
ADD COLUMN comment_status VARCHAR(255) NOT NULL COMMENT '文章的评论状态';

-- 为 comments 表添加 user_id 字段
ALTER TABLE comments 
ADD COLUMN user_id INT NOT NULL DEFAULT 1 COMMENT '文章评论的发布人';