-- Create the database if it doesn't exist
CREATE DATABASE IF NOT EXISTS `app-db` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Use the app-db database
USE `app-db`;

-- Create a table for users
CREATE TABLE IF NOT EXISTS `users` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `username` VARCHAR(255) NOT NULL,
    `password` VARCHAR(255) NOT NULL,
    `role` VARCHAR(50) NOT NULL,
    `name` VARCHAR(100) NOT NULL
);

-- Create a table for user sessions
CREATE TABLE IF NOT EXISTS `user_sessions` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `user_id` INT NOT NULL,
    `token` VARCHAR(255) NOT NULL,
    FOREIGN KEY (`user_id`) REFERENCES `users`(`id`)
);

-- Insert a new user if a user with the username "root" doesn't already exist
INSERT IGNORE INTO `users` (`username`, `password`, `role`, `name`)
VALUES ('root', 'root', 'admin', 'Mr Admin Administrador');
