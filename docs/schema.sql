-- develop_tools MySQL schema
-- 与 internal/model 业务表结构对齐（utf8mb4）
-- 新建库直接执行本文件。

CREATE TABLE IF NOT EXISTS `my_user` (
  `id`         INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
  `name`       VARCHAR(128) NOT NULL DEFAULT '' COMMENT '用户名（唯一）',
  `created_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_my_user_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

CREATE TABLE IF NOT EXISTS `my_user_key` (
  `id`          INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
  `user_id`     INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户 ID，关联 my_user.id',
  `browser_key` VARCHAR(128) NOT NULL DEFAULT '' COMMENT '浏览器指纹 / 本地密钥',
  `user_agent`  VARCHAR(1024) NOT NULL DEFAULT '' COMMENT '最近一次 User-Agent',
  `created_at`  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_my_user_key_user_browser` (`user_id`, `browser_key`),
  KEY `idx_my_user_key_browser_key` (`browser_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户浏览器密钥表（登录态关联）';

CREATE TABLE IF NOT EXISTS `my_share` (
  `id`         INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
  `op_type`    INT          NOT NULL DEFAULT 0 COMMENT '操作类型：1收藏 2分享 3下载',
  `user_id`    INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '所属用户 ID，关联 my_user.id',
  `uuid`       VARCHAR(64)  NOT NULL DEFAULT '' COMMENT '前端本地用户 UUID',
  `path`       VARCHAR(255) NOT NULL DEFAULT '' COMMENT '页面路径，如 /json',
  `name`       VARCHAR(255) NOT NULL DEFAULT '' COMMENT '分享/收藏名称',
  `token`      VARCHAR(128) NOT NULL DEFAULT '' COMMENT '访问令牌',
  `data`       LONGTEXT     NOT NULL COMMENT '业务快照 JSON',
  `status`     INT          NOT NULL DEFAULT 0 COMMENT '状态：0正常 1已删除（软删）',
  `created_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_my_share_path_token_status` (`path`, `token`, `status`),
  KEY `idx_my_share_user_op_status` (`user_id`, `op_type`, `status`),
  KEY `idx_my_share_updated_at` (`updated_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='收藏/分享/下载记录表';

CREATE TABLE IF NOT EXISTS `my_dsp` (
  `id`                   INT UNSIGNED  NOT NULL AUTO_INCREMENT COMMENT '主键',
  `name`                 VARCHAR(255)  NOT NULL DEFAULT '' COMMENT 'DSP 名称',
  `unique_key`           VARCHAR(64)   NOT NULL DEFAULT '' COMMENT '对外唯一标识（竞价 URL 路径参数）',
  `is_cn`                TINYINT       NOT NULL DEFAULT 0 COMMENT '市场：0海外 1国内',
  `request_id`           VARCHAR(128)  NOT NULL DEFAULT '' COMMENT '响应 request id，空或 {REQUEST_ID} 时用请求 id',
  `price`                DOUBLE        NOT NULL DEFAULT 0 COMMENT '出价价格，0 时使用默认价',
  `adm`                  LONGTEXT      NOT NULL COMMENT '广告物料 / 竞价响应体',
  `crid`                 VARCHAR(128)  NOT NULL DEFAULT '' COMMENT '创意 ID',
  `bundle`               VARCHAR(255)  NOT NULL DEFAULT '' COMMENT '应用包名 / Bundle',
  `deeplink`             VARCHAR(1024) NOT NULL DEFAULT '' COMMENT 'DeepLink',
  `deeplinkfallbackurl`  VARCHAR(1024) NOT NULL DEFAULT '' COMMENT 'DeepLink 失败回落 URL',
  `fallback`             VARCHAR(1024) NOT NULL DEFAULT '' COMMENT '通用回落 URL',
  `created_at`           DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`           DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_my_dsp_unique_key` (`unique_key`),
  KEY `idx_my_dsp_is_cn_updated_at` (`is_cn`, `updated_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='ADX 模拟 DSP 配置表';

CREATE TABLE IF NOT EXISTS `my_dsp_notice` (
  `id`          INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
  `dsp_id`      INT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'DSP ID，关联 my_dsp.id',
  `notice_type` INT          NOT NULL DEFAULT 0 COMMENT '回调类型：1burl 2nurl 3lurl 4tpnurl 5tplurl 6tpburl 7tpimpurl 8tpclkurl',
  `ip`          VARCHAR(64)  NOT NULL DEFAULT '' COMMENT '回调来源 IP',
  `ua`          VARCHAR(1024) NOT NULL DEFAULT '' COMMENT '回调 User-Agent',
  `created_at`  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间（回调发生时间）',
  `updated_at`  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_my_dsp_notice_dsp_id_id` (`dsp_id`, `id`),
  KEY `idx_my_dsp_notice_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='ADX DSP 竞价回调/曝光点击通知日志';
