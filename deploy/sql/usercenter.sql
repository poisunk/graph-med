DROP TABLE IF EXISTS `user`;

CREATE TABLE `user` (
    `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `user_id` varchar(64) NOT NULL COMMENT '用户ID',
    `nickname` varchar(64) NOT NULL COMMENT '用户名',
    `password` varchar(64) NOT NULL COMMENT '密码',
    `mobile` varchar(64) DEFAULT NULL COMMENT '手机号',
    `email` varchar(64) DEFAULT NULL COMMENT '邮箱',
    `sex` int(11) DEFAULT NULL COMMENT '性别 0:未知 1:男 2:女',
    `avatar` varchar(255) DEFAULT NULL COMMENT '头像',
    `info` varchar(255) DEFAULT NULL COMMENT '用户信息',
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_user_id` (`user_id`),
    INDEX `idx_nickname` (`nickname`),
    INDEX `idx_email` (`email`),
    INDEX `idx_mobile` (`mobile`)
) ENGINE=InnoDB AUTO_INCREMENT=100 DEFAULT CHARSET=utf8 COMMENT='用户表';
