-- Create syntax for TABLE 'sus_short'
CREATE TABLE `sus_short`
(
    `id`            int(11)     NOT NULL AUTO_INCREMENT,
    `short`         varchar(16) NOT NULL COMMENT '短链接',
    `origin`        tinytext    NOT NULL COMMENT '默认源地址',
    `hash`          char(32)    NOT NULL COMMENT '源地址哈希',
    `is_active`     tinyint(1)  NOT NULL DEFAULT 1 COMMENT '是否开启',
    `is_statistics` tinyint(1)  NOT NULL DEFAULT 0 COMMENT '是否记录PV/UV',
    `is_access_log` tinyint(1)  NOT NULL DEFAULT 1 COMMENT '是否记录AccessLog',
    `activity_id`   int(11)     NOT NULL DEFAULT 0 COMMENT '归属活动',
    `created_at`    datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`    datetime    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uni_short` (`short`),
    KEY `idx_activity_id` (`activity_id`),
    KEY `idx_hash` (`hash`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8 COMMENT ='短链接';

-- Create syntax for TABLE 'sus_statistic'
CREATE TABLE `sus_statistic`
(
    `id`  int(11) NOT NULL AUTO_INCREMENT,
    `sid` int(11) NOT NULL COMMENT '短链接ID（如果为0表示当日所有）',
    `day` int(11) NOT NULL COMMENT '日期',
    `pv`  int(11) NOT NULL DEFAULT 0 COMMENT 'PV',
    `uv`  int(11) NOT NULL DEFAULT 0 COMMENT 'UV',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uni_sid_day` (`sid`, `day`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8 COMMENT ='短链接统计';

-- Create syntax for TABLE 'sus_config'
CREATE TABLE `sus_config`
(
    `name`       varchar(16)  NOT NULL DEFAULT '' COMMENT '配置名',
    `note`       varchar(255) NOT NULL DEFAULT '' COMMENT '配置备注',
    `data`       text         NOT NULL COMMENT '配置JSON字符串',
    `schema`     text         NOT NULL COMMENT '配置schema',
    `created_at` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`name`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8 COMMENT ='配置表';
