/*
 Navicat Premium Dump SQL

 Source Server         : localhost
 Source Server Type    : MySQL
 Source Server Version : 90300 (9.3.0)
 Source Host           : localhost:3306
 Source Schema         : web3_task3

 Target Server Type    : MySQL
 Target Server Version : 90300 (9.3.0)
 File Encoding         : 65001

 Date: 29/05/2025 13:46:59
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for books
-- ----------------------------
DROP TABLE IF EXISTS `books`;
CREATE TABLE `books`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `author` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `price` decimal(10, 2) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 7 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of books
-- ----------------------------
INSERT INTO `books` VALUES (1, 'Go in Action', 'William Kennedy', 45.00);
INSERT INTO `books` VALUES (2, 'Advanced Go Programming', 'Cai Go', 60.00);
INSERT INTO `books` VALUES (3, 'Database Design for Mere Mortals', 'Michael J. Hernandez', 70.00);
INSERT INTO `books` VALUES (4, 'The Pragmatic Programmer', 'Andrew Hunt', 35.00);
INSERT INTO `books` VALUES (5, 'Clean Code', 'Robert C. Martin', 55.00);
INSERT INTO `books` VALUES (6, 'Design Patterns: Elements of Reusable Object-Oriented Software', 'Erich Gamma', 80.00);

SET FOREIGN_KEY_CHECKS = 1;
