-- ============================================
-- 修复 Pending 状态交易记录的 SQL 脚本
-- ============================================
-- 说明：此脚本用于处理之前因为代码逻辑问题而遗漏的 Pending 状态交易记录
-- 需要处理的 TransactionID：
-- CS1962479145732902739, CS2006270514421859244, CS2013151030424668983,
-- CS2013151030424669020, CS2013151030428863328, CS2013151030433057754,
-- CS2013151030437252075

-- ============================================
-- 第一步：查询需要处理的交易记录信息（用于验证）
-- ============================================
SELECT 
    ctr.id,
    ctr.transaction_id,
    ctr.order_id,
    ctr.card_id,
    ctr.status,
    ctr.amount,
    ctr.currency,
    ctr.transaction_type,
    ctr.transaction_time,
    ctr.client_id,
    ctr.iam_id,
    pc.card_no,
    pc.client_id as card_client_id,
    pc.card_bin_id,
    w.balance as current_wallet_balance
FROM card_transaction_record ctr
LEFT JOIN client_card pc ON pc.card_id = ctr.card_id
LEFT JOIN wallets w ON w.client_id = pc.client_id
WHERE ctr.transaction_id IN (
    'CS1962479145732902739',
    'CS2006270514421859244',
    'CS2013151030424668983',
    'CS2013151030424669020',
    'CS2013151030428863328',
    'CS2013151030433057754',
    'CS2013151030437252075'
)
AND ctr.status = 'Pending'
AND ctr.transaction_type = 'Card_Withdraw';

-- ============================================
-- 第二步：开始事务处理
-- ============================================
START TRANSACTION;

-- ============================================
-- 第三步：更新交易记录状态和补充信息
-- ============================================
UPDATE card_transaction_record ctr
INNER JOIN client_card pc ON pc.card_id = ctr.card_id
SET 
    ctr.status = 'Success',
    ctr.client_id = pc.client_id,
    ctr.iam_id = pc.iam_id
WHERE ctr.transaction_id IN (
    'CS1962479145732902739',
    'CS2006270514421859244',
    'CS2013151030424668983',
    'CS2013151030424669020',
    'CS2013151030428863328',
    'CS2013151030433057754',
    'CS2013151030437252075'
)
AND ctr.status = 'Pending'
AND ctr.transaction_type = 'Card_Withdraw';

-- ============================================
-- 第四步：更新钱包余额（全额入账，不扣除手续费）
-- ============================================
-- 注意：先按 client_id 分组求和，确保同一客户的多笔交易金额能正确累加
UPDATE wallets w
INNER JOIN (
    SELECT client_id, SUM(amount) as total_amount
    FROM card_transaction_record
    WHERE transaction_id IN (
        'CS1962479145732902739',
        'CS2006270514421859244',
        'CS2013151030424668983',
        'CS2013151030424669020',
        'CS2013151030428863328',
        'CS2013151030433057754',
        'CS2013151030437252075'
    )
    AND status = 'Success'
    AND transaction_type = 'Card_Withdraw'
    GROUP BY client_id
) ctr_sum ON ctr_sum.client_id = w.client_id
SET w.balance = w.balance + ctr_sum.total_amount;

-- ============================================
-- 第五步：插入 WalletHistory 记录（入账记录）
-- ============================================
INSERT INTO wallet_history (
    created_at,
    updated_at,
    deleted_at,
    order_id,
    client_id,
    iam_id,
    is_fee,
    amount,
    amount_currency,
    balance,
    currency,
    transaction_type,
    reference_id,
    card_no
)
SELECT 
    NOW() as created_at,
    NOW() as updated_at,
    NULL as deleted_at,
    ctr.order_id,
    ctr.client_id,
    ctr.iam_id,
    0 as is_fee,
    ctr.amount,
    ctr.currency as amount_currency,
    w.balance as balance,  -- 入账后的余额
    ctr.currency,
    'Card_Withdraw' as transaction_type,
    ctr.order_id as reference_id,
    pc.card_no
FROM card_transaction_record ctr
INNER JOIN client_card pc ON pc.card_id = ctr.card_id
INNER JOIN wallets w ON w.client_id = ctr.client_id
WHERE ctr.transaction_id IN (
    'CS1962479145732902739',
    'CS2006270514421859244',
    'CS2013151030424668983',
    'CS2013151030424669020',
    'CS2013151030428863328',
    'CS2013151030433057754',
    'CS2013151030437252075'
)
AND ctr.status = 'Success'
AND ctr.transaction_type = 'Card_Withdraw'
AND NOT EXISTS (
    SELECT 1 FROM wallet_history wh 
    WHERE wh.reference_id = ctr.order_id 
    AND wh.transaction_type = 'Card_Withdraw'
    AND wh.is_fee = 0
);

-- ============================================
-- 第六步：更新或插入 ClientDailyReport
-- ============================================
INSERT INTO client_daily_report (
    report_day,
    client_id,
    card_withdraw_count,
    card_withdraw_amount,
    fee_amount,
    created_at
)
SELECT 
    DATE(ctr.transaction_time) as report_day,
    ctr.client_id,
    1 as card_withdraw_count,
    ctr.amount as card_withdraw_amount,
    0 as fee_amount,  -- 不收取手续费
    NOW() as created_at
FROM card_transaction_record ctr
WHERE ctr.transaction_id IN (
    'CS1962479145732902739',
    'CS2006270514421859244',
    'CS2013151030424668983',
    'CS2013151030424669020',
    'CS2013151030428863328',
    'CS2013151030433057754',
    'CS2013151030437252075'
)
AND ctr.status = 'Success'
AND ctr.transaction_type = 'Card_Withdraw'
ON DUPLICATE KEY UPDATE
    card_withdraw_count = card_withdraw_count + 1,
    card_withdraw_amount = card_withdraw_amount + VALUES(card_withdraw_amount);

-- ============================================
-- 第八步：验证处理结果
-- ============================================
SELECT 
    ctr.transaction_id,
    ctr.status,
    ctr.client_id,
    w.balance as wallet_balance,
    (SELECT COUNT(*) FROM wallet_history wh WHERE wh.reference_id = ctr.order_id AND wh.transaction_type = 'Card_Withdraw') as history_count
FROM card_transaction_record ctr
LEFT JOIN wallets w ON w.client_id = ctr.client_id
WHERE ctr.transaction_id IN (
    'CS1962479145732902739',
    'CS2006270514421859244',
    'CS2013151030424668983',
    'CS2013151030424669020',
    'CS2013151030428863328',
    'CS2013151030433057754',
    'CS2013151030437252075'
);

-- ============================================
-- 如果验证无误，提交事务；如果有问题，回滚
-- ============================================
-- COMMIT;
-- ROLLBACK;
