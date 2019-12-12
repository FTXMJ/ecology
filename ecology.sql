-- --------------------------------------------------------
-- 主机:                           127.0.0.1
-- 服务器版本:                        5.6.46 - MySQL Community Server (GPL)
-- 服务器OS:                        Linux
-- HeidiSQL 版本:                  10.2.0.5599
-- --------------------------------------------------------

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!50503 SET NAMES utf8mb4 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;


-- Dumping database structure for ecology
CREATE DATABASE IF NOT EXISTS `ecology` /*!40100 DEFAULT CHARACTER SET utf8 COLLATE utf8_swedish_ci */;
USE `ecology`;

-- Dumping structure for table ecology.account
CREATE TABLE IF NOT EXISTS `account` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` varchar(50) COLLATE utf8_swedish_ci NOT NULL DEFAULT '' COMMENT '用户 id',
  `balance` double DEFAULT '0' COMMENT '交易结余',
  `currency` varchar(50) COLLATE utf8_swedish_ci NOT NULL DEFAULT 'USDD' COMMENT '货币',
  `bocked_balance` double NOT NULL DEFAULT '0' COMMENT '铸币结余',
  `level` varchar(50) COLLATE utf8_swedish_ci DEFAULT NULL COMMENT '等级',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=57 DEFAULT CHARSET=utf8 COLLATE=utf8_swedish_ci COMMENT='生态（钱包）';

-- Data exporting was unselected.

-- Dumping structure for table ecology.account_detail
CREATE TABLE IF NOT EXISTS `account_detail` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` varchar(50) COLLATE utf8_swedish_ci NOT NULL,
  `current_revenue` double NOT NULL DEFAULT '0' COMMENT '本期收入',
  `current_outlay` double NOT NULL DEFAULT '0' COMMENT '本期支出',
  `opening_balance` double NOT NULL DEFAULT '0' COMMENT '上期余额',
  `current_balance` double NOT NULL DEFAULT '0' COMMENT '本期余额',
  `create_date` datetime NOT NULL COMMENT '发生交易日期',
  `comment` varchar(200) COLLATE utf8_swedish_ci NOT NULL COMMENT '评论',
  `tx_id` varchar(200) COLLATE utf8_swedish_ci NOT NULL COMMENT '交易唯一id',
  `account` int(11) NOT NULL COMMENT '生态钱包id',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=25 DEFAULT CHARSET=utf8 COLLATE=utf8_swedish_ci COMMENT='交易记录表';

-- Data exporting was unselected.

-- Dumping structure for table ecology.blocked_detail
CREATE TABLE IF NOT EXISTS `blocked_detail` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` varchar(50) COLLATE utf8_swedish_ci NOT NULL,
  `current_revenue` double NOT NULL DEFAULT '0' COMMENT '本期收入',
  `current_outlay` double NOT NULL DEFAULT '0' COMMENT '本期支出',
  `opening_balance` double NOT NULL DEFAULT '0' COMMENT '上期余额',
  `current_balance` double NOT NULL DEFAULT '0' COMMENT '本期余额',
  `create_date` datetime NOT NULL COMMENT '发生交易时间',
  `comment` varchar(200) COLLATE utf8_swedish_ci NOT NULL COMMENT '评论',
  `tx_id` varchar(200) COLLATE utf8_swedish_ci NOT NULL COMMENT '交易唯一id',
  `account` int(11) NOT NULL COMMENT '生态钱包id',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=utf8 COLLATE=utf8_swedish_ci COMMENT='铸币表';

-- Data exporting was unselected.

-- Dumping structure for table ecology.calculation_power
CREATE TABLE IF NOT EXISTS `calculation_power` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `datetime` varchar(50) COLLATE utf8_swedish_ci NOT NULL DEFAULT '' COMMENT '日期',
  `principal_calculation` double NOT NULL DEFAULT '0' COMMENT '自由算力',
  `direct_calculation_force` double NOT NULL DEFAULT '0' COMMENT '加速算力',
  `dynamic_calculation` double NOT NULL DEFAULT '0' COMMENT '动态算力',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_swedish_ci COMMENT='算力表';

-- Data exporting was unselected.

-- Dumping structure for table ecology.daily_dividend_tasks
CREATE TABLE IF NOT EXISTS `daily_dividend_tasks` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `time` varchar(50) COLLATE utf8_swedish_ci NOT NULL COMMENT '更新时间的记录，只记录 2006-01-02 格式',
  `state` varchar(50) COLLATE utf8_swedish_ci NOT NULL DEFAULT 'false' COMMENT 'false=未完成 , true=完成',
  `completion_time` datetime DEFAULT NULL COMMENT '完成任务的具体时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_swedish_ci COMMENT='每日释放和分红的任务表';

-- Data exporting was unselected.

-- Dumping structure for table ecology.force_table
CREATE TABLE IF NOT EXISTS `force_table` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `level` varchar(50) COLLATE utf8_swedish_ci NOT NULL COMMENT '等级',
  `low_hold` int(11) NOT NULL COMMENT '等级区间(最低)　－－　单位　USDD',
  `high_hold` int(11) NOT NULL COMMENT '等级区间(最高)　－－　单位　USDD',
  `return_multiple` double NOT NULL COMMENT '杠杆',
  `hold_return_rate` double NOT NULL COMMENT '本金自由算力',
  `recommend_return_rate` double NOT NULL COMMENT '直推加速算力',
  `team_return_rate` double NOT NULL COMMENT '团队加速',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8 COLLATE=utf8_swedish_ci COMMENT='后台管理的统一算力定义表';

-- Data exporting was unselected.

-- Dumping structure for table ecology.formula
CREATE TABLE IF NOT EXISTS `formula` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `ecology_id` int(11) NOT NULL DEFAULT '0' COMMENT '对应的生态项目id',
  `level` varchar(50) COLLATE utf8_swedish_ci NOT NULL,
  `low_hold` int(11) NOT NULL COMMENT '低位',
  `high_hold` int(11) NOT NULL COMMENT '高位',
  `return_multiple` double NOT NULL DEFAULT '1' COMMENT '杠杆',
  `hold_return_rate` double NOT NULL COMMENT '本金自由算力',
  `recommend_return_rate` double NOT NULL COMMENT '加速算力',
  `team_return_rate` double NOT NULL COMMENT '动态算力',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=27 DEFAULT CHARSET=utf8 COLLATE=utf8_swedish_ci COMMENT='（公式表？？？）';

-- Data exporting was unselected.

-- Dumping structure for table ecology.super_force_table
CREATE TABLE IF NOT EXISTS `super_force_table` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `level` varchar(50) COLLATE utf8_swedish_ci NOT NULL COMMENT '超级节点的等级',
  `coin_number_rule` int(11) NOT NULL COMMENT '持币要求',
  `bonus_calculation` double NOT NULL DEFAULT '0' COMMENT '享受全网当日总算力的分红',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8 COLLATE=utf8_swedish_ci COMMENT='超级节点的算力规定表';

-- Data exporting was unselected.

-- Dumping structure for table ecology.super_peer_table
CREATE TABLE IF NOT EXISTS `super_peer_table` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` varchar(50) COLLATE utf8_swedish_ci NOT NULL DEFAULT '',
  `coin_number` double NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8 COLLATE=utf8_swedish_ci COMMENT='超级节点表';

-- Data exporting was unselected.

-- Dumping structure for table ecology.tx_id_list
CREATE TABLE IF NOT EXISTS `tx_id_list` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `state` varchar(50) COLLATE utf8_swedish_ci NOT NULL DEFAULT 'false' COMMENT '任务完成状态 默认为 ___  "false"->表示没有完成这个任务，"true"->完成',
  `tx_id` varchar(50) COLLATE utf8_swedish_ci NOT NULL,
  `user_id` varchar(50) COLLATE utf8_swedish_ci NOT NULL,
  `create_time` datetime NOT NULL,
  `expenditure` double NOT NULL DEFAULT '0' COMMENT '支出金额',
  `income` double NOT NULL DEFAULT '0' COMMENT '收入金额',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=200 DEFAULT CHARSET=utf8 COLLATE=utf8_swedish_ci COMMENT='交易唯一标示id表';

-- Data exporting was unselected.

-- Dumping structure for table ecology.user
CREATE TABLE IF NOT EXISTS `user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` varchar(200) COLLATE utf8_swedish_ci NOT NULL COMMENT '对应 monggodb 的user_id',
  `father_id` varchar(50) COLLATE utf8_swedish_ci DEFAULT NULL COMMENT '父亲ｉｄ',
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_id` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8 COLLATE=utf8_swedish_ci COMMENT='用户信息表';

-- Data exporting was unselected.

/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IF(@OLD_FOREIGN_KEY_CHECKS IS NULL, 1, @OLD_FOREIGN_KEY_CHECKS) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
