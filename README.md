# PixelCard

Pixelcard 模块总览
=
后台功能模块
=
- Dashboard 仪表盘
  - 总充值数据， 总支出数据，nagtive Balance warning(花到负数)
- KYC
  - 查询用户信息, 审核用户状态
- 用户列表
  - 列表展示，用户信息展示， 冻结，激活， RO登录
- 交易记录
  - 所有用户的交易记录查询
  - negative 余额记录，查询
  - 提现记录
- card center 卡中心
  - card bin 管理
  - card management 卡片管理
  - card holder 持卡人信息列表
- configue 配置
  - card fee 卡费用设置（全局默认配置， 针对用户设置， 月费设置）
  - inbound fee 入金费用设置（全局默认配置， 针对用户设置）
  - inbound/ deposit settings (出金，入金调整设置)
- SMS Code 
  - 手机验证码历史记录

数据库设计
=
#### clients 用户表
|字段名|字段类型|是否为空|备注|
|---  |-------|------|----|
|id|int|no|用户ID，主键
|client_id|int|yes|客户ID，唯一索引|
|client_no|int|yes|客户编码|
|client_email|varchar(128)|yes|邮箱|
|client_name|varchar(128)|yes|客户名|
|client_en_name|varchar(128)|yes|英文名|
|client_type|varchar(64)|yes|Individual/Enterprise 个人/企业|
|client_location|varchar(64)|yes| 国家|
|client_status| varchar(64)| yes | Active/Review/Suspend 激活/审核中/冻结|
|client_review_status|varchar(64)|yes| Unreview/Reviewing/Completed 未审核/审核中/审核结束|
|client_nick_name|varchar|yes|客户别名|
|account_manager|int|no|客户经理|
|register_source|varchar(64)|yes|注册渠道
|contact_type|varchar|yes|联系方式
|contact|varchar|yes|联系号码/id|
|client_register_time|yes|注册时间|
|remark|text|yes|备注
|created_at| datetime|yes|创建时间
|updated_at|datetime|yes|更新时间（自动时间戳）





#### client_due_diligence 用户尽职调查

| 字段名                               | 字段类型         | 是否为空 | 备注             |
| --------------------------------- | ------------ | ---- | -------------- |
| ent\_enterprise\_type             | varchar(64)  | yes  | 企业类型           |
| ent\_enterprise\_chinese\_name    | varchar(128) | no   | 企业中文名          |
| ent\_enterprise\_english\_name    | varchar(128) | yes  | 企业英文名          |
| ent\_business\_registration\_form | varchar(255) | yes  | 营业注册表（附件链接）    |
| ent\_business\_registration\_no   | varchar(64)  | no   | 营业登记编号         |
| ent\_business\_address\_proof     | varchar(255) | yes  | 地址证明           |
| ent\_date\_of\_establishment      | date         | yes  | 成立日期           |
| ent\_date\_of\_expiration         | date         | yes  | 营业有效期          |
| ent\_local\_business\_premise     | varchar(255) | no   | 实体经营地址         |
| ent\_province                     | varchar(64)  | yes  | 所在省份           |
| ent\_city                         | varchar(64)  | yes  | 城市             |
| ent\_listed\_company              | boolean      | yes  | 是否上市公司         |
| ent\_state\_owned                 | boolean      | yes  | 是否国有           |
| ent\_foreign\_invested            | boolean      | yes  | 是否外资           |
| ent\_shareholder\_structure       | varchar(255) | yes  | 股东结构           |
| ent\_is\_b2b                      | boolean      | yes  | 是否 B2B         |
| ent\_operation\_address           | varchar(255) | yes  | 实际经营地址         |
| ent\_registered\_capital          | varchar(64)  | yes  | 注册资本           |
| ent\_intended\_business\_industry | varchar(128) | yes  | 预期业务领域         |
| ind\_country\_or\_region          | varchar(64)  | no   | 国家或地区          |
| ind\_position                     | varchar(64)  | yes  | 职位（如 Director） |
| ind\_id\_type                     | varchar(16)  | no   | 证件类型（如 PP）     |
| ind\_chinese\_name                | varchar(128) | no   | 中文姓名           |
| ind\_english\_name                | varchar(128) | yes  | 英文姓名           |
| ind\_id\_front\_end               | varchar(255) | yes  | 证件正面（附件）       |
| ind\_id\_back\_end                | varchar(255) | yes  | 证件背面（附件）       |
| ind\_identification\_no           | varchar(64)  | no   | 身份证号码/护照号      |
| ind\_issue\_date                  | date         | yes  | 证件签发日期         |
| ind\_expiration\_date             | date         | yes  | 证件有效期          |
| ind\_date\_of\_birth              | date         | yes  | 出生日期           |
| ind\_province\_or\_state          | varchar(64)  | yes  | 所在省份或州         |
| ind\_city                         | varchar(64)  | yes  | 城市             |
| ind\_residential\_address         | varchar(255) | yes  | 居住地址           |
| ind\_reliability\_of\_documents   | varchar(128) | yes  | 文件可信度评价        |



### 交易记录 card_transactions_records

| 字段名                   | 类型            | 是否为空     | 说明              |
| --------------------- | ------------- | -------- | --------------- |
| partner\_order\_id    | varchar(64)   | NOT NULL | 用户请求ID          |
| card\_id              | varchar(64)   | NOT NULL | 卡ID             |
| transaction\_id       | varchar(64)   | NOT NULL | 交易流水号           |
| transaction\_time     | bigint        | NOT NULL | 交易时间（时间戳）       |
| create\_time          | bigint        | NOT NULL | 创建时间（时间戳）       |
| transaction\_currency | varchar(16)   | NOT NULL | 交易币种            |
| transaction\_amount   | decimal(18,6) | NOT NULL | 交易金额            |
| billing\_currency     | varchar(16)   | NOT NULL | 账单币种            |
| billing\_amount       | decimal(18,6) | NOT NULL | 账单金额            |
| auth\_code            | varchar(64)   | YES      | 授权码             |
| transaction\_type     | varchar(32)   | NOT NULL | 交易类型            |
| transaction\_status   | varchar(32)   | NOT NULL | 交易状态            |
| result\_code          | varchar(32)   | NOT NULL | 交易结果 code       |
| fail\_reason          | varchar(255)  | YES      | 失败原因            |
| merchant\_name        | varchar(255)  | YES      | 商户名称            |
| reference\_id         | varchar(64)   | YES      | 关联ID            |
| mcc                   | varchar(16)   | YES      | 商户类别代码          |
| cross\_board\_type    | char(1)       | YES      | 是否跨境：0 否 / 1 是  |
| fund\_account\_type   | varchar(32)   | YES      | 资金账户类型          |
| fund\_direct          | int           | YES      | 资金方向（入账/出账）     |
| merchant\_fee         | JSON          | YES      | 原始费率信息（JSON 格式） |
| fee\_amount           | decimal(18,6) | YES      | 我方手续费金额         |
| fee\_currency         | varchar(16)   | YES      | 我方手续费币种         |
| created_at	          | datetime	    | no	     |本地记录首次创建时间（入库时间）
| updated_at            | datetime      | YES      ｜本地记录最后更新时间 |

### wallet_withdraw

| 字段名              | 类型（建议）   | 是否为空 | 描述说明                      |
| ---------------- | -------- | ---- | ------------------------- |
| `id`             | int      | 否    | 自增主键             |
| `order_id`       | string   | 否    | 提现请求订单号，唯一标识              |
| `created_time`   | datetime | 否    | 提现申请创建时间                  |
| `finish_time`    | datetime | 是    | 提现处理完成时间                  |
| `unit_id`        | string   | 否    | 用户账户标识（系统唯一 ID）           |
| `unit_nickname`  | string   | 否    | 用户昵称/邮箱                   |
| `account_type`   | string   | 否    | 提现账户类型（如：TronAddress 等）   |
| `account_number` | string   | 否    | 提现目标账号（如钱包地址、银行账号）        |
| `amount`         | decimal  | 否    | 提现金额                      |
| `currency`       | string   | 否    | 币种（如 USD）                 |
| `status`         | string   | 否    | 提现状态（Success / Failure）   |
| `operator`       | string   | 是    | 操作人（人工审核处理者）              |
| `action_status`  | string   | 是    | 审核操作动作（Proceed / Decline） |
| `remarks`        | text     | 是    | 审核备注                      |
| `created_at`     | datetime | 否    | 本地记录创建时间（一般为入库时间）         |
| `updated_at`     | datetime | 否    | 本地记录最后更新时间（可自动更新）         |


### card_bin

| 字段名       | 类型（建议）        | 是否为空 | 描述说明                              |
| ------------------------------- | ------------- | ---- | --------------------------------- |
| `card_bin_id`                   | string        | 否    | 卡 BIN 配置唯一标识 ID（如 43612002）       |
| `card_bin`                      | string        | 否    | 实际 BIN 段（如 436120）                |
| `channel_card_bin_id`           | string        | 是    | 通道端卡 BIN ID（如有对接通道系统）             |
| `card_brand`                    | string        | 否    | 卡品牌，如 Visa、MasterCard             |
| `card_type`                     | string        | 否    | 卡类型，如 Virtual、Physical            |
| `card_model`                    | string        | 是    | 卡模型类型，如 SHARE                     |
| `card_type_level`               | string        | 是    | 卡等级类型，如 Virtual                   |
| `currency`                      | string        | 否    | 币种，如 USD                          |
| `region`                        | string        | 否    | 区域，如 US                           |
| `channel`                       | string        | 否    | 通道名称，如 Slash、VoYaPay              |
| `qty_issuance_limit_cardbin`    | int           | 否    | 每个 BIN 的发卡总限额                     |
| `qty_issuance_limit_cardholder` | int           | 否    | 每个持卡人可拥有的卡数量限制                    |
| `create_recharge_limit`         | int           | 否    | 创建卡时的初始充值金额上限                     |
| `auth_amount_limit`             | decimal       | 否    | 授权金额限制（如 1 表示最多授权 1 美元）           |
| `min_balance`                   | decimal       | 否    | 最低余额门槛（如卡不可低于该金额）                 |
| `description`                   | text          | 是    | 备注信息                              |
| `support_platform`              | text          | 是    | 支持的平台说明（如 Facebook 广告、Paypal 绑定等） |
| `issuer_available`              | boolean (Y/N) | 否    | 是否允许发卡                            |
| `top_up`                        | boolean (Y/N) | 否    | 是否允许充值                            |
| `customer_available`            | boolean (Y/N) | 否    | 是否对最终客户可见                         |
| `cardholder_required`           | boolean (Y/N) | 否    | 是否要求持卡人实名信息                       |
| `bin_status`                    | boolean (Y/N) | 否    | BIN 启用状态                          |
| `cancel_card`                   | boolean (Y/N) | 否    | 是否允许销卡                            |
| `withdrawal`                    | boolean (Y/N) | 否    | 是否允许提现                            |
| `support_freezing`              | boolean (Y/N) | 否    | 是否支持冻结操作                          |
| `channel_auto_cancel`           | boolean (Y/N) | 否    | 是否启用通道端自动销卡                       |
| `created_at`                    | datetime      | 否    | 配置记录创建时间                          |
| `updated_at`                    | datetime      | 否    | 配置记录最后更新时间（支持自动更新）                |

### card_bin_global_fee_config

| 字段名                | 类型（建议）          | 是否为空 | 描述说明                                           |
| ------------------ | --------------- | ---- | ---------------------------------------------- |
| `id`               | string / bigint | 否    | 配置唯一标识                                         |
| `fee_type`         | string          | 否    | 交易类型（如：RefundTransaction、AuthorizationQuery 等） |
| `card_bin`         | string          | 否    | 所匹配卡 BIN（如为 ALL 表示适用于所有 BIN）                   |
| `card_model`       | string          | 是    | 卡模型（如 SHARE、ALL）                               |
| `card_type`        | string          | 否    | 卡类型（如 Virtual、Physical）                        |
| `fee_rate`         | decimal(5,2)    | 是    | 按百分比收费（如 2.00 表示 2%）                           |
| `fee_fix`          | decimal(10,2)   | 是    | 固定收费金额（如 1.00 表示固定收 1 美元）                      |
| `fee_currency`     | string(8)       | 否    | 收费币种（如 USD）                                    |
| `active_status`    | enum            | 否    | 激活状态（ENABLE / DISABLE / EXPIRED）               |
| `description`      | text            | 是    | 配置说明                                           |
| `created_at`       | datetime        | 否    | 创建时间                                           |
| `updated_at`       | datetime        | 否    | 最后更新时间                                         |

 
### card_bin_user_fee_config

| 字段名                | 类型（建议）          | 是否为空 | 描述说明                                           |
| ------------------ | --------------- | ---- | ---------------------------------------------- |
| `id`               | string / bigint | 否    | 配置唯一标识                                         |
| `client_id`        | varchar         | 否    | id|
| `fee_type`         | string          | 否    | 交易类型（如：RefundTransaction、AuthorizationQuery 等） |
| `card_bin`         | string          | 否    | 所匹配卡 BIN（如为 ALL 表示适用于所有 BIN）                   |
| `card_model`       | string          | 是    | 卡模型（如 SHARE、ALL）                               |
| `card_type`        | string          | 否    | 卡类型（如 Virtual、Physical）                        |
| `fee_rate_percent` | decimal(5,2)    | 是    | 按百分比收费（如 2.00 表示 2%）                           |
| `fee_fix_amount`   | decimal(10,2)   | 是    | 固定收费金额（如 1.00 表示固定收 1 美元）                      |
| `fee_currency`     | string(8)       | 否    | 收费币种（如 USD）                                    |
| `active_status`    | enum            | 否    | 激活状态（ENABLE / DISABLE / EXPIRED）               |
| `description`      | text            | 是    | 配置说明                                           |
| `created_at`       | datetime        | 否    | 创建时间                                           |
| `updated_at`       | datetime        | 否    | 最后更新时间                                         |

 ### card_bin_monthly_fee



### cards

| 字段名                 | 类型（建议）   | 是否为空 | 描述说明                              |
| ------------------- | -------- | ---- | --------------------------------- |
| `user_id`           | string   | 否    | 用户 ID（发卡用户）                       |
| `card_id`           | string   | 否    | 系统内部卡唯一标识                         |
| `card_no`           | string   | 否    | 卡号 |
| `channel_card_id`   | string   | 是    | 渠道返回的卡 ID（如有）|
| `request_status`    | string   | 否    | 卡请求状态（如 激活、失败、处理中） |
| `request_time`      | datetime | 否    | 发卡请求时间  |
| `card_status`       | string   | 否    | 卡状态（如 Active、Inactive、Terminated） |
| `cardholder`        | string   | 是    | 卡持有人（可为空）                         |
| `channel`           | string   | 否    | 通道名称，如 VoYaPay                    |
| `brand`             | string   | 否    | 卡品牌，如 MasterCard                  |
| `currency`          | string   | 否    | 币种，如 USD                          |
| `available_balance` | decimal  | 否    | 卡片当前可用余额                          |
| `valid_date`        | date     | 否    | 生效日期（卡启用开始）                       |
| `expire_date`       | date     | 否    | 到期日期（卡自然失效时间）                     |
| `termination_date`  | date     | 是    | 实际终止时间（如提前注销）                     |
| `created_at`        | datetime | 否    | 本记录创建时间                           |
| `updated_at`        | datetime | 否    | 本记录最后更新时间                         |

### card_holder

| 字段名                 | 类型（建议）            | 是否为空 | 描述说明                 |
| ------------------- | ----------------- | ---- | -------------------- |
| `client_id`         | string            | 否    | 客户id |
| `card_holder_id`    | string            | 否    | 持卡人 ID（平台生成，内部唯一）    |
| `region`            | string(3)         | 否    | 申请地区，3位字母码（如 US）     |
| `first_name`        | string            | 否    | 姓                    |
| `last_name`         | string            | 否    | 名                    |
| `email`             | string            | 否    | 持卡人邮箱                |
| `mobile_prefix`     | string            | 否    | 手机号国家前缀，如 +86        |
| `mobile`            | string            | 否    | 手机号码                 |
| `birth_date`        | date (yyyy-MM-dd) | 否    | 出生日期                 |
| `country_code`      | string(3)         | 否    | 居住国家代码（3字母 ISO）      |
| `state`             | string            | 是    | 所在州/省                |
| `city`              | string            | 是    | 所在城市                 |
| `postcode`          | string            | 是    | 邮政编码                 |
| `address`           | string            | 是    | 居住地址（详细地址）           |
| `created_at`        | datetime          | 否    | 记录创建时间               |
| `updated_at`        | datetime          | 否    | 记录最后更新时间             |


### inbound_fee_config

| 字段名             | 类型（建议）        | 是否为空 | 描述说明                                     |
| --------------- | ------------- | ---- | ---------------------------------------- |
| `id`            | bigint        | 否    | 主键，自增或唯一标识                               |
| `business_type` | string        | 否    | 业务类型（如 A2P、Direct Payment、PyTransfer）    |
| `currency`      | string(8)     | 否    | 币种（如 USD、VND）                            |
| `fee_rate`      | decimal(5,4)  | 否    | 百分比收费比例（如 0.0100 表示 1%）                  |
| `fix_fee`       | decimal(10,2) | 否    | 固定收费金额                                   |
| `min_fee`       | decimal(10,2) | 是    | 最小手续费（如有）                                |
| `max_fee`       | decimal(10,2) | 是    | 最高手续费（如有）                                |
| `comment`       | text          | 是    | 备注说明                                     |
| `operator`      | string        | 否    | 配置操作者（如 admin、Aaron）                     |
| `status`        | enum          | 否    | 状态：如 `Pending` / `Approved` / `Rejected` |
| `created_at`    | datetime      | 否    | 配置创建时间                                   |
| `updated_at`    | datetime      | 否    | 配置更新时间                                   |

### inbound_deposit_management

| 字段名                            | 类型            | 描述                      |
| ------------------------------ | ------------- | ----------------------- |
| id                             | BIGINT        | 主键，自增                   |
| order\_id                      | VARCHAR(64)   | 系统内部入金订单号               |
| client\_id                     | VARCHAR(64)   | 客户ID                    |
| unit\_id                       | VARCHAR(64)   | 单元ID                    |
| unit\_nickname                 | VARCHAR(64)   | 单元昵称                    |
| customer\_create\_time         | DATETIME      | 客户创建时间                  |
| type                           | VARCHAR(32)   | 入金类型（如 PY / Chain）      |
| channel\_inbound\_id           | VARCHAR(64)   | 通道入金ID                  |
| receive\_account\_name         | VARCHAR(128)  | 接收账户名                   |
| receive\_account\_no           | VARCHAR(128)  | 接收账号                    |
| receive\_account\_address      | VARCHAR(255)  | 接收账户地址                  |
| remit\_bank\_name              | VARCHAR(128)  | 汇款银行名称                  |
| remit\_bank\_account           | VARCHAR(128)  | 汇款银行账号                  |
| remit\_time                    | DATETIME      | 汇款时间                    |
| remit\_reference               | VARCHAR(128)  | 汇款参考号                   |
| fee\_rate                      | DECIMAL(10,4) | 手续费率（如0.005）            |
| fix\_fee                       | DECIMAL(18,2) | 固定费用                    |
| final\_fee                     | DECIMAL(18,2) | 实际费用                    |
| original\_deposit\_amount      | DECIMAL(18,6) | 原始入金金额                  |
| remit\_amount                  | DECIMAL(18,6) | 汇款金额                    |
| final\_deposit\_amount         | DECIMAL(18,6) | 最终到账金额                  |
| associated\_channel\_order\_id | VARCHAR(64)   | 关联通道订单号                 |
| currency                       | VARCHAR(16)   | 币种（如 USD, USDT）         |
| status                         | VARCHAR(32)   | 状态（如 Success / Pending） |
| comment                        | VARCHAR(255)  | 备注                      |
| operator                       | VARCHAR(64)   | 操作人                     |
| created\_at                    | DATETIME      | 创建时间                    |
| updated\_at                    | DATETIME      | 更新时间                    |
