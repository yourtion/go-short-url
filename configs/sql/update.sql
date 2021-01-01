ALTER TABLE `sus_short`
    ADD `is_statistics` TINYINT(1) NOT NULL DEFAULT '0' COMMENT '是否记录PV/UV' AFTER `is_active`;
ALTER TABLE `sus_short`
    ADD `is_access_log` TINYINT(1) NOT NULL DEFAULT '1' COMMENT '是否记录AccessLog' AFTER `is_statistics`;
