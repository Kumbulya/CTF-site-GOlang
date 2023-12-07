CREATE TABLE `users` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `login` varchar(50) NOT NULL,
    `password` varchar(50) NOT NULL,
    `page` varchar(50) NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY (`id`)
) AUTO_INCREMENT = 1 ;

INSERT INTO `users` (`login`, `password`, `page`) VALUES ('admin','admin','1');

CREATE TABLE `katalog`(
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `product_name` varchar(50) NOT NULL,
    `category` varchar(50) NOT NULL,
    `seller` varchar(50) NOT NULL,
    `description` varchar(50) NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY (`id`)
)   AUTO_INCREMENT = 1 ;

