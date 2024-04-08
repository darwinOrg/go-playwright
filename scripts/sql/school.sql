CREATE TABLE `world_school`
(
    `id`           bigint       NOT NULL AUTO_INCREMENT,
    `cn_name`      varchar(255) NOT NULL COMMENT '中文学校名称',
    `en_name`      varchar(255) NOT NULL COMMENT '英文学校名称',
    `region`       varchar(255) NOT NULL COMMENT '国家/地区',
    `qs`           int          NOT NULL COMMENT 'QS排名',
    `detail_url`   varchar(255) COMMENT '详情url',
    `logo`         varchar(1024) COMMENT 'logo',
    `upload_logo`  varchar(1024) COMMENT '上传logo',
    `abbreviation` varchar(1024) COMMENT '缩写/别名',
    `nature`       varchar(255) COMMENT '学校性质',
    `semester`     varchar(255) COMMENT '学期制度',
    `website`      varchar(255) COMMENT '学校网址',
    `phone`        varchar(255) COMMENT '电话',
    `address`      varchar(255) COMMENT '学校所在地',
    `intro`        text COMMENT '介绍',
    `created_at`   datetime(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci COMMENT ='世界大学排名';

CREATE TABLE `school_major`
(
    `id`                     bigint       NOT NULL AUTO_INCREMENT,
    `school_name`            varchar(255) NOT NULL COMMENT '学校名称',
    `dual_class_majors`      json COMMENT '双一流建设学科',
    `country_special_majors` json COMMENT '国家特色专业',
    `created_at`             datetime(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_school_name` (`school_name`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci COMMENT ='学校专业';

CREATE TABLE `school_logo`
(
    `id`          bigint       NOT NULL AUTO_INCREMENT,
    `school_name` varchar(255) NOT NULL COMMENT '学校名称',
    `logo`        varchar(1024) COMMENT 'logo',
    `created_at`  datetime(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_school_name` (`school_name`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci COMMENT ='学校logo';

CREATE TABLE `gaokao_school`
(
    `id`                     bigint       NOT NULL AUTO_INCREMENT,
    `school_name`            varchar(255) NOT NULL COMMENT '学校名称',
    `logo`                   varchar(1024) COMMENT 'logo',
    `upload_logo`            varchar(1024) COMMENT '上传logo',
    `dual_class_majors`      json COMMENT '双一流建设学科',
    `country_special_majors` json COMMENT '国家特色专业',
    `created_year`           varchar(16) COMMENT '创建年份',
    `created_at`             datetime(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_school_name` (`school_name`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci COMMENT ='高考学校';
