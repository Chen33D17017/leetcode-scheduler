USE `leetcode_scheduler`;

SET NAMES utf8;
set character_set_client = utf8mb4;

DROP TABLE IF EXISTS  `leetcode_user`;
CREATE TABLE `leetcode_user` (
	`user_id` INT AUTO_INCREMENT,
    `user_name` VARCHAR(255) NOT NULL,
    `email` VARCHAR(255) UNIQUE NOT NULL,
    `password` binary(60) NOT NULL,
    `session` VARCHAR(40) UNIQUE,
    PRIMARY KEY (`user_id`)
) ENGINE=INNODB;