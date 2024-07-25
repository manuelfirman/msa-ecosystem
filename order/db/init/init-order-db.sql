DROP DATABASE IF EXISTS `ms_order`;
CREATE DATABASE `ms_order`;
USE `ms_order`;


CREATE TABLE `orders` (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    `order_date` DATETIME NOT NULL,
    `user_id` BIGINT NOT NULL,
    `total_price` DECIMAL(10, 2) NOT NULL,
    `status` ENUM('PENDING', 'SHIPPED', 'DELIVERED') NOT NULL,
    `payment_method` ENUM('CREDIT_CARD', 'DEBIT_CARD', 'PAYPAL', 'CASH') NOT NULL,
    `payment_status` ENUM('PAID', 'UNPAID') NOT NULL
);


CREATE TABLE `order_items` (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    `order_id` BIGINT,
    `product_id` BIGINT,
    `quantity` INT NOT NULL,
    `price` DECIMAL(10, 2) NOT NULL,
    FOREIGN KEY (`order_id`) REFERENCES `orders`(id)
);
