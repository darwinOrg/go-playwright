update `scraped_choice_question` set import_status = 1 where source_module = 2 and question_content NOT REGEXP '<[^>]+>';
update `scraped_interview_question` set import_status = 1 where source_module = 2 and CHAR_LENGTH(question_title) > 4 and question_content NOT REGEXP '<[^>]+>';
update `scraped_interview_question` set import_status = 0 where source_module = 2 and question_title like 'Please repeat the sentence you heard%';
# update `scraped_interview_question` set import_status = 1 where CHAR_LENGTH(question_title) > 4 and (question_title not like '%...' or question_content NOT REGEXP '<[^>]+>');

# update scraped_interview_question set question_title = trim(question_title) where question_title;
# update scraped_interview_question set question_title = replace(question_title, '\n', '') where question_title like '%\n%';

alter table scraped_interview_question add `question_id` bigint NOT NULL DEFAULT 0 COMMENT '问题id' after `paper_name`;
alter table scraped_interview_question add `question_content`     text         NOT NULL COMMENT '问题内容' after `question_title`;
alter table scraped_interview_question add `level`                tinyint      NOT NULL DEFAULT '1' COMMENT '难度级别' after `company_name`;
alter table scraped_interview_question add `tags`                 json COMMENT '考察点标签' after `level`;
alter table scraped_interview_question add `skill_tags`           json COMMENT '技能点' after `tags`;

alter table scraped_interview_question add `source_module`        tinyint(1)   NOT NULL DEFAULT 4 COMMENT '来源模块，0-笔试，1-专项练习, 2-题库中心，3-面试题库, 4-面试试卷' after `source_url`;
alter table scraped_choice_question modify `source_module`        tinyint(1)   NOT NULL DEFAULT 0 COMMENT '来源模块，0-笔试，1-专项练习, 2-题库中心，3-面试题库, 4-面试试卷';

update `scraped_interview_question` set question_content = question_title where source_module = 4;
alter table scraped_interview_question drop key `uk_paper_id_question_title_source_channel`;

delete from scraped_interview_question where source_module in (2,3);
delete from scraped_choice_question where source_module in (2,3);

truncate scraped_interview_question;
truncate scraper.scraped_interview_question;
insert into scraped_interview_question select * from scraper.scraped_interview_question;

truncate scraped_choice_question;
truncate scraper.scraped_choice_question;
insert into scraped_choice_question select * from scraper.scraped_choice_question;

alter table scraped_company add `ext_fields`         json COMMENT '扩展字段' after `work_time`;

update common_task t set t.content = json_set(t.content, '$.id', (select ws.id from world_school ws where ws.detail_url = t.content -> '$.detailUrl')) where t.type = 'scrape_world_school_detail';

alter table world_school
    add `logo`         varchar(1024) COMMENT 'logo' after detail_url,
    add `upload_logo`  varchar(1024) COMMENT '上传logo' after logo;

alter table scraped_company add `raw_logo`           varchar(1024) NOT NULL DEFAULT '' COMMENT '原始logo' after `detail_url`;
