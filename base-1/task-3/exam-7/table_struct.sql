-- 执行本脚本前，请先按顺序执行：exam5的CodeFirst操作 ==》base-1\task-3\exam-6\data.sql

-- 为 users 表添加 post_count 字段
ALTER TABLE users 
ADD COLUMN post_count INT NOT NULL DEFAULT 0 COMMENT '用户的文章数量统计';

