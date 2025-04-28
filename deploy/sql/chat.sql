DROP TABLE IF EXISTS `chat_message`;

CREATE TABLE `chat_message` (
    `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `dialogue_id` int(11) NOT NULL COMMENT '消息ID',
    `message_type` varchar(64) NOT NULL COMMENT '消息类型',
    `content` longtext NOT NULL COMMENT '内容',
    `function_name` varchar(64) DEFAULT NULL COMMENT '函数名称',
    `function_args` longtext DEFAULT NULL COMMENT '函数参数',
    `function_result` longtext DEFAULT NULL COMMENT '函数结果',
    `other` longtext DEFAULT NULL COMMENT '其他信息',
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    INDEX `idx_dialogue_id` (`dialogue_id`)
) ENGINE=InnoDB AUTO_INCREMENT=100 DEFAULT CHARSET=utf8 COMMENT='消息表';

DROP TABLE IF EXISTS `chat_dialogue`;

CREATE TABLE `chat_dialogue` (
    `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `dialogue_id` int(11) NOT NULL COMMENT '消息ID',
    `parent_dialogue_id` int(11) NOT NULL DEFAULT 0 COMMENT '父ID',
    `chat_session_id` varchar(64) NOT NULL COMMENT '会话ID',
    `user_id` varchar(64) NOT NULL COMMENT '用户ID',
    `model_name` varchar(100) NOT NULL DEFAULT 'doubao-1-5-lite-32k-250115' COMMENT '模型名称',
    `role` ENUM('user', 'assistant', 'system') NOT NULL COMMENT '角色',
    `content` longtext NOT NULL COMMENT '内容',
    `operator` varchar(64) DEFAULT NULL COMMENT '操作人',
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_dialogue_id` (`dialogue_id`),
    INDEX `idx_chat_session_id` (`chat_session_id`),
    INDEX `idx_user_id` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=100 DEFAULT CHARSET=utf8 COMMENT='对话表';

DROP TABLE IF EXISTS `chat_session`;

CREATE TABLE `chat_session` (
    `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `session_id` varchar(64) NOT NULL COMMENT '会话ID',
    `type_id` varchar(64) NOT NULL DEFAULT 'default' COMMENT '类型ID',
    `user_id` varchar(64) NOT NULL COMMENT '用户ID',
    `title` varchar(100) NOT NULL DEFAULT '' COMMENT '标题',
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_session_id` (`session_id`),
    INDEX `idx_user_id` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=100 DEFAULT CHARSET=utf8 COMMENT='对话会话表';

DROP TABLE IF EXISTS `chat_dialogue_feedback`;

CREATE TABLE `chat_dialogue_feedback` (
    `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `dialogue_id` int(11) NOT NULL COMMENT '消息ID',
    `user_id` varchar(64) NOT NULL COMMENT '用户ID',
    `chat_session_id` varchar(64) NOT NULL COMMENT '会话ID',
    `feedback` varchar(255) DEFAULT NULL COMMENT '反馈',
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    INDEX `idx_chat_session_id` (`chat_session_id`),
    INDEX `idx_dialogue_id` (`dialogue_id`),
    INDEX `idx_user_id` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=100 DEFAULT CHARSET=utf8 COMMENT='对话消息反馈表';

DROP TABLE IF EXISTS `chat_type`;

CREATE TABLE `chat_type` (
    `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `type_id` varchar(64) NOT NULL COMMENT '类型ID',
    `type_name` varchar(64) NOT NULL COMMENT '类型名称',
    `mcp_ids` varchar(255) NOT NULL DEFAULT '' COMMENT 'MCP IDs (逗号,分隔)',
    `model_name` varchar(64) NOT NULL DEFAULT 'doubao-1-5-lite-32k-250115' COMMENT '模型名称',
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_type_id` (`type_id`)
) ENGINE=InnoDB AUTO_INCREMENT=100 DEFAULT CHARSET=utf8 COMMENT='对话类型表';

DROP TABLE IF EXISTS `mcp_service`;

CREATE TABLE `mcp_service` (
    `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `mcp_type` varchar(64) NOT NULL COMMENT 'MCP类型',
    `name` varchar(64) NOT NULL COMMENT 'MCP名称',
    `args` text DEFAULT NULL COMMENT '运行参数',
    `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=100 DEFAULT CHARSET=utf8 COMMENT='MCP服务表';

INSERT INTO chat_type (type_id, type_name, mcp_ids, model_name, created_at, updated_at)
VALUES ('default', '知识图谱对话', '', 'doubao-1-5-lite-32k-250115', DEFAULT, DEFAULT);