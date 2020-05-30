drop table if exists catalog_auth;
	 CREATE TABLE `catalog_auth` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `res_type` varchar(50) NOT NULL DEFAULT '' COMMENT '资源类型',
  `res_id` varchar(50) NOT NULL DEFAULT '' COMMENT '资源id',
  `auth` tinyint(16) unsigned NOT NULL DEFAULT '0' COMMENT '资源可访问权限:2访问权,4使用权,7归属权,6访问权-使用权',
  `acc_addr` varchar(50) NOT NULL DEFAULT '' COMMENT '可访问资源机构地址',
  `own_addr` varchar(50) NOT NULL DEFAULT '' COMMENT '资源归属机构地址',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk` (`res_type`,`res_id`,`acc_addr`,`own_addr`) USING BTREE,
  KEY `idx_type_accAddr` (`res_type`,`acc_addr`) USING BTREE,
  KEY `idx_type_ownAddr` (`res_type`,`own_addr`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8