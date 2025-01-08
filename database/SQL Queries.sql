CREATE DATABASE `ecommercedb`;

CREATE TABLE `ecommercedb`.`products` (
  `productId` INT NOT NULL AUTO_INCREMENT,  
  `pricePerUnit` DECIMAL(13,2) NOT NULL,  
  `productName` VARCHAR(255) NOT NULL,
  `productBrand` VARCHAR(255) NOT NULL,
  `description` VARCHAR(255) NOT NULL DEFAULT "This is a sample description",
  `stockQuantity` INT NOT NULL DEFAULT 10,
  `category` VARCHAR(255) NOT NULL DEFAULT "household",
  `subCategory` VARCHAR(255) NOT NULL DEFAULT "household",
  `imageURL` VARCHAR(255) NOT NULL DEFAULT "/images/sample-image.jpg",
  PRIMARY KEY (`productId`));

INSERT INTO `ecommercedb`.`products` (`pricePerUnit`, `productName`, `productBrand`) VALUES (497.45,"sticky note","Johns-Jenkins");  
INSERT INTO `ecommercedb`.`products` (`pricePerUnit`, `productName`, `productBrand`) VALUES (282.29,"leg warmers","Hessel, Schimmel and Feeney");  
INSERT INTO `ecommercedb`.`products` (`pricePerUnit`, `productName`, `productBrand`) VALUES (436.26,"lamp shade","Swaniawski, Bartoletti and Bruen");  
INSERT INTO `ecommercedb`.`products` (`pricePerUnit`, `productName`, `productBrand`) VALUES (537.90,"flowers","Runolfsdottir, Littel and Dicki");  
INSERT INTO `ecommercedb`.`products` (`pricePerUnit`, `productName`, `productBrand`) VALUES (112.10,"clamp","Kuhn, Cronin and Spencer");  
INSERT INTO `ecommercedb`.`products` (`pricePerUnit`, `productName`, `productBrand`) VALUES (593.53,"twister","Quigley, Casper and Boyer");  
INSERT INTO `ecommercedb`.`products` (`pricePerUnit`, `productName`, `productBrand`) VALUES (88.97,"clay pot","Gutmann and Sons");  
INSERT INTO `ecommercedb`.`products` (`pricePerUnit`, `productName`, `productBrand`) VALUES (933.35,"tooth picks","Bins-Hansen");  
INSERT INTO `ecommercedb`.`products` (`pricePerUnit`, `productName`, `productBrand`) VALUES (426.23,"mirror","Jones, Braun and Ratke");  
INSERT INTO `ecommercedb`.`products` (`pricePerUnit`, `productName`, `productBrand`) VALUES (630.61,"rug","Upton-Mraz");  
INSERT INTO `ecommercedb`.`products` (`pricePerUnit`, `productName`, `productBrand`) VALUES (13.67,"headphones","Schneider, Douglas and Franecki");


CREATE TABLE `ecommercedb`.`variantProducts` (
  `variantId` VARCHAR(255) PRIMARY KEY,  
  `productId` INT NOT NULL,
  `variantName` VARCHAR(255) NOT NULL,
  `description` VARCHAR(255) NOT NULL DEFAULT "This is a sample variant description",
  `pricePerUnit` DECIMAL(13,2) NOT NULL,
  `stockQuantity` INT NOT NULL DEFAULT 10,
  `imageURL` VARCHAR(255) NOT NULL DEFAULT "/images/sample-image.jpg",
  `color` VARCHAR(255) NOT NULL DEFAULT "Any",
  FOREIGN KEY (`productId`) REFERENCES `ecommercedb`.`products`(`productId`) ON DELETE CASCADE
  );

INSERT INTO `ecommercedb`.`variantProducts` (`variantId`, `productId`, `variantName`, `description`, `pricePerUnit`, `color`)
VALUES
-- Product 1
("1A", 1, "sticky note_Red", "This is a Red Variant", 497.45, "Red"),
("1B", 1, "sticky note_Blue", "This is a Blue Variant", 497.45, "Blue"),
("1C", 1, "sticky note_Black", "This is a Black Variant", 497.45, "Black"),

-- Product 2
("2A", 2, "leg warmers_Red", "This is a Red Variant", 700.29, "Red"),
("2B", 2, "leg warmers_Blue", "This is a Blue Variant", 700.29, "Blue"),
("2C", 2, "leg warmers_Black", "This is a Black Variant", 700.29, "Black"),

-- Product 3
("3A", 3, "lamp shade_Red", "This is a Red Variant", 436.26, "Red"),
("3B", 3, "lamp shade_Blue", "This is a Blue Variant", 436.26, "Blue"),
("3C", 3, "lamp shade_Black", "This is a Black Variant", 436.26, "Black"),

-- Product 4
("4A", 4, "flowers_Red", "This is a Red Variant", 537.9, "Red"),
("4B", 4, "flowers_Blue", "This is a Blue Variant", 537.9, "Blue"),
("4C", 4, "flowers_Black", "This is a Black Variant", 537.9, "Black"),

-- Product 5
("5A", 5, "clamp_Red", "This is a Red Variant", 112.1, "Red"),
("5B", 5, "clamp_Blue", "This is a Blue Variant", 112.1, "Blue"),
("5C", 5, "clamp_Black", "This is a Black Variant", 112.1, "Black"),

-- Product 6
("6A", 6, "twister_Red", "This is a Red Variant", 593.53, "Red"),
("6B", 6, "twister_Blue", "This is a Blue Variant", 593.53, "Blue"),
("6C", 6, "twister_Black", "This is a Black Variant", 593.53, "Black"),

-- Product 7
("7A", 7, "clay pot_Red", "This is a Red Variant", 88.97, "Red"),
("7B", 7, "clay pot_Blue", "This is a Blue Variant", 88.97, "Blue"),
("7C", 7, "clay pot_Black", "This is a Black Variant", 88.97, "Black"),

-- Product 8
("8A", 8, "tooth picks_Red", "This is a Red Variant", 933.35, "Red"),
("8B", 8, "tooth picks_Blue", "This is a Blue Variant", 933.35, "Blue"),
("8C", 8, "tooth picks_Black", "This is a Black Variant", 933.35, "Black"),

-- Product 9
("9A", 9, "mirror_Red", "This is a Red Variant", 426.23, "Red"),
("9B", 9, "mirror_Blue", "This is a Blue Variant", 426.23, "Blue"),
("9C", 9, "mirror_Black", "This is a Black Variant", 426.23, "Black"),

-- Product 10
("10A", 10, "rug_Red", "This is a Red Variant", 630.61, "Red"),
("10B", 10, "rug_Blue", "This is a Blue Variant", 630.61, "Blue"),
("10C", 10, "rug_Black", "This is a Black Variant", 630.61, "Black"),

-- Product 11
("11A", 11, "headphones_Red", "This is a Red Variant", 13.67, "Red"),
("11B", 11, "headphones_Blue", "This is a Blue Variant", 13.67, "Blue"),
("11C", 11, "headphones_Black", "This is a Black Variant", 13.67, "Black");

CREATE TABLE `ecommercedb`.`users` (
  `userId` INT NOT NULL AUTO_INCREMENT,  
  `email` VARCHAR(255) NOT NULL,  
  `password` VARCHAR(255) NOT NULL,
  `name` VARCHAR(255) DEFAULT "SAMPLE_NAME",
  `isAdmin` TINYINT(1) NOT NULL DEFAULT 0,
  PRIMARY KEY (`userId`));


CREATE TABLE `ecommercedb`.`address` (
    `id` INT AUTO_INCREMENT PRIMARY KEY, 
    `userId` INT NOT NULL, 
    `houseNo` VARCHAR(255) NOT NULL, 
    `landmark` VARCHAR(255),
    `city` VARCHAR(100) NOT NULL,
    `state` VARCHAR(100) NOT NULL, 
    `pincode` VARCHAR(10) NOT NULL, 
    `phoneNumber` VARCHAR(15), 
    `createdAt` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updatedAt` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`userId`) REFERENCES `ecommercedb`.`users`(`userId`) ON DELETE CASCADE 
);


CREATE TABLE `ecommercedb`.`carts` (
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `user_id` INT UNIQUE, -- Each user has one cart
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`id`) REFERENCES `ecommercedb`.`users`(`userId`)
);


CREATE TABLE `ecommercedb`.`cart_items` (
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `cart_id` INT NOT NULL,         -- Relates to the cart
    `variantId` VARCHAR(255) NOT NULL,      -- Relates to the variant product
    `quantity` INT NOT NULL,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`cart_id`) REFERENCES `ecommercedb`.`carts`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`variantId`) REFERENCES `ecommercedb`.`variantProducts`(`variantId`) ON DELETE CASCADE,
    UNIQUE KEY `unique_cart_product` (`cart_id`, `variantId`) -- Enforce uniqueness
);


CREATE TABLE `ecommercedb`.`orders` (
    `orderId` INT AUTO_INCREMENT PRIMARY KEY,
    `userId` INT NOT NULL,   
    `createdAt` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updatedAt` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, 
    `deliveryMode` VARCHAR(100) NOT NULL,       -- Mode of delivery (e.g., "Standard", "Express")
    `paymentMode` VARCHAR(50) NOT NULL,         -- Payment mode (e.g., "COD")
    `orderValue` DECIMAL(10, 2) NOT NULL,       -- Total value of the order (before shipping)
    `shippingAddress` VARCHAR(1024) NOT NULL,
    `orderTotal` DECIMAL(10, 2) NOT NULL,       -- Total value of the order (including shipping)
    FOREIGN KEY (`userId`) REFERENCES `ecommercedb`.`users`(`userId`) ON DELETE CASCADE  -- Foreign key referencing users table
);

CREATE TABLE `ecommercedb`.`order_items` (
    `orderItemId` INT PRIMARY KEY AUTO_INCREMENT,
    `orderId` INT NOT NULL,
    `variantId` INT NOT NULL, 
    `quantity` INT NOT NULL, 
    `priceperunit` DECIMAL(10, 2) NOT NULL,
    `totalPrice` DECIMAL(10, 2) NOT NULL,
    `createdAt` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updatedAt` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`orderId`) REFERENCES `ecommercedb`.`orders`(`orderId`) ON DELETE CASCADE,
    FOREIGN KEY (`variantId`) REFERENCES `ecommercedb`.`variantProducts`(`variantId`) ON DELETE CASCADE  -- Foreign key to the products table
);

CREATE TABLE `ecommercedb`.`productPrices` (
    productId INT NOT NULL,                    -- Unique identifier for the product (foreign key from products table)
    currencyCode VARCHAR(3) NOT NULL,           -- Currency code (e.g., 'USD', 'EUR', 'GBP')
    price DECIMAL(10, 2) NOT NULL,              -- Price of the product in the specified currency
    PRIMARY KEY (productId, currencyCode),     -- Composite primary key (productId + currencyCode)
    FOREIGN KEY (productId) REFERENCES `ecommercedb`.`products`(productId) -- productId references products table
    ON DELETE CASCADE                           -- Ensure that if a product is deleted, its price record is also deleted
);

INSERT INTO `ecommercedb`.`productPrices` (productId, currencyCode, price)
VALUES
(1, 'USD', 19.99),
(1, 'EUR', 17.50),
(1, 'GBP', 15.00),
(2, 'USD', 29.99),
(2, 'EUR', 26.50),
(2, 'GBP', 22.00),
(3, 'USD', 99.99),
(3, 'EUR', 89.00),
(3, 'GBP', 79.00),
(4, 'USD', 249.99),
(4, 'EUR', 220.00),
(4, 'GBP', 190.00),
(5, 'USD', 10.50),
(5, 'EUR', 9.00),
(5, 'GBP', 8.00),
(6, 'USD', 199.99),
(6, 'EUR', 179.00),
(6, 'GBP', 160.00),
(7, 'USD', 15.00),
(7, 'EUR', 13.00),
(7, 'GBP', 11.50),
(8, 'USD', 49.99),
(8, 'EUR', 44.00),
(8, 'GBP', 38.00),
(9, 'USD', 9.99),
(9, 'EUR', 8.50),
(9, 'GBP', 7.50),
(10, 'USD', 199.99),
(10, 'EUR', 180.00),
(10, 'GBP', 160.00),
(11, 'USD', 25.00),
(11, 'EUR', 22.00),
(11, 'GBP', 19.50);