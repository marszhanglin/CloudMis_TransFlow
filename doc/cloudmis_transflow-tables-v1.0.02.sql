/*
Navicat MySQL Data Transfer

Source Server         : 本机localhost
Source Server Version : 50528
Source Host           : localhost:3306
Source Database       : cloudmis_transflow

Target Server Type    : MYSQL
Target Server Version : 50528
File Encoding         : 65001

Date: 2018-09-21 15:36:04
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for trans_flow
-- ----------------------------
DROP TABLE IF EXISTS `trans_flow`;
CREATE TABLE `trans_flow` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `flow_no` varchar(128) DEFAULT NULL COMMENT '上送流水号',
  `device_sn` varchar(128) DEFAULT NULL COMMENT '设备sn',
  `upload_time` varchar(128) DEFAULT NULL COMMENT '流水上送时间 YYYYMMDDHHmmss',
  `trans_time` varchar(128) DEFAULT NULL COMMENT '流水上送时间 YYYYMMDDHHmmss',
  `trans_type` varchar(128) DEFAULT NULL COMMENT '交易类型',
  `channel_id` varchar(128) DEFAULT NULL COMMENT '支付渠道',
  `merchant_id` varchar(128) DEFAULT NULL COMMENT '商户号',
  `terminal_id` varchar(128) DEFAULT NULL COMMENT '终端号',
  `merchant_name` varchar(128) DEFAULT NULL COMMENT '商户名称',
  `amount` bigint(20) DEFAULT NULL COMMENT '交易金额(分)',
  `trans_amount` bigint(20) DEFAULT NULL COMMENT '实际交易金额(分)',
  `currency_code` varchar(128) DEFAULT '156' COMMENT '货币类型',
  `out_order_no` varchar(128) DEFAULT NULL COMMENT '外部订单号',
  `voucher_no` varchar(128) DEFAULT NULL COMMENT '终端流水号',
  `reference_no` varchar(128) DEFAULT NULL COMMENT '系统流水号',
  `auth_code` varchar(128) DEFAULT NULL COMMENT '授权码',
  `ori_out_order_no` varchar(128) DEFAULT NULL COMMENT '原外部流水号',
  `ori_voucher_no` varchar(128) DEFAULT NULL COMMENT '原终端流水号',
  `ori_reference_no` varchar(128) DEFAULT NULL COMMENT '原系统流水号',
  `ori_auth_code` varchar(128) DEFAULT NULL COMMENT '原授权码',
  `card_no` varchar(128) DEFAULT NULL COMMENT '卡号',
  `operator_no` varchar(128) DEFAULT NULL COMMENT '操作员编号',
  `combination_no` varchar(128) DEFAULT NULL COMMENT '组合支付编号',
  `cardType` varchar(128) DEFAULT NULL COMMENT '银行卡类型',
  `remark` varchar(512) DEFAULT NULL COMMENT '备注',
  `extendParams` varchar(1024) DEFAULT NULL COMMENT '扩展字段',
  `shop_id` bigint(20) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8 COMMENT='流水表';

-- ----------------------------
-- Records of trans_flow
-- ----------------------------
INSERT INTO `trans_flow` VALUES ('9', '1221', '12345678901111', '2016-01-02 23:04:03', '2016-01-02 23:04:03', '1', 'wxpay', 'merchantId', 'terminalId', 'merchantName', '0', '5', 'currencyCode', 'outOrderNo', 'voucherNo', 'referenceNo', 'authCode', 'oriOutOrderNo', 'oriVoucherNo', 'oriReferenceNo', 'oriAuthCode', 'cardNo', 'operatorNo', 'combinationNo', 'cardType', 'remark', 'extendParams1', '1');
INSERT INTO `trans_flow` VALUES ('10', '123121231213123234123123121231213', '12345678901111', '2016-01-02 23:04:03', '2016-01-02 23:04:03', '2', 'wxpay', 'merchantId', 'terminalId', 'merchantName', '0', '5', 'currencyCode', 'outOrderNo', 'voucherNo', 'referenceNo', 'authCode', 'oriOutOrderNo', 'oriVoucherNo', 'oriReferenceNo', 'oriAuthCode', 'cardNo', 'operatorNo', 'combinationNo', 'cardType', 'remark', 'extendParams1', '1');

-- ----------------------------
-- Table structure for t_operator
-- ----------------------------
DROP TABLE IF EXISTS `t_operator`;
CREATE TABLE `t_operator` (
  `operator_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `operator_name` varchar(128) DEFAULT NULL COMMENT '客户名称',
  `address` varchar(128) DEFAULT NULL COMMENT '地址',
  `linkman` varchar(128) DEFAULT NULL COMMENT '联系人',
  `mobile` varchar(128) DEFAULT NULL COMMENT '手机号',
  `operator_status` bigint(20) DEFAULT NULL COMMENT '状态：\r\n            0 - 锁定\r\n            1 - 正常\r\n            2 - 待审核',
  `cre_time` datetime DEFAULT NULL,
  `upd_time` datetime DEFAULT NULL,
  `remark` varchar(512) DEFAULT NULL,
  PRIMARY KEY (`operator_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COMMENT='运营商信息表';

-- ----------------------------
-- Records of t_operator
-- ----------------------------
INSERT INTO `t_operator` VALUES ('1', '新大陆测试运营商', '福州马尾儒江', '新大陆联系人', '18912345678', '1', '2016-12-11 11:12:34', '2016-12-11 11:12:57', '审核通过');

-- ----------------------------
-- Table structure for t_role
-- ----------------------------
DROP TABLE IF EXISTS `t_role`;
CREATE TABLE `t_role` (
  `role_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `role_code` varchar(128) NOT NULL COMMENT '角色标识（英文）',
  `role_name` varchar(128) NOT NULL COMMENT '角色名称',
  `role_remark` varchar(128) DEFAULT NULL,
  `cre_time` datetime NOT NULL,
  `upd_time` datetime DEFAULT NULL,
  PRIMARY KEY (`role_id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of t_role
-- ----------------------------
INSERT INTO `t_role` VALUES ('1', 'sys_admin', '系统管理员', null, '2015-11-21 00:00:00', null);
INSERT INTO `t_role` VALUES ('2', 'operator_admin', '运营商管理员', null, '2015-11-21 00:00:00', null);
INSERT INTO `t_role` VALUES ('3', 'shop_admin', '店铺管理员', null, '2015-11-11 00:00:00', null);
INSERT INTO `t_role` VALUES ('4', 'base_admin', '操作员', null, '2015-11-11 00:00:00', null);

-- ----------------------------
-- Table structure for t_shop
-- ----------------------------
DROP TABLE IF EXISTS `t_shop`;
CREATE TABLE `t_shop` (
  `shop_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `shop_name` varchar(120) DEFAULT NULL COMMENT '店铺名称',
  `address` varchar(120) DEFAULT NULL COMMENT '地址',
  `operator_id` bigint(20) NOT NULL COMMENT '所属运营商信息表主键',
  `cre_time` datetime DEFAULT NULL,
  `upd_time` datetime DEFAULT NULL,
  `remark` varchar(512) DEFAULT NULL,
  PRIMARY KEY (`shop_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COMMENT='店铺信息表';

-- ----------------------------
-- Records of t_shop
-- ----------------------------
INSERT INTO `t_shop` VALUES ('1', '新大陆测试门店', '福州马尾儒江', '1', '2016-12-11 11:12:34', '2016-12-11 11:12:57', '');

-- ----------------------------
-- Table structure for t_user
-- ----------------------------
DROP TABLE IF EXISTS `t_user`;
CREATE TABLE `t_user` (
  `user_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `user_name` varchar(128) DEFAULT NULL COMMENT '登陆账号',
  `password` varchar(128) DEFAULT NULL COMMENT '登陆密码',
  `user_status` bigint(20) DEFAULT NULL COMMENT '用户状态:\r\n            1 - 正常\r\n            0 - 锁定\r\n            ',
  `user_type` bigint(20) NOT NULL COMMENT '用户类型:     管理员|运营商|店长|操作员\r\n             1 - 管理员\r\n             2 - 运营商\r\n             3 - 店长\r\n             4 - 操作员\r\n            ',
  `operator_id` bigint(20) DEFAULT NULL COMMENT '所属运营商信息表主键',
  `shop_id` bigint(20) DEFAULT NULL COMMENT '所属店铺信息表主键',
  `cre_time` datetime DEFAULT NULL,
  `upd_time` datetime DEFAULT NULL,
  `remark` varchar(512) DEFAULT NULL,
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=153499085791887 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of t_user
-- ----------------------------
INSERT INTO `t_user` VALUES ('1', 'admin', '88888888', '1', '1', null, null, '2015-11-03 17:15:41', '2017-01-23 11:25:53', null);
INSERT INTO `t_user` VALUES ('2', 'newlandyys', '88888888', '1', '2', '1', null, '2016-12-11 11:12:12', '2017-01-15 15:17:22', null);
INSERT INTO `t_user` VALUES ('3', 'newlanddz', '88888888', '1', '3', '1', '1', '2016-12-11 11:12:12', '2017-01-15 15:17:22', null);
INSERT INTO `t_user` VALUES ('4', 'newlandczy', '88888888', '1', '4', '1', '1', '2016-12-11 11:12:12', '2017-01-15 15:17:22', null);
INSERT INTO `t_user` VALUES ('5', 'newlandczy02', '77777777', '1', '4', '1', '1', '2018-08-07 16:33:22', '2018-08-07 16:33:22', '添加者：newlanddz');
INSERT INTO `t_user` VALUES ('6', 'newlandczy01', '77777777', '1', '4', '1', '1', '2018-08-07 16:43:43', '2018-08-07 16:43:43', '添加者：newlanddz');
INSERT INTO `t_user` VALUES ('7', 'newlandczy03', '88888888', '1', '4', '1', '1', '2018-08-10 09:27:58', '2018-08-10 09:27:58', '添加者：newlanddz');
INSERT INTO `t_user` VALUES ('8', 'newlandczy04', '88888888', '1', '4', '1', '1', '2018-08-10 09:30:34', '2018-08-10 09:30:34', '添加者：newlanddz');
INSERT INTO `t_user` VALUES ('9', 'newlandczy05', '88888888', '1', '4', '1', '1', '2018-08-10 09:31:09', '2018-08-10 09:31:09', '添加者：newlanddz');
INSERT INTO `t_user` VALUES ('10', '18000000001', '12345678', '1', '4', '1', '1', '2018-08-20 16:21:51', '2018-08-20 16:21:51', '添加者：newlanddz');
INSERT INTO `t_user` VALUES ('11', '18000000002', '88888888', '1', '4', '1', '1', '2018-08-20 16:42:45', '2018-08-20 16:42:45', '添加者：newlanddz');
INSERT INTO `t_user` VALUES ('153499085791886', 'newlandczy09', '88888888', '1', '4', '1', '1', '2018-08-23 10:20:57', '2018-08-23 10:20:57', '添加者：newlanddz');
