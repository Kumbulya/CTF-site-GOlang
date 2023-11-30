CREATE TABLE `users` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `login` varchar(50) NOT NULL,
    `password` varchar(50) NOT NULL,
    PRIMARY KEY (`user_id`),
    UNIQUE KEY (`user_id`)
) AUTO_INCREMENT = 1 ;

CREATE TABLE `products`(
    `product_id` int(11) NOT NULL AUTO_INCREMENT,
    `fabric_id` int(11) NOT NULL, 
    `product_name`varchar(50) NOT NULL,
    `product_description` text NOT NULL,
    `product_price` decimal(20,2) NOT NULL,
    `product_image` varchar(255) NOT NULL,
    PRIMARY KEY (`product_id`),
    UNIQUE KEY (`product_id`)
) AUTO_INCREMENT = 1 ; 

CREATE  TABLE `product_properties` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
     `product_id` int(11) NOT NULL,
    `property_name` varchar(255) NOT NULL,
    `property_value` varchar(255) NOT NULL,
    `property_price` decimal(20,2) NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY  (`id`)
) AUTO_INCREMENT=1 ;

CREATE TABLE product_images(
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `product_id` int(11) NOT NULL,
    `image` varchar(255) NOT NULL,
    `title` varchar(255) NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY (`id`)
) AUTO_INCREMENT = 1 ;
