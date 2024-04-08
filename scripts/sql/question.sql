CREATE TABLE `new_coder_interview_paper`
(
    `id`           bigint       NOT NULL AUTO_INCREMENT,
    `job_category` varchar(255) COMMENT '职位类别',
    `job_title`    varchar(255) COMMENT '职位标题',
    `company_name` varchar(255) COMMENT '企业名称',
    `paper_id`     bigint       NOT NULL COMMENT '试卷id',
    `paper_name`   varchar(255) NOT NULL COMMENT '试卷名称',
    `content`      text         NOT NULL COMMENT '内容',
    `created_at`   datetime(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_paper_id` (`paper_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci COMMENT ='牛客面试试卷';

CREATE TABLE `new_coder_interview_question`
(
    `id`                bigint       NOT NULL AUTO_INCREMENT,
    `paper_id`          bigint       NOT NULL COMMENT '试卷id',
    `question_title`    varchar(512) NOT NULL COMMENT '题干',
    `question_solution` text COMMENT '解题思路',
    `url`               varchar(255) COMMENT 'url地址',
    `task_id`           bigint       NOT NULL COMMENT '任务id',
    `created_at`        datetime(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_paper_id_question_title` (`paper_id`, `question_title`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci COMMENT ='牛客面试试题';

CREATE TABLE `new_coder_writing_paper`
(
    `id`               bigint       NOT NULL AUTO_INCREMENT,
    `job_id`           bigint       NOT NULL COMMENT '职位id',
    `job_category`     varchar(255) NOT NULL COMMENT '职位类别',
    `job_title`        varchar(255) NOT NULL COMMENT '职位标题',
    `company_name`     varchar(255) COMMENT '企业名称',
    `paper_id`         bigint       NOT NULL COMMENT '试卷id',
    `paper_name`       varchar(255) NOT NULL COMMENT '试卷名称',
    `url`              varchar(255) NOT NULL COMMENT 'url地址',
    `href`             varchar(255) NOT NULL COMMENT 'href地址',
    `raw_content`      text         NOT NULL COMMENT '原始内容',
    `question_content` longtext COMMENT '问题内容',
    `answer_content`   longtext COMMENT '回答内容',
    `answer_html`      longtext COMMENT '回答html',
    `task_id`          bigint COMMENT '任务id',
    `created_at`       datetime(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `modified_at`      datetime(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_paper_id` (`paper_id`),
    UNIQUE KEY `uk_href` (`href`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci COMMENT ='牛客笔试试卷内容';

CREATE TABLE `new_coder_special_exercise`
(
    `id`              bigint       NOT NULL AUTO_INCREMENT,
    `job_category`    varchar(16)  NOT NULL COMMENT '职位类别',
    `first_category`  varchar(16)  NOT NULL COMMENT '一级类别',
    `second_category` varchar(32)  NOT NULL COMMENT '二级类别',
    `famous_company`  tinyint      NOT NULL default 0 COMMENT '是否名企',
    `url`             varchar(255) NOT NULL COMMENT 'url地址',
    `created_at`      datetime(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_new_coder_special_exercise` (`job_category`, `first_category`, `second_category`, `famous_company`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci COMMENT ='牛客专项练习';

CREATE TABLE `new_coder_special_exercise_paper`
(
    `id`              bigint      NOT NULL AUTO_INCREMENT,
    `exercise_id`     bigint      NOT NULL COMMENT '练习id',
    `job_category`    varchar(16) NOT NULL COMMENT '职位类别',
    `first_category`  varchar(16) NOT NULL COMMENT '一级类别',
    `second_category` varchar(32) NOT NULL COMMENT '二级类别',
    `famous_company`  tinyint     NOT NULL default 0 COMMENT '是否名企',
    `paper_url`       varchar(4096) COMMENT '试卷url地址',
    `paper_content`   longtext COMMENT '试卷内容',
    `answer_url`      varchar(4096) COMMENT '回答url地址',
    `answer_content`  longtext COMMENT '回答内容',
    `answer_html`     longtext COMMENT '回答html',
    `task_id`         bigint COMMENT '任务id',
    `created_at`      datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci COMMENT ='牛客专项练习试卷';

CREATE TABLE `new_coder_online_coding`
(
    `id`               bigint       NOT NULL AUTO_INCREMENT,
    `first_category`   varchar(32)  NOT NULL COMMENT '一级类别',
    `second_category`  varchar(128) NOT NULL COMMENT '二级类别',
    `member`           tinyint      NOT NULL default 0 COMMENT '是否会员',
    `url`              varchar(255) NOT NULL COMMENT 'url地址',
    `page`             int                   default 1 COMMENT '页号',
    `response_type`    varchar(16)  NOT NULL COMMENT '响应类型',
    `response_content` longtext COMMENT '响应内容',
    `created_at`       datetime(3)  NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci COMMENT ='牛客在线编程';

CREATE TABLE `new_coder_online_coding_question`
(
    `id`                   bigint        NOT NULL AUTO_INCREMENT,
    `first_category`       varchar(32)   NOT NULL COMMENT '一级类别',
    `second_category`      varchar(128)  NOT NULL COMMENT '二级类别',
    `member`               tinyint       NOT NULL default 0 COMMENT '是否会员',
    `question_id`          bigint        NOT NULL COMMENT '题目id',
    `question_title`       varchar(1024) NOT NULL COMMENT '问题标题',
    `question_category`    varchar(32) COMMENT '问题类别',
    `difficulty`           int           NOT NULL COMMENT '难度',
    `tips_url`             varchar(1024) COMMENT '答题建议url地址',
    `parent_category_id`   bigint COMMENT '父类型id',
    `parent_category_name` varchar(255) COMMENT '父类型名称',
    `url`                  varchar(1024) NOT NULL COMMENT 'url地址',
    `raw_content`          text COMMENT '原始内容',
    `response_content`     text COMMENT '试题内容',
    `tips_html`            longtext COMMENT '答题建议html',
    `task_id`              bigint COMMENT '任务id',
    `created_at`           datetime(3)   NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `modified_at`          datetime(3)   NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_question_id` (`question_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci COMMENT ='牛客在线编程题目';

CREATE TABLE `new_coder_online_coding_question`
(
    `id`                   bigint        NOT NULL AUTO_INCREMENT,
    `first_category`       varchar(32)   NOT NULL COMMENT '一级类别',
    `second_category`      varchar(128)  NOT NULL COMMENT '二级类别',
    `member`               tinyint       NOT NULL default 0 COMMENT '是否会员',
    `question_id`          bigint        NOT NULL COMMENT '题目id',
    `question_title`       varchar(1024) NOT NULL COMMENT '问题标题',
    `question_category`    varchar(32) COMMENT '问题类别',
    `difficulty`           int           NOT NULL COMMENT '难度',
    `tips_url`             varchar(1024) COMMENT '答题建议url地址',
    `parent_category_id`   bigint COMMENT '父类型id',
    `parent_category_name` varchar(255) COMMENT '父类型名称',
    `url`                  varchar(1024) NOT NULL COMMENT 'url地址',
    `raw_content`          text COMMENT '原始内容',
    `response_content`     text COMMENT '试题内容',
    `tips_html`            longtext COMMENT '答题建议html',
    `task_id`              bigint COMMENT '任务id',
    `created_at`           datetime(3)   NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    `modified_at`          datetime(3)   NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_question_id` (`question_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci COMMENT ='牛客在线编程题目';

CREATE TABLE `new_coder_question_center`
(
    `id`               bigint      NOT NULL AUTO_INCREMENT,
    `page_no`          int         NOT NULL default 1 COMMENT '页号',
    `response_content` longtext    NOT NULL COMMENT '试题内容',
    `created_at`       datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci COMMENT ='牛客题库中心';

CREATE TABLE `new_coder_interview_library`
(
    `id`               bigint      NOT NULL AUTO_INCREMENT,
    `page_no`          int         NOT NULL default 1 COMMENT '页号',
    `response_content` longtext    NOT NULL COMMENT '试题内容',
    `created_at`       datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci COMMENT ='牛客面试题库';
