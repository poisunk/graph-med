DROP TABLE IF EXISTS `chat_session`;

CREATE TABLE `chat_session` (
    `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `type_id` varchar(64) NOT NULL DEFAULT 'default' COMMENT '类型ID',
    `session_id` varchar(64) NOT NULL COMMENT '会话ID',
    `user_id` varchar(64) NOT NULL COMMENT '用户ID',
    `title` varchar(100) DEFAULT NULL COMMENT '标题',
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_session_id` (`session_id`),
    INDEX `idx_user_id` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=100 DEFAULT CHARSET=utf8 COMMENT='对话会话表';

DROP TABLE IF EXISTS `chat_message`;

CREATE TABLE `chat_message` (
    `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `message_id` int(11) NOT NULL COMMENT '消息ID',
    `parent_message_id` int(11) DEFAULT NULL COMMENT '父ID',
    `user_id` varchar(64) NOT NULL COMMENT '用户ID',
    `chat_session_id` varchar(64) NOT NULL COMMENT '会话ID',
    `model_name` varchar(100) DEFAULT NULL COMMENT '模型名称',
    `role` ENUM('user', 'assistant', 'system') NOT NULL COMMENT '角色',
    `content` longtext DEFAULT NULL COMMENT '内容',
    `operator` varchar(64) DEFAULT NULL COMMENT '操作人',
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_chat_session_id` (`chat_session_id`),
    INDEX `idx_message_id` (`message_id`),
    INDEX `idx_user_id` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=100 DEFAULT CHARSET=utf8 COMMENT='对话消息表';

DROP TABLE IF EXISTS `chat_message_feedback`;

CREATE TABLE `chat_message_feedback` (
    `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `message_id` int(11) NOT NULL COMMENT '消息ID',
    `user_id` varchar(64) NOT NULL COMMENT '用户ID',
    `chat_session_id` varchar(64) NOT NULL COMMENT '会话ID',
    `feedback` varchar(255) DEFAULT NULL COMMENT '反馈',
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_chat_session_id` (`chat_session_id`),
    INDEX `idx_message_id` (`message_id`),
    INDEX `idx_user_id` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=100 DEFAULT CHARSET=utf8 COMMENT='对话消息反馈表';

DROP TABLE IF EXISTS `chat_user`;

CREATE TABLE `chat_user` (
    `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `user_id` varchar(64) NOT NULL COMMENT '用户ID',
    `username` varchar(64) NOT NULL COMMENT '用户名',
    `password` varchar(64) NOT NULL COMMENT '密码',
    `email` varchar(64) DEFAULT NULL COMMENT '邮箱',
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_user_id` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=100 DEFAULT CHARSET=utf8 COMMENT='用户表';

DROP TABLE IF EXISTS `chat_type`;

CREATE TABLE `chat_type` (
    `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `type_id` varchar(64) NOT NULL COMMENT '类型ID',
    `type_name` varchar(64) NOT NULL COMMENT '类型名称',
    `mcp_ids` varchar(255) DEFAULT NULL COMMENT 'MCP IDs (逗号,分隔)',
    `model_name` varchar(64) DEFAULT NULL COMMENT '模型名称',
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_type_id` (`type_id`)
) ENGINE=InnoDB AUTO_INCREMENT=100 DEFAULT CHARSET=utf8 COMMENT='对话类型表';

DROP TABLE IF EXISTS `node_info`;

CREATE TABLE `node_info` (
     `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
     `label` varchar(100) NOT NULL COMMENT '标签',
     `name` varchar(100) NOT NULL COMMENT '名称',
     `attr_name` varchar(100) NOT NULL COMMENT '属性名',
     `attr_value` longtext NOT NULL COMMENT '属性值',
     `group` varchar(100) NOT NULL COMMENT '分组',
     `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
     `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
     `deleted_at` datetime DEFAULT NULL,
     PRIMARY KEY (`id`),
     INDEX `idx_label` (`label`),
     INDEX `idx_name` (`name`),
     INDEX `idx_attr_name` (`attr_name`)
) ENGINE=InnoDB AUTO_INCREMENT=100 DEFAULT CHARSET=utf8 COMMENT='节点信息表';

DROP TABLE IF EXISTS `mcp_service`;

CREATE TABLE `mcp_service` (
    `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `name` varchar(64) NOT NULL COMMENT 'MCP名称',
    `mcp_id` varchar(64) NOT NULL COMMENT 'MCP ID',
    `type` varchar(64) NOT NULL COMMENT '类型',
    `args` longtext DEFAULT NULL COMMENT '参数',
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=100 DEFAULT CHARSET=utf8 COMMENT='MCP服务表';


# casbin 生成的表
DROP TABLE IF EXISTS `casbin_rule`;

# 初始化一些数据
INSERT INTO chat_type (type_id, type_name, mcp_ids, model_name, created_at, updated_at, deleted_at)
VALUES ('default', '知识图谱对话', 'default', 'doubao-1-5-lite-32k-250115', DEFAULT, DEFAULT, null);
INSERT INTO mcp_service (name, mcp_id, type, args, created_at, updated_at, deleted_at)
VALUES ('knowledge_graph_assistant', 'default', 'sse', '{"baseURL":"http://localhost:8081"}', DEFAULT, DEFAULT, null);