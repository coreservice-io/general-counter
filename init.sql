SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for g_counter
-- ----------------------------
DROP TABLE IF EXISTS `g_counter`;
CREATE TABLE `g_counter` (
    `sql_id` bigint(20) NOT NULL AUTO_INCREMENT,
    `id` varchar(255) DEFAULT NULL,
    `gkey` varchar(255) DEFAULT NULL,
    `gtype` varchar(255) DEFAULT NULL,
    `amount` bigint(20) DEFAULT NULL,
    PRIMARY KEY (`sql_id`),
    UNIQUE KEY `id` (`id`),
    KEY `idx_g_counter_id` (`id`),
    KEY `idx_g_counter_gkey` (`gkey`),
    KEY `idx_g_counter_gtype` (`gtype`),
    KEY `idx_g_counter_amount` (`amount`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for g_counter_daily_agg
-- ----------------------------
DROP TABLE IF EXISTS `g_counter_daily_agg`;
CREATE TABLE `g_counter_daily_agg` (
    `sql_id` bigint(20) NOT NULL AUTO_INCREMENT,
    `id` varchar(255) DEFAULT NULL,
    `gkey` varchar(255) DEFAULT NULL,
    `gtype` varchar(255) DEFAULT NULL,
    `date` date DEFAULT NULL,
    `amount` bigint(20) DEFAULT NULL,
    PRIMARY KEY (`sql_id`),
    UNIQUE KEY `id` (`id`),
    KEY `idx_g_counter_daily_agg_date` (`date`),
    KEY `idx_g_counter_daily_agg_id` (`id`),
    KEY `idx_g_counter_daily_agg_gkey` (`gkey`),
    KEY `idx_g_counter_daily_agg_gtype` (`gtype`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for g_counter_detail
-- ----------------------------
DROP TABLE IF EXISTS `g_counter_detail`;
CREATE TABLE `g_counter_detail` (
    `sql_id` bigint(20) NOT NULL AUTO_INCREMENT,
    `id` varchar(255) DEFAULT NULL,
    `gkey` varchar(255) DEFAULT NULL,
    `gtype` varchar(255) DEFAULT NULL,
    `datetime` datetime(6) DEFAULT NULL,
    `amount` bigint(20) DEFAULT NULL,
    `msg` longtext,
    PRIMARY KEY (`sql_id`),
    UNIQUE KEY `id` (`id`),
    KEY `idx_g_counter_detail_gtype` (`gtype`),
    KEY `idx_g_counter_detail_datetime` (`datetime`),
    KEY `idx_g_counter_detail_id` (`id`),
    KEY `idx_g_counter_detail_gkey` (`gkey`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;

SET FOREIGN_KEY_CHECKS = 1;
