CREATE DATABASE `ecommercedb` 


CREATE TABLE `ecommercedb`.`address` (
    id INT AUTO_INCREMENT PRIMARY KEY, 
    userId INT NOT NULL, 
    houseNo VARCHAR(255) NOT NULL, 
    landmark VARCHAR(255),
    city VARCHAR(100) NOT NULL,
    state VARCHAR(100) NOT NULL, 
    pincode VARCHAR(10) NOT NULL, 
    phoneNumber VARCHAR(15), 
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (userId) REFERENCES `ecommercedb`.`users`(userId) ON DELETE CASCADE 
);


CREATE TABLE `ecommercedb`.`cart_items` (
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `cart_id` INT NOT NULL,         -- Relates to the cart
    `product_id` INT NOT NULL,      -- Relates to the product
    `quantity` INT NOT NULL,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`cart_id`) REFERENCES `ecommercedb`.`carts`(`id`) ON DELETE CASCADE,
    FOREIGN KEY (`product_id`) REFERENCES `ecommercedb`.`products`(`productId`) ON DELETE CASCADE,
    UNIQUE KEY `unique_cart_product` (`cart_id`, `product_id`) -- Enforce uniqueness
);

  
  CREATE TABLE `ecommercedb`.`carts` (
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `user_id` INT UNIQUE, -- Each user has one cart
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`id`) REFERENCES `ecommercedb`.`users`(`userId`)
);


CREATE TABLE `ecommercedb`.`order_items` (
    `orderItemId` INT PRIMARY KEY AUTO_INCREMENT,
    `orderId` INT NOT NULL,
    `productId` INT NOT NULL, 
    `quantity` INT NOT NULL, 
    `priceperunit` DECIMAL(10, 2) NOT NULL,
    `totalPrice` DECIMAL(10, 2) NOT NULL,
    `createdAt` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updatedAt` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`orderId`) REFERENCES `ecommercedb`.`orders`(`orderId`) ON DELETE CASCADE,
    FOREIGN KEY (`productId`) REFERENCES `ecommercedb`.`products`(`productId`) ON DELETE CASCADE  -- Foreign key to the products table
);


CREATE TABLE `ecommercedb`.`orders` (
    orderId INT AUTO_INCREMENT PRIMARY KEY,
    userId INT NOT NULL,   
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, 
    deliveryMode VARCHAR(100) NOT NULL,       -- Mode of delivery (e.g., "Standard", "Express")
    paymentMode VARCHAR(50) NOT NULL,         -- Payment mode (e.g., "COD")
    orderValue DECIMAL(10, 2) NOT NULL,       -- Total value of the order (before shipping)
    shippingAddress VARCHAR(1024) NOT NULL,
    orderTotal DECIMAL(10, 2) NOT NULL,       -- Total value of the order (including shipping)
    FOREIGN KEY (userId) REFERENCES `ecommercedb`.`users`(userId) ON DELETE CASCADE  -- Foreign key referencing users table
);


CREATE TABLE `ecommercedb`.`products` (
  `productId` INT NOT NULL AUTO_INCREMENT,  
  `pricePerUnit` DECIMAL(13,2) NOT NULL,  
  `productName` VARCHAR(255) NOT NULL,
  `productBrand` VARCHAR(255) NOT NULL,
  `description` VARCHAR(255) NOT NULL DEFAULT "This is a sample description",
  `stockQuantity` INT NOT NULL DEFAULT 10,
  PRIMARY KEY (`productId`));


CREATE TABLE `ecommercedb`.`users` (
  `userId` INT NOT NULL AUTO_INCREMENT,  
  `email` VARCHAR(255) NOT NULL,  
  `password` VARCHAR(255) NOT NULL,
  `isAdmin` TINYINT(1) NOT NULL DEFAULT 0,
  PRIMARY KEY (`userId`));
