CREATE TABLE `scraped_interview_question`
(
    `id`                   bigint       NOT NULL AUTO_INCREMENT,
    `category_name`        varchar(255) NOT NULL DEFAULT '' COMMENT '职类名称',
    `raw_category_name`    varchar(255) NOT NULL DEFAULT '' COMMENT '原始职类名称',
    `job_title`            varchar(255) COMMENT '职位标题',
    `paper_id`             bigint       NOT NULL COMMENT '试卷id',
    `paper_name`           varchar(255) NOT NULL COMMENT '试卷名称',
    `question_id`          bigint       NOT NULL DEFAULT 0 COMMENT '问题id',
    `question_title`       varchar(512) NOT NULL COMMENT '题干',
    `question_content`     text         NOT NULL COMMENT '问题内容',
    `question_solution`    text COMMENT '解题思路',
    `company_name`         varchar(255) COMMENT '企业名称',
    `level`                tinyint      NOT NULL DEFAULT '1' COMMENT '难度级别',
    `tags`                 json COMMENT '考察点标签',
    `skill_tags`           json COMMENT '技能点',
    `source_channel`       tinyint(1)   NOT NULL DEFAULT 1 COMMENT '来源渠道',
    `source_url`           varchar(4096) COMMENT '来源url地址',
    `source_module`        tinyint(1)   NOT NULL DEFAULT 4 COMMENT '来源你模块，0-笔试，1-专项练习, 2-题库中心，3-面试题库, 4-面试试卷',
    `digital_man_media_id` bigint       NOT NULL DEFAULT 0 COMMENT '数字人媒体id',
    `digital_man_task_id`  bigint       NOT NULL DEFAULT 0 COMMENT '数字人任务id',
    `import_status`        tinyint(1)   NOT NULL DEFAULT 0 COMMENT '导入状态，0-初始化，1-待导入，2：已导入',
    `imported_question_id` bigint       NOT NULL DEFAULT 0 COMMENT '导入后的问题id',
    `created_at`           datetime(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `modified_at`          datetime(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_question_id_source_channel` (`question_id`, `source_channel`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci COMMENT ='已爬取的面试题';

CREATE TABLE `scraped_choice_question`
(
    `id`                   bigint        NOT NULL AUTO_INCREMENT,
    `category_name`        varchar(255)  NOT NULL DEFAULT '' COMMENT '职类名称',
    `raw_category_name`    varchar(255)  NOT NULL DEFAULT '' COMMENT '原始职类名称',
    `job_id`               bigint        NOT NULL DEFAULT 0 COMMENT '职位id',
    `job_title`            varchar(255)  NOT NULL DEFAULT '' COMMENT '职位标题',
    `match_job_names`      json COMMENT '匹配职位名称',
    `paper_id`             bigint        NOT NULL DEFAULT 0 COMMENT '试卷id',
    `paper_name`           varchar(255)  NOT NULL COMMENT '试卷名称',
    `question_id`          bigint        NOT NULL COMMENT '问题id',
    `question_title`       varchar(512)  NOT NULL COMMENT '题干',
    `question_content`     text          NOT NULL COMMENT '问题内容',
    `question_style`       tinyint       NOT NULL COMMENT '试题样式',
    `reference_answer`     text COMMENT '引用答案',
    `analysis`             text COMMENT '思路分析',
    `company_name`         varchar(255) COMMENT '企业名称',
    `level`                tinyint       NOT NULL DEFAULT '1' COMMENT '难度级别',
    `tags`                 json COMMENT '考察点标签',
    `skill_tags`           json COMMENT '技能点',
    `knowledge_tags`       json COMMENT '知识点',
    `options`              json          NOT NULL COMMENT '选项',
    `correct_answer`       text          NOT NULL COMMENT '正确答案',
    `source_channel`       tinyint(1)    NOT NULL DEFAULT 1 COMMENT '来源渠道',
    `source_module`        tinyint(1)    NOT NULL DEFAULT 0 COMMENT '来源你模块，0-笔试，1-专项练习, 2-题库中心，3-面试题库, 4-面试试卷',
    `source_url`           varchar(4096) NOT NULL DEFAULT '' COMMENT '来源url地址',
    `import_status`        tinyint(1)    NOT NULL DEFAULT 0 COMMENT '导入状态，0-初始化，1-待导入，2：已导入',
    `imported_question_id` bigint        NOT NULL DEFAULT 0 COMMENT '导入后的问题id',
    `created_at`           datetime(3)   NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `modified_at`          datetime(3)   NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_question_id_source_channel` (`question_id`, `source_channel`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci COMMENT ='已爬取的笔试选择题';
