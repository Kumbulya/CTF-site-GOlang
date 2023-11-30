CREATE TABLE `users` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `login` varchar(50) NOT NULL,
    `password` varchar(50) NOT NULL,
    `page` varchar(50) NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY (`id`)
) AUTO_INCREMENT = 1 ;

INSERT INTO `users` (`login`, `password`, `page`) VALUES ('admin','admin','1');



