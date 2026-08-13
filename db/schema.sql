-- 用户表
CREATE TABLE `user` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '用户ID',
    `username` varchar(50) NOT NULL COMMENT '用户名',
    `password` varchar(255) NOT NULL COMMENT '密码',
    `phone` varchar(20) NOT NULL COMMENT '手机号',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` datetime DEFAULT NULL COMMENT '逻辑删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_username` (`username`),
    UNIQUE KEY `uk_phone` (`phone`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';
--商品表
CREATE TABLE `product` (
   `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
   `name` VARCHAR(100) NOT NULL COMMENT '商品名称',
   `description` TEXT COMMENT '商品描述',
   `price` DECIMAL(10,2) NOT NULL COMMENT '售价',
   `original_price` DECIMAL(10,2) DEFAULT 0 COMMENT '原价',
   `stock` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '库存',
   `category_id` BIGINT UNSIGNED DEFAULT 0 COMMENT '分类ID',
   `main_image` VARCHAR(255) DEFAULT '' COMMENT '主图URL',
   `images` JSON COMMENT '商品图片列表',
   `status` TINYINT NOT NULL DEFAULT 1 COMMENT '1=上架 0=下架',
   `sales` INT UNSIGNED DEFAULT 0 COMMENT '销量',
   `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
   `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   PRIMARY KEY (`id`),
   KEY `idx_category_id` (`category_id`),
   KEY `idx_status` (`status`),
   KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品表';
--商品分类表
CREATE TABLE `product_category` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(50) NOT NULL COMMENT '分类名称',
    `parent_id` BIGINT UNSIGNED DEFAULT 0 COMMENT '父分类ID，0=一级分类',
    `sort` INT DEFAULT 0 COMMENT '排序',
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_parent_id` (`parent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品分类表';
--购物车表
CREATE TABLE `cart_item` (
     `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
     `user_id` BIGINT UNSIGNED NOT NULL,
     `product_id` BIGINT UNSIGNED NOT NULL,
     `quantity` INT UNSIGNED NOT NULL DEFAULT 1 COMMENT '数量',
     `checked` TINYINT NOT NULL DEFAULT 1 COMMENT '1=选中 0=未选中',
     `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
     `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
     PRIMARY KEY (`id`),
     UNIQUE KEY `uk_user_product` (`user_id`, `product_id`),
     KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='购物车表';
--订单表
CREATE TABLE `order` (
     `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
     `order_no` VARCHAR(64) NOT NULL COMMENT '订单编号（唯一，对外暴露）',
     `user_id` BIGINT UNSIGNED NOT NULL,
     `total_amount` DECIMAL(10,2) NOT NULL COMMENT '订单总金额',
     `pay_amount` DECIMAL(10,2) NOT NULL COMMENT '实付金额（扣优惠后）',
     `coupon_id` BIGINT UNSIGNED DEFAULT 0 COMMENT '使用的优惠券ID',
     `coupon_amount` DECIMAL(10,2) DEFAULT 0 COMMENT '优惠金额',
     `status` TINYINT NOT NULL DEFAULT 0 COMMENT '0=待支付 1=已支付 2=已发货 3=已完成 4=已取消 5=已退款',
     `address_id` BIGINT UNSIGNED DEFAULT 0 COMMENT '收货地址ID',
     `remark` VARCHAR(200) DEFAULT '' COMMENT '订单备注',
     `paid_at` DATETIME DEFAULT NULL COMMENT '支付时间',
     `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
     `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
     PRIMARY KEY (`id`),
     UNIQUE KEY `uk_order_no` (`order_no`),
     KEY `idx_user_id` (`user_id`),
     KEY `idx_status` (`status`),
     KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单主表';
--订单明细表
CREATE TABLE `order_item` (
      `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
      `order_id` BIGINT UNSIGNED NOT NULL,
      `order_no` VARCHAR(64) NOT NULL,
      `product_id` BIGINT UNSIGNED NOT NULL,
      `product_name` VARCHAR(100) NOT NULL COMMENT '快照：商品名称',
      `product_image` VARCHAR(255) DEFAULT '' COMMENT '快照：商品图片',
      `price` DECIMAL(10,2) NOT NULL COMMENT '快照：购买时单价',
      `quantity` INT UNSIGNED NOT NULL COMMENT '购买数量',
      `total_price` DECIMAL(10,2) NOT NULL COMMENT '小计',
      `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
      PRIMARY KEY (`id`),
      KEY `idx_order_id` (`order_id`),
      KEY `idx_order_no` (`order_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单明细表';
--优惠券模板表
CREATE TABLE `coupon` (
      `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
      `name` VARCHAR(100) NOT NULL COMMENT '优惠券名称',
      `type` TINYINT NOT NULL COMMENT '1=满减 2=折扣 3=无门槛',
      `value` DECIMAL(10,2) NOT NULL COMMENT '减免金额 或 折扣率（如 0.85 = 85折）',
      `min_amount` DECIMAL(10,2) DEFAULT 0 COMMENT '满减门槛（0=无门槛）',
      `total_count` INT UNSIGNED DEFAULT 0 COMMENT '发放总量（0=不限）',
      `received_count` INT UNSIGNED DEFAULT 0 COMMENT '已领取数量',
      `per_limit` INT UNSIGNED DEFAULT 1 COMMENT '每人限领数量',
      `valid_days` INT UNSIGNED DEFAULT 7 COMMENT '领取后有效天数',
      `start_time` DATETIME DEFAULT NULL COMMENT '领取开始时间',
      `end_time` DATETIME DEFAULT NULL COMMENT '领取结束时间',
      `use_start_time` DATETIME DEFAULT NULL COMMENT '使用开始时间',
      `use_end_time` DATETIME DEFAULT NULL COMMENT '使用结束时间',
      `status` TINYINT DEFAULT 1 COMMENT '1=生效 0=失效',
      `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
      PRIMARY KEY (`id`),
      KEY `idx_status` (`status`),
      KEY `idx_end_time` (`end_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='优惠券模板表';
--用户优惠券表
CREATE TABLE `user_coupon` (
   `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
   `user_id` BIGINT UNSIGNED NOT NULL,
   `coupon_id` BIGINT UNSIGNED NOT NULL,
   `order_id` BIGINT UNSIGNED DEFAULT 0 COMMENT '使用的订单ID（0=未使用）',
   `status` TINYINT NOT NULL DEFAULT 0 COMMENT '0=未使用 1=已使用 2=已过期',
   `expire_at` DATETIME NOT NULL COMMENT '过期时间',
   `used_at` DATETIME DEFAULT NULL COMMENT '使用时间',
   `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
   PRIMARY KEY (`id`),
   UNIQUE KEY `uk_user_coupon` (`user_id`, `coupon_id`),
   KEY `idx_user_id` (`user_id`),
   KEY `idx_status` (`status`),
   KEY `idx_expire_at` (`expire_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户优惠券表';
--收货地址表
CREATE TABLE `address` (
   `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
   `user_id` BIGINT UNSIGNED NOT NULL,
   `receiver_name` VARCHAR(50) NOT NULL COMMENT '收货人姓名',
   `receiver_phone` VARCHAR(20) NOT NULL COMMENT '收货人电话',
   `province` VARCHAR(50) DEFAULT '' COMMENT '省',
   `city` VARCHAR(50) DEFAULT '' COMMENT '市',
   `district` VARCHAR(50) DEFAULT '' COMMENT '区',
   `detail` VARCHAR(200) DEFAULT '' COMMENT '详细地址',
   `is_default` TINYINT DEFAULT 0 COMMENT '1=默认地址',
   `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
   `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
   PRIMARY KEY (`id`),
   KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='收货地址表';