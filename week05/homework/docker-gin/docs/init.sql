CREATE DATABASE IF NOT EXISTS `smart_vocab` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `smart_vocab`;

-- 用户表
CREATE TABLE `users` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `username` VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名',
    `password_hash` VARCHAR(255) NOT NULL COMMENT '加密后的密码',
    `created_at` DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    `deleted_at` DATETIME(3) DEFAULT NULL COMMENT '软删除时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户认证表';

-- 单词记录表
CREATE TABLE `words` (
    `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '关联用户ID',
    `word` VARCHAR(100) NOT NULL COMMENT '查询的单词',
    `meaning` TEXT NOT NULL COMMENT 'AI生成的精准释义',
    `sentences` JSON NOT NULL COMMENT 'AI生成的3条例句 (JSON数组存储)',
    `ai_provider` VARCHAR(50) NOT NULL COMMENT '使用的AI模型 (如 DeepSeek, Qwen)',
    `created_at` DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3),
    `updated_at` DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    `deleted_at` DATETIME(3) DEFAULT NULL COMMENT '软删除时间',
    INDEX `idx_user_id` (`user_id`),
    INDEX `idx_deleted_at` (`deleted_at`),
    CONSTRAINT `fk_words_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='单词学习记录表';