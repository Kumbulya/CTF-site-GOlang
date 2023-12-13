CREATE TABLE `users` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `login` varchar(50) NOT NULL,
    `password` varchar(50) NOT NULL,
    `page` varchar(50) NOT NULL,
    `isAdmin` tinyint(1),
    `balance` float,
    PRIMARY KEY (`id`),
    UNIQUE KEY (`id`)
) AUTO_INCREMENT = 1 ;

INSERT INTO `users` (`login`, `password`, `page`, `isAdmin`, `balance`) VALUES ('admin','admin','1', 1, 0);

CREATE TABLE `katalog`(
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `product_name` varchar(50) NOT NULL,
    `category` varchar(50) NOT NULL,
    `seller` varchar(50) NOT NULL,
    `description` varchar(255) NOT NULL,
    `cost` float NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY (`id`)
)   AUTO_INCREMENT = 1 ;

INSERT INTO `katalog` (`product_name`, `category`,`seller`, `description`, `cost`) VALUES ('Senko-san','Fox-wife','1', 'This very cute little fox will always be waiting for you at home.', 9999);

CREATE TABLE `basket`(
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `basketID` int(11),
    `productID` int(11),
    PRIMARY KEY (`id`),
    UNIQUE KEY (`id`)
)   AUTO_INCREMENT = 1 ;