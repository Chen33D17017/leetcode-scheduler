USE `leetcode_scheduler`;

SET NAMES utf8;
set character_set_client = utf8mb4;

DROP TABLE IF EXISTS  `problem_log`;
DROP TABLE IF EXISTS  `problem_category`;
DROP TABLE IF EXISTS  `category`;

CREATE TABLE `problem_log` (
	`id` INT AUTO_INCREMENT,
	`user_id` INT,
	`problem_id` INT,
	`problem_id` INT,
    `date` DATE,
    `review_level` INT,
    `time` CHAR(20),
    `done` BOOLEAN,
    PRIMARY KEY (`id`),
    UNIQUE KEY (`user_id`, `problem_id`, `date`)
) ENGINE=INNODB;


CREATE TABLE `category` (
	`id` INT AUTO_INCREMENT,
    `category_name` CHAR(30),
    PRIMARY KEY (`id`)
) ENGINE=INNODB;


CREATE TABLE `problem_category` (
	`id` INT AUTO_INCREMENT,
    `category_id` INT,
    `problem_id` INT,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`category_id`)
        REFERENCES `category` (`id`),
	FOREIGN KEY (`problem_id`)
        REFERENCES `leetcode_problem` (`id`)
) ENGINE=INNODB;

