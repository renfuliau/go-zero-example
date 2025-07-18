-- order/rpc/model/order.sql

CREATE TABLE `orders` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '訂單ID',
  `order_sn` varchar(255) NOT NULL DEFAULT '' COMMENT '訂單編號',
  `user_id` bigint NOT NULL COMMENT '用戶ID',
  `total_amount` decimal(10,2) NOT NULL COMMENT '訂單總金額',
  `status` tinyint NOT NULL DEFAULT '0' COMMENT '訂單狀態 0:待支付 1:已支付 2:已發貨 3:已完成 4:已取消',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_order_sn` (`order_sn`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;