-- 初始化数据库
CREATE DATABASE IF NOT EXISTS web_data DEFAULT CHARSET utf8 COLLATE utf8_general_ci;

-- 豌豆荚app信息表
CREATE TABLE IF NOT EXISTS `web_data`.`tb_wdj_apk_info` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `app_id` VARCHAR(120) NOT NULL COMMENT '豌豆荚appid, 例如279979',
  `app_vid` VARCHAR(120) NOT NULL DEFAULT '-1' COMMENT '豌豆荚appvid, 例如400732794',
  `app_name` VARCHAR(120) NOT NULL DEFAULT '-1' COMMENT 'App名称, 例如支付宝',
  `app_pname` VARCHAR(120) NOT NULL DEFAULT '-1' COMMENT 'app包名,  唯一索引, 例如com.eg.android.AlipayGphone',
  `app_vname` VARCHAR(120) NOT NULL DEFAULT '-1' COMMENT '豌豆荚appvname, 例如10.1.55.6000',
  `app_vcode` VARCHAR(120) NOT NULL DEFAULT '-1' COMMENT '豌豆荚appvcode, 例如137',
  `app_category_id` VARCHAR(120) NOT NULL DEFAULT '-1' COMMENT '豌豆荚app分组, 例如5023',
  `app_category_name` VARCHAR(120) NOT NULL DEFAULT '-1' COMMENT '豌豆荚app分组名称, 例如金融理财',
  `app_rtype` VARCHAR(120) NOT NULL DEFAULT '-1' COMMENT '豌豆荚apprtype, 例如0',
  `app_install_count` VARCHAR(120) NOT NULL DEFAULT '-1' COMMENT '豌豆荚app安装情况, 例如525.9万人安装',
  `app_tags` VARCHAR(120) NOT NULL DEFAULT '-1' COMMENT '豌豆荚app标签分组, 以英文逗号分隔',
  `app_tag_link` VARCHAR(300) NOT NULL DEFAULT '-1' COMMENT '豌豆荚app标签链接',
  `app_icon_link` VARCHAR(300) NOT NULL DEFAULT '-1' COMMENT '豌豆荚app的icon链接',
  `status` INT NOT NULL DEFAULT 1 COMMENT '0: 抓取完成, 1: 待抓取tag',
  `create_time` DATETIME NOT NULL COMMENT '创建时间',
  `update_time` DATETIME NOT NULL COMMENT '更新时间',
  `version` INT NOT NULL DEFAULT 0 COMMENT '版本号',
  `remark` VARCHAR(500) NOT NULL DEFAULT '-1' COMMENT '备注',
  PRIMARY KEY (`id`),
  INDEX `inx_twai_appid` (`app_id` ASC),
  INDEX `inx_twai_appvid` (`app_vid` ASC),
  INDEX `inx_twai_name` (`app_name` ASC),
  INDEX `inx_twai_pname` (`app_pname` ASC),
  INDEX `inx_twai_cateid` (`app_category_id` ASC),
  INDEX `inx_twai_ctime` (`create_time` ASC),
  INDEX `inx_twai_utime` (`update_time` ASC),
  UNIQUE INDEX `app_pname_UNIQUE` (`app_pname` ASC))
  ENGINE = InnoDB
