
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for g_counter
-- ----------------------------
DROP TABLE IF EXISTS `g_counter`;
CREATE TABLE `g_counter` (
  `sql_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `id` varchar(512) NOT NULL,
  `gkey` varchar(512) NOT NULL,
  `gtype` varchar(512) NOT NULL,
  `amount` decimal(60,0) NOT NULL,
  PRIMARY KEY (`sql_id`),
  UNIQUE KEY `idx_g_counter_id` (`id`),
  KEY `idx_g_counter_gkey` (`gkey`),
  KEY `idx_g_counter_gtype` (`gtype`),
  KEY `idx_g_counter_amount` (`amount`)
) ENGINE=InnoDB AUTO_INCREMENT=79 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for g_counter_daily_agg
-- ----------------------------
DROP TABLE IF EXISTS `g_counter_daily_agg`;
CREATE TABLE `g_counter_daily_agg` (
  `sql_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `id` varchar(512) NOT NULL,
  `gkey` varchar(512) NOT NULL,
  `gtype` varchar(512) NOT NULL,
  `date` date NOT NULL,
  `amount` decimal(60,0) NOT NULL,
  `status` varchar(32) NOT NULL,
  PRIMARY KEY (`sql_id`),
  UNIQUE KEY `idx_g_counter_daily_agg_id` (`id`),
  KEY `idx_g_counter_daily_agg_date` (`date`),
  KEY `idx_g_counter_daily_agg_gkey` (`gkey`),
  KEY `idx_g_counter_daily_agg_gtype` (`gtype`),
  KEY `idx_g_counter_daily_agg_status` (`status`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=4184 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for g_counter_detail
-- ----------------------------
DROP TABLE IF EXISTS `g_counter_detail`;
CREATE TABLE `g_counter_detail` (
  `sql_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `id` varchar(512) NOT NULL,
  `gkey` varchar(512) NOT NULL,
  `gtype` varchar(512) NOT NULL,
  `datetime` datetime(6) NOT NULL,
  `amount` decimal(60,0) NOT NULL,
  `msg` longtext,
  PRIMARY KEY (`sql_id`),
  UNIQUE KEY `idx_g_counter_detail_id` (`id`),
  KEY `idx_g_counter_detail_gkey` (`gkey`),
  KEY `idx_g_counter_detail_gtype` (`gtype`),
  KEY `idx_g_counter_detail_datetime` (`datetime`)
) ENGINE=InnoDB AUTO_INCREMENT=29 DEFAULT CHARSET=utf8mb4;


SET FOREIGN_KEY_CHECKS = 1;
