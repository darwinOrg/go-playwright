create table common_task
(
    `id`              bigint auto_increment primary key,
    `type`            varchar(64) NOT NULL DEFAULT '' COMMENT '任务类型',
    `channel`         varchar(64) NOT NULL DEFAULT '' COMMENT '渠道',
    `content`         text COMMENT '任务内容',
    `status`          tinyint     NOT NULL DEFAULT 0 COMMENT '状态（0：初始化，1：处理中，2：成功，3：失败）',
    `reason`          varchar(64) COMMENT '失败原因',
    `start_at`        datetime(3) COMMENT '开始时间',
    `end_at`          datetime(3) COMMENT '结束时间',
    `priority`        int         NOT NULL default 100 COMMENT '优先级',
    `processed_count` int         NOT NULL default 0 COMMENT '已处理次数',
    `locked_at`       datetime(3) COMMENT '锁定时间',
    `lock_until`      datetime(3) COMMENT '锁定到',
    `locked_by`       varchar(64) COMMENT '锁定者',
    `created_at`      datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `modified_at`     datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)
) engine = innodb
  character set = utf8mb4 comment '通用任务';

create table common_kv
(
    `id`          bigint auto_increment primary key,
    `category`    varchar(32)  NOT NULL DEFAULT '' COMMENT '类别',
    `k`           varchar(128) NOT NULL DEFAULT '' COMMENT 'key',
    `v`           varchar(255) NOT NULL COMMENT 'value',
    `expires_at`  datetime(3),
    `created_at`  datetime(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `modified_at` datetime(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    unique key uk_k (`k`)
) engine = innodb
  character set = utf8mb4 comment '通用kv';

create table common_account
(
    `id`          bigint auto_increment primary key,
    `category`    varchar(32)  NOT NULL DEFAULT '' COMMENT '类别',
    `account`     varchar(128) NOT NULL DEFAULT '' COMMENT '账号',
    `password`    varchar(128)          DEFAULT '' COMMENT '密码',
    `mark`        varchar(128)          DEFAULT '' COMMENT '标记',
    `day_limit`   int                   DEFAULT '100' COMMENT '日限额',
    `created_at`  datetime(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `modified_at` datetime(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    unique key uk_category_account (`category`, `account`)
) engine = innodb
  character set = utf8mb4 comment '通用账号';
