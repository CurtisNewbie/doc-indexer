CREATE DATABASE IF NOT EXISTS docindexer;

CREATE TABLE docindexer.bookmark (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `user_no` varchar(32) NOT NULL COMMENT 'user no',
  `icon` text COMMENT 'icon encoded blob',
  `name` varchar(512) NOT NULL DEFAULT '' COMMENT 'bookmark name',
  `href` varchar(1024) NOT NULL DEFAULT '' COMMENT 'bookmark href',
  `md5` varchar(32) not null default '' comment 'md5',
  PRIMARY KEY (`id`),
  KEY `idx_user_no` (`user_no`),
  KEY `idx_name` (name),
  UNIQUE `uk_md5` (md5)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Bookmark';