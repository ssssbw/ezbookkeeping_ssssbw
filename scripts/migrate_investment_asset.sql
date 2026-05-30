-- 数据迁移脚本：从 InvestmentAsset 迁移到 Asset + UserAsset
-- 使用前请备份数据库！

-- 1. 创建 Asset 表（如果不存在）
CREATE TABLE IF NOT EXISTS `asset` (
    `asset_id` BIGINT NOT NULL,
    `code` VARCHAR(20) NOT NULL,
    `market` INT NOT NULL,
    `name` VARCHAR(64) NOT NULL,
    `category` VARCHAR(20) NOT NULL,
    `currency` VARCHAR(3) NOT NULL,
    `industry` VARCHAR(20) DEFAULT '',
    `tags` TEXT,
    `extra_info` TEXT,
    `created_unix_time` BIGINT,
    `updated_unix_time` BIGINT,
    PRIMARY KEY (`asset_id`),
    INDEX `IDX_asset_code_market` (`code`, `market`),
    INDEX `IDX_asset_category` (`category`),
    INDEX `IDX_asset_industry` (`industry`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 2. 创建 UserAsset 表（如果不存在）
CREATE TABLE IF NOT EXISTS `user_asset` (
    `id` BIGINT NOT NULL,
    `uid` BIGINT NOT NULL,
    `asset_id` BIGINT NOT NULL,
    `is_active` TINYINT(1) NOT NULL DEFAULT 1,
    `added_unix_time` BIGINT,
    PRIMARY KEY (`id`),
    INDEX `IDX_user_asset_uid_asset_id` (`uid`, `asset_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 3. 从 InvestmentAsset 迁移数据到 Asset（去重）
INSERT IGNORE INTO `asset` (`asset_id`, `code`, `market`, `name`, `category`, `currency`, `industry`, `tags`, `extra_info`, `created_unix_time`, `updated_unix_time`)
SELECT 
    MIN(`asset_id`) as `asset_id`,
    `code`,
    `market`,
    `name`,
    CASE 
        WHEN `type` = 1 THEN 'equity'
        WHEN `type` = 2 THEN 'equity'
        WHEN `type` = 3 THEN 'equity'
        WHEN `type` = 4 THEN 'fixed_income'
        WHEN `type` = 5 THEN 'digital'
        ELSE 'equity'
    END as `category`,
    `currency`,
    '' as `industry`,
    '[]' as `tags`,
    `extra_info`,
    MIN(`created_unix_time`) as `created_unix_time`,
    MAX(`updated_unix_time`) as `updated_unix_time`
FROM `investment_asset`
WHERE `deleted` = 0
GROUP BY `code`, `market`, `name`, `type`, `currency`;

-- 4. 创建 UserAsset 关联
INSERT IGNORE INTO `user_asset` (`id`, `uid`, `asset_id`, `is_active`, `added_unix_time`)
SELECT 
    ia.`asset_id` as `id`,
    ia.`uid`,
    a.`asset_id`,
    ia.`is_active`,
    ia.`created_unix_time`
FROM `investment_asset` ia
JOIN `asset` a ON ia.`code` = a.`code` AND ia.`market` = a.`market`
WHERE ia.`deleted` = 0;

-- 5. 验证迁移结果
SELECT 
    (SELECT COUNT(*) FROM `investment_asset` WHERE `deleted` = 0) as `old_count`,
    (SELECT COUNT(*) FROM `asset`) as `new_asset_count`,
    (SELECT COUNT(*) FROM `user_asset`) as `new_user_asset_count`;
