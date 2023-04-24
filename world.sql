-- phpMyAdmin SQL Dump
-- version 5.2.0
-- https://www.phpmyadmin.net/
--
-- Хост: 127.0.0.1
-- Время создания: Апр 24 2023 г., 19:47
-- Версия сервера: 10.4.27-MariaDB
-- Версия PHP: 8.0.25

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- База данных: `world`
--

-- --------------------------------------------------------

--
-- Структура таблицы `comments`
--

CREATE TABLE `comments` (
  `comm_id` int(11) DEFAULT NULL,
  `productid` int(11) DEFAULT NULL,
  `userid` int(11) DEFAULT NULL,
  `comment` varchar(9) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Дамп данных таблицы `comments`
--

INSERT INTO `comments` (`comm_id`, `productid`, `userid`, `comment`) VALUES
(1, 6, 4, 'comment 5'),
(2, 1, 7, 'comment 4'),
(3, 1, 2, 'comment 4'),
(4, 9, 6, 'comment 5'),
(5, 8, 6, 'comment 1'),
(NULL, 0, 5, 'text '),
(NULL, 0, 5, 'text'),
(NULL, 0, 5, 'fghjk'),
(NULL, 3, 5, 'fgh'),
(NULL, 3, 5, 'my commen'),
(NULL, 1, 12, 'my commen'),
(NULL, 1, 12, 'new ones\r'),
(NULL, 10, 5, 'some comm'),
(NULL, 10, 5, 'new one'),
(NULL, 10, 5, 'new one'),
(NULL, 10, 13, 'my commen'),
(NULL, 3, 13, 'ghjkl'),
(NULL, 10, 5, 'new one'),
(NULL, 5, 5, 'new fghj'),
(NULL, 8, 0, 'midterm2 '),
(NULL, 8, 0, 'midterm2 '),
(NULL, 1, 0, 'sdfgh'),
(NULL, 1, 5, 'df'),
(NULL, 1, 5, 'ghjkl;'),
(NULL, 8, 5, 'zhzh'),
(NULL, 10, 5, 'fghjk'),
(NULL, 5, 5, 'fghj'),
(NULL, 6, 5, 'asdfgk'),
(NULL, 6, 5, 'asdfg'),
(NULL, 8, 5, 'asdfgh'),
(NULL, 3, 5, 'asdfg');

-- --------------------------------------------------------

--
-- Структура таблицы `products`
--

CREATE TABLE `products` (
  `id` int(11) DEFAULT NULL,
  `car_name` varchar(50) DEFAULT NULL,
  `details` text DEFAULT NULL,
  `price` int(11) DEFAULT NULL,
  `rating` varchar(10) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Дамп данных таблицы `products`
--

INSERT INTO `products` (`id`, `car_name`, `details`, `price`, `rating`) VALUES
(1, 'GX', 'non pretium quis lectus suspendisse potenti in eleifend quam a odio in hac habitasse', 20, '3.6'),
(2, 'S40', 'venenatis non sodales sed tincidunt eu felis fusce posuere felis sed lacus morbi', 23, '3.33'),
(3, 'Seville', 'elit sodales scelerisque mauris sit amet eros suspendisse accumsan tortor quis turpis sed ante', 79, '3.5'),
(4, 'Turbo Firefly', 'vivamus tortor duis mattis egestas metus aenean fermentum donec ut mauris eget massa tempor convallis', 92, '3'),
(5, 'Passat', 'ipsum dolor sit amet consectetuer adipiscing elit proin risus praesent lectus vestibulum quam sapien varius', 77, '2.83'),
(6, 'Eclipse', 'tincidunt in leo maecenas pulvinar lobortis est phasellus sit amet erat nulla tempus vivamus in felis', 97, '3'),
(7, 'Town Car', 'nulla suscipit ligula in lacus curabitur at ipsum ac tellus', 91, '3'),
(8, 'Hearse', 'luctus rutrum nulla tellus in sagittis dui vel nisl duis ac nibh fusce lacus purus aliquet at feugiat non', 92, '3.25'),
(9, 'Sierra 3500', 'ultrices phasellus id sapien in sapien iaculis congue vivamus metus arcu adipiscing molestie hendrerit at vulputate vitae', 37, '3.5'),
(10, 'ES', 'nec sem duis aliquam convallis nunc proin at turpis a pede', 100, '3.6');

-- --------------------------------------------------------

--
-- Структура таблицы `ratings`
--

CREATE TABLE `ratings` (
  `rate` varchar(1) DEFAULT NULL,
  `productId` int(11) DEFAULT NULL,
  `userId` int(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Дамп данных таблицы `ratings`
--

INSERT INTO `ratings` (`rate`, `productId`, `userId`) VALUES
('5', 9, 6),
('4', 10, 6),
('1', 2, 1),
('3', 8, 7),
('4', 10, 5),
('4', 6, 4),
('2', 8, 3),
('5', 3, 9),
('2', 5, 4),
('4', 3, 8),
('4', 5, 7),
('3', 1, 5),
('4', 4, 1),
('3', 1, 7),
('1', 4, 8),
('1', 5, 2),
('2', 3, 1),
('5', 2, 2),
('2', 9, 4),
('1', 6, 4),
('1', 7, 3),
('4', 1, 5),
('4', 1, 0),
('4', 1, 12),
('4', 2, 12),
('2', 3, 12),
('3', 10, 5),
('4', 10, 13),
('3', 3, 13),
('3', 10, 5),
('5', 7, 5),
('4', 5, 5),
('4', 6, 5),
('4', 8, 5),
('4', 4, 5),
('4', 8, 5),
('1', 5, 5),
('5', 5, 5),
('5', 3, 5);

--
-- Триггеры `ratings`
--
DELIMITER $$
CREATE TRIGGER `update_product_rating` AFTER INSERT ON `ratings` FOR EACH ROW BEGIN
  UPDATE products p
  SET p.rating = (SELECT ROUND(AVG(rate),2) FROM ratings WHERE productId = NEW.productId)
  WHERE p.id = NEW.productId;
END
$$
DELIMITER ;

-- --------------------------------------------------------

--
-- Структура таблицы `users`
--

CREATE TABLE `users` (
  `userid` int(11) NOT NULL,
  `name` varchar(255) NOT NULL,
  `password` varchar(255) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Дамп данных таблицы `users`
--

INSERT INTO `users` (`userid`, `name`, `password`) VALUES
(2, 'zhaniya', 'd'),
(3, 'zhaniya2', 's'),
(4, 'Zhaniy2', 'z'),
(5, 'zh', 'zh'),
(6, 'zh1', 'zh'),
(7, 'z1', 'z'),
(8, 'z2', 'z'),
(9, 'z3', 'z'),
(10, 'asd', 'asd'),
(11, 'as', 'as'),
(12, 'new', 'new'),
(13, 'new4', 'new4');

--
-- Индексы сохранённых таблиц
--

--
-- Индексы таблицы `users`
--
ALTER TABLE `users`
  ADD PRIMARY KEY (`userid`);

--
-- AUTO_INCREMENT для сохранённых таблиц
--

--
-- AUTO_INCREMENT для таблицы `users`
--
ALTER TABLE `users`
  MODIFY `userid` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=14;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
