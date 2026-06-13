-- ============================================================
-- 投资模块测试数据种子脚本
-- 用户 UID: 3772517497039749120
-- 前置条件: asset 表已有数据
-- 请逐段执行，或使用支持多语句的客户端（如 mysql CLI）
-- ============================================================

-- STEP 1: 初始化变量
SET @uid = 3772517497039749120;
SET @today = 1781254258;
SELECT @max_id := IFNULL(MAX(account_id), 1000000000000000000) FROM account;
SET @parent_account_id = @max_id + 1;
SET @pool1 = @max_id + 2;
SET @pool2 = @max_id + 3;

-- user_asset 和 investment_transaction 的主键 ID（无自增，需手动指定）
SELECT @ua_max := IFNULL(MAX(id), 1000000000000000000) FROM user_asset;
SELECT @tx_max := IFNULL(MAX(transaction_id), 1000000000000000000) FROM investment_transaction;

-- STEP 2: 获取资产 ID
SELECT @asset1 := asset_id FROM asset WHERE code = '000001' AND market = 1 LIMIT 1;
SELECT @asset2 := asset_id FROM asset WHERE code = '000009' AND market = 1 LIMIT 1;
SELECT @asset3 := asset_id FROM asset WHERE code = '000011' AND market = 1 LIMIT 1;
SELECT @asset4 := asset_id FROM asset WHERE code = '000003' AND market = 1 LIMIT 1;
SELECT @asset5 := asset_id FROM asset WHERE code = '000006' AND market = 1 LIMIT 1;
SELECT @asset6 := asset_id FROM asset WHERE code = '000008' AND market = 1 LIMIT 1;

-- STEP 3: 创建账户
INSERT INTO account (account_id, uid, deleted, type, category, name, parent_account_id, icon, color, currency, balance, comment, display_order, hidden, created_unix_time, updated_unix_time)
VALUES (@parent_account_id, @uid, 0, 2, 7, '资金池测试账户', 0, 0, '4CAF50', 'CNY', 0, '测试用投资总资金池', 100, 0, @today, @today);

INSERT INTO account (account_id, uid, deleted, type, category, name, parent_account_id, icon, color, currency, balance, comment, display_order, hidden, created_unix_time, updated_unix_time)
VALUES (@pool1, @uid, 0, 1, 7, '基金池 1 测试', @parent_account_id, 0, '2196F3', 'CNY', 0, '测试用子资金池1', 1, 0, @today, @today);

INSERT INTO account (account_id, uid, deleted, type, category, name, parent_account_id, icon, color, currency, balance, comment, display_order, hidden, created_unix_time, updated_unix_time)
VALUES (@pool2, @uid, 0, 1, 7, '基金池 2 测试', @parent_account_id, 0, 'FF9800', 'CNY', 0, '测试用子资金池2', 2, 0, @today, @today);

-- STEP 4: 插入 user_asset 持仓记录
INSERT INTO user_asset (id, uid, asset_id, deleted, is_active, is_watchlist, comment, created_unix_time, updated_unix_time) VALUES
    (@ua_max + 1, @uid, @asset1, 0, 1, 0, '持仓-华夏成长', @today, @today),
    (@ua_max + 2, @uid, @asset2, 0, 1, 0, '持仓-易方达货币', @today, @today),
    (@ua_max + 3, @uid, @asset3, 0, 1, 0, '持仓-华夏大盘', @today, @today),
    (@ua_max + 4, @uid, @asset6, 0, 1, 0, '持仓-嘉实中证500', @today, @today);

-- STEP 5: 插入 user_asset 自选记录
INSERT INTO user_asset (id, uid, asset_id, deleted, is_active, is_watchlist, comment, created_unix_time, updated_unix_time) VALUES
    (@ua_max + 5, @uid, @asset4, 0, 1, 1, '自选-中海可转债', @today, @today),
    (@ua_max + 6, @uid, @asset5, 0, 1, 1, '自选-西部利得', @today, @today);

-- STEP 6: 行情数据 - 华夏成长混合
INSERT INTO market_data (asset_id, date, price, volume, created_unix_time, updated_unix_time) VALUES
    (@asset1, UNIX_TIMESTAMP(DATE_SUB(CURDATE(), INTERVAL 4 DAY)) * 1000, 12400, 0, @today, @today),
    (@asset1, UNIX_TIMESTAMP(DATE_SUB(CURDATE(), INTERVAL 3 DAY)) * 1000, 12450, 0, @today, @today),
    (@asset1, UNIX_TIMESTAMP(DATE_SUB(CURDATE(), INTERVAL 2 DAY)) * 1000, 12500, 0, @today, @today),
    (@asset1, UNIX_TIMESTAMP(DATE_SUB(CURDATE(), INTERVAL 1 DAY)) * 1000, 12480, 0, @today, @today),
    (@asset1, UNIX_TIMESTAMP(CURDATE()) * 1000, 12520, 0, @today, @today);

-- STEP 7: 行情数据 - 易方达货币
INSERT INTO market_data (asset_id, date, price, volume, created_unix_time, updated_unix_time) VALUES
    (@asset2, UNIX_TIMESTAMP(DATE_SUB(CURDATE(), INTERVAL 4 DAY)) * 1000, 10000, 0, @today, @today),
    (@asset2, UNIX_TIMESTAMP(DATE_SUB(CURDATE(), INTERVAL 3 DAY)) * 1000, 10001, 0, @today, @today),
    (@asset2, UNIX_TIMESTAMP(DATE_SUB(CURDATE(), INTERVAL 2 DAY)) * 1000, 10001, 0, @today, @today),
    (@asset2, UNIX_TIMESTAMP(DATE_SUB(CURDATE(), INTERVAL 1 DAY)) * 1000, 10002, 0, @today, @today),
    (@asset2, UNIX_TIMESTAMP(CURDATE()) * 1000, 10002, 0, @today, @today);

-- STEP 8: 行情数据 - 华夏大盘精选
INSERT INTO market_data (asset_id, date, price, volume, created_unix_time, updated_unix_time) VALUES
    (@asset3, UNIX_TIMESTAMP(DATE_SUB(CURDATE(), INTERVAL 4 DAY)) * 1000, 84000, 0, @today, @today),
    (@asset3, UNIX_TIMESTAMP(DATE_SUB(CURDATE(), INTERVAL 3 DAY)) * 1000, 84500, 0, @today, @today),
    (@asset3, UNIX_TIMESTAMP(DATE_SUB(CURDATE(), INTERVAL 2 DAY)) * 1000, 85000, 0, @today, @today),
    (@asset3, UNIX_TIMESTAMP(DATE_SUB(CURDATE(), INTERVAL 1 DAY)) * 1000, 84800, 0, @today, @today),
    (@asset3, UNIX_TIMESTAMP(CURDATE()) * 1000, 85200, 0, @today, @today);

-- STEP 9: 行情数据 - 嘉实中证500
INSERT INTO market_data (asset_id, date, price, volume, created_unix_time, updated_unix_time) VALUES
    (@asset6, UNIX_TIMESTAMP(DATE_SUB(CURDATE(), INTERVAL 4 DAY)) * 1000, 15600, 0, @today, @today),
    (@asset6, UNIX_TIMESTAMP(DATE_SUB(CURDATE(), INTERVAL 3 DAY)) * 1000, 15700, 0, @today, @today),
    (@asset6, UNIX_TIMESTAMP(DATE_SUB(CURDATE(), INTERVAL 2 DAY)) * 1000, 15750, 0, @today, @today),
    (@asset6, UNIX_TIMESTAMP(DATE_SUB(CURDATE(), INTERVAL 1 DAY)) * 1000, 15800, 0, @today, @today),
    (@asset6, UNIX_TIMESTAMP(CURDATE()) * 1000, 15850, 0, @today, @today);

-- STEP 10: 行情数据 - 自选资产
INSERT INTO market_data (asset_id, date, price, volume, created_unix_time, updated_unix_time) VALUES
    (@asset4, UNIX_TIMESTAMP(CURDATE()) * 1000, 11200, 0, @today, @today),
    (@asset5, UNIX_TIMESTAMP(CURDATE()) * 1000, 15800, 0, @today, @today);

-- STEP 11: 基金池 1 测试 - 交易记录
INSERT INTO investment_transaction (transaction_id, uid, deleted, asset_id, account_id, type, trade_time, confirm_time, quantity, price, amount, fee, related_transaction_id, timezone_utc_offset, comment, created_unix_time, updated_unix_time) VALUES
    (@tx_max + 1, @uid, 0, @asset1, @pool1, 1, UNIX_TIMESTAMP('2025-01-15 09:30:00'), UNIX_TIMESTAMP('2025-01-16 09:00:00'), 100000000, 12000, 120000000, 120000, 0, 480, '池1-买入华夏成长混合', @today, @today),
    (@tx_max + 2, @uid, 0, @asset1, @pool1, 1, UNIX_TIMESTAMP('2025-03-10 14:00:00'), UNIX_TIMESTAMP('2025-03-11 09:00:00'), 50000000, 12500, 62500000, 62500, 0, 480, '池1-加仓华夏成长混合', @today, @today),
    (@tx_max + 3, @uid, 0, @asset2, @pool1, 1, UNIX_TIMESTAMP('2025-02-01 10:00:00'), UNIX_TIMESTAMP('2025-02-02 09:00:00'), 500000000, 10000, 500000000, 0, 0, 480, '池1-买入货币基金', @today, @today),
    (@tx_max + 4, @uid, 0, @asset3, @pool1, 1, UNIX_TIMESTAMP('2025-04-20 11:00:00'), UNIX_TIMESTAMP('2025-04-21 09:00:00'), 20000000, 80000, 160000000, 160000, 0, 480, '池1-买入华夏大盘精选', @today, @today),
    (@tx_max + 5, @uid, 0, @asset1, @pool1, 2, UNIX_TIMESTAMP('2025-06-01 15:00:00'), UNIX_TIMESTAMP('2025-06-02 09:00:00'), 30000000, 12500, 37500000, 37500, 0, 480, '池1-部分卖出华夏成长混合', @today, @today);

-- STEP 12: 基金池 2 测试 - 交易记录
INSERT INTO investment_transaction (transaction_id, uid, deleted, asset_id, account_id, type, trade_time, confirm_time, quantity, price, amount, fee, related_transaction_id, timezone_utc_offset, comment, created_unix_time, updated_unix_time) VALUES
    (@tx_max + 6, @uid, 0, @asset1, @pool2, 1, UNIX_TIMESTAMP('2025-02-10 10:00:00'), UNIX_TIMESTAMP('2025-02-11 09:00:00'), 80000000, 11800, 94400000, 94400, 0, 480, '池2-买入华夏成长混合', @today, @today),
    (@tx_max + 7, @uid, 0, @asset6, @pool2, 1, UNIX_TIMESTAMP('2025-03-05 09:30:00'), UNIX_TIMESTAMP('2025-03-06 09:00:00'), 150000000, 15000, 225000000, 225000, 0, 480, '池2-买入嘉实中证500', @today, @today),
    (@tx_max + 8, @uid, 0, @asset6, @pool2, 1, UNIX_TIMESTAMP('2025-05-12 14:00:00'), UNIX_TIMESTAMP('2025-05-13 09:00:00'), 100000000, 15500, 155000000, 155000, 0, 480, '池2-加仓嘉实中证500', @today, @today),
    (@tx_max + 9, @uid, 0, @asset3, @pool2, 1, UNIX_TIMESTAMP('2025-04-25 11:00:00'), UNIX_TIMESTAMP('2025-04-26 09:00:00'), 15000000, 82000, 123000000, 123000, 0, 480, '池2-买入华夏大盘精选', @today, @today);

-- ============================================================
-- 预期结果:
--
--   华夏成长混合: 总 17000 份
--     ├─ 基金池 1 测试: 12000 份
--     └─ 基金池 2 测试:  8000 份
--
--   易方达货币: 总 50000 份
--     └─ 基金池 1 测试: 50000 份
--
--   华夏大盘精选: 总 3500 份
--     ├─ 基金池 1 测试: 2000 份
--     └─ 基金池 2 测试: 1500 份
--
--   嘉实中证500: 总 25000 份
--     └─ 基金池 2 测试: 25000 份
-- ============================================================
