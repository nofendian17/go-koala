-- Adminer 4.7.7 MySQL dump
use db_app;
SET NAMES utf8;
SET time_zone = '+00:00';
SET foreign_key_checks = 0;
SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';

DROP TABLE IF EXISTS `customers`;
CREATE TABLE `customers`
(
    `customer_id`   varchar(64)  NOT NULL,
    `customer_name` varchar(80)  NOT NULL,
    `email`         varchar(50)  NOT NULL,
    `phone_number`  varchar(20)  NOT NULL,
    `dob`           date         NOT NULL,
    `sex`           tinyint(1)   NOT NULL DEFAULT '1',
    `salt`          varchar(80)  NOT NULL,
    `password`      varchar(400) NOT NULL,
    `created_at`    datetime     NOT NULL,
    PRIMARY KEY (`customer_id`),
    UNIQUE KEY `customers_customer_id_uindex` (`customer_id`),
    UNIQUE KEY `customers_email_uindex` (`email`)
) ENGINE = InnoDB
  DEFAULT CHARSET = latin1;

INSERT INTO `customers` (`customer_id`, `customer_name`, `email`, `phone_number`, `dob`, `sex`, `salt`, `password`,
                         `created_at`)
VALUES ('489ad21d-2569-4c06-958d-6677da1ebc60', 'customer 1', 'customer1@example.com', '08123456789', '2020-12-04', 1,
        'secret_salt', '123456', '2020-12-04 12:57:13'),
       ('4cd4c6eb-2ba0-4d8a-967b-94883a4d3d11', 'customer 2', 'customer2@example.com', '08123456789', '2020-12-04', 1,
        'secret_salt', '123456', '2020-12-04 12:57:13'),
       ('d0eaa74b-4bab-4fbc-aed4-e720e480d813', 'customer 3', 'customer3@example.com', '08123456789', '2020-12-04', 1,
        'secret_salt', '123456', '2020-12-04 12:57:13');

DROP TABLE IF EXISTS `orders`;
CREATE TABLE `orders`
(
    `order_id`          varchar(64) NOT NULL,
    `customer_id`       varchar(64) NOT NULL,
    `order_number`      varchar(40) NOT NULL,
    `order_date`        datetime    NOT NULL,
    `payment_method_id` varchar(64) NOT NULL,
    PRIMARY KEY (`order_id`),
    UNIQUE KEY `orders_order_id_uindex` (`order_id`),
    KEY `orders_customers_customer_id_fk` (`customer_id`),
    KEY `orders_payment_methods_payment_method_id_fk` (`payment_method_id`),
    CONSTRAINT `orders_customers_customer_id_fk` FOREIGN KEY (`customer_id`) REFERENCES `customers` (`customer_id`),
    CONSTRAINT `orders_payment_methods_payment_method_id_fk` FOREIGN KEY (`payment_method_id`) REFERENCES `payment_methods` (`payment_method_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = latin1;

INSERT INTO `orders` (`order_id`, `customer_id`, `order_number`, `order_date`, `payment_method_id`)
VALUES ('2bc9e664-26c8-45dd-abe0-e386fe589173', 'd0eaa74b-4bab-4fbc-aed4-e720e480d813', 'PO-003/IX/2020',
        '2020-12-04 13:12:30', 'f693dced-cc1d-4aa9-baa3-978eb3148a73'),
       ('7ece3018-a4b7-4a3d-9913-875afc4ac8a3', '489ad21d-2569-4c06-958d-6677da1ebc60', 'PO-004/IX/2020',
        '2020-12-04 13:12:30', 'f693dced-cc1d-4aa9-baa3-978eb3148a73'),
       ('7f841b5f-ed8d-4619-9e79-8274238beb06', '489ad21d-2569-4c06-958d-6677da1ebc60', 'PO-001/IX/2020',
        '2020-12-04 13:12:30', '017c30fb-8e64-4b5d-ac5a-7ba335a1cb36'),
       ('f325cfe7-7647-4f8b-bb70-26714419fe8a', '4cd4c6eb-2ba0-4d8a-967b-94883a4d3d11', 'PO-002/IX/2020',
        '2020-12-04 13:12:30', 'f693dced-cc1d-4aa9-baa3-978eb3148a73');

DROP TABLE IF EXISTS `order_details`;
CREATE TABLE `order_details`
(
    `order_detail_id` varchar(64) NOT NULL,
    `order_id`        varchar(64) NOT NULL,
    `product_id`      varchar(64) NOT NULL,
    `qty`             int(11)     NOT NULL,
    `created_at`      datetime    NOT NULL,
    PRIMARY KEY (`order_detail_id`),
    UNIQUE KEY `order_details_order_detail_id_uindex` (`order_detail_id`),
    KEY `order_details_orders_order_id_fk` (`order_id`),
    KEY `order_details_products_product_id_fk` (`product_id`),
    CONSTRAINT `order_details_orders_order_id_fk` FOREIGN KEY (`order_id`) REFERENCES `orders` (`order_id`),
    CONSTRAINT `order_details_products_product_id_fk` FOREIGN KEY (`product_id`) REFERENCES `products` (`product_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = latin1;

INSERT INTO `order_details` (`order_detail_id`, `order_id`, `product_id`, `qty`, `created_at`)
VALUES ('4a41d3ea-4608-4153-90d0-ea79a8b92b6d', 'f325cfe7-7647-4f8b-bb70-26714419fe8a',
        '4bf36408-82da-475a-a3da-86a610f70f87', 1, '2020-12-04 13:14:18'),
       ('51cf3afa-6a54-4fdd-849a-7dc0bf54e88c', '7f841b5f-ed8d-4619-9e79-8274238beb06',
        '238868ab-57f6-4674-89fe-e6d6dd3ebd92', 1, '2020-12-04 13:14:18'),
       ('99721687-8e58-45a1-909d-0364524a348a', '2bc9e664-26c8-45dd-abe0-e386fe589173',
        '542c6767-95b9-4740-8667-14f2c0cf2199', 2, '2020-12-04 13:14:18'),
       ('9e7b2bd4-3641-4a3e-889d-18fb5fc5dd22', '7f841b5f-ed8d-4619-9e79-8274238beb06',
        '4bf36408-82da-475a-a3da-86a610f70f87', 1, '2020-12-04 13:14:18'),
       ('a6d05890-aa6b-4fdc-bb31-3d7c02970738', 'f325cfe7-7647-4f8b-bb70-26714419fe8a',
        '542c6767-95b9-4740-8667-14f2c0cf2199', 2, '2020-12-04 13:14:18'),
       ('aea15d27-27a3-48fa-856e-e5f229055b4c', '7ece3018-a4b7-4a3d-9913-875afc4ac8a3',
        '4bf36408-82da-475a-a3da-86a610f70f87', 2, '2020-12-04 13:14:18');

DROP TABLE IF EXISTS `payment_methods`;
CREATE TABLE `payment_methods`
(
    `payment_method_id` varchar(64) NOT NULL,
    `method_name`       varchar(70) NOT NULL,
    `code`              varchar(10) NOT NULL,
    `created_at`        datetime    NOT NULL,
    PRIMARY KEY (`payment_method_id`),
    UNIQUE KEY `payment_methods_code_uindex` (`code`),
    UNIQUE KEY `payment_methods_payment_method_id_uindex` (`payment_method_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = latin1;

INSERT INTO `payment_methods` (`payment_method_id`, `method_name`, `code`, `created_at`)
VALUES ('017c30fb-8e64-4b5d-ac5a-7ba335a1cb36', 'Virtual Account', 'VA', '2020-12-04 13:07:06'),
       ('3e857f69-d943-4ece-94dc-5e5fb0fc6b8e', 'Credit Card', 'CC', '2020-12-04 13:08:28'),
       ('f693dced-cc1d-4aa9-baa3-978eb3148a73', 'Bank Transfer', 'BT', '2020-12-04 13:06:16');

DROP TABLE IF EXISTS `products`;
CREATE TABLE `products`
(
    `product_id`   varchar(64) NOT NULL,
    `product_name` varchar(80) NOT NULL,
    `basic_price`  decimal(10, 2) DEFAULT NULL,
    `created_at`   datetime    NOT NULL,
    UNIQUE KEY `products_product_id_uindex` (`product_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = latin1;

INSERT INTO `products` (`product_id`, `product_name`, `basic_price`, `created_at`)
VALUES ('238868ab-57f6-4674-89fe-e6d6dd3ebd92', 'grape', 20000.00, '2020-12-04 13:00:24'),
       ('4bf36408-82da-475a-a3da-86a610f70f87', 'banana', 15000.00, '2020-12-04 13:00:24'),
       ('542c6767-95b9-4740-8667-14f2c0cf2199', 'apple', 10000.00, '2020-12-04 13:00:24');

-- 2020-12-04 06:51:00
