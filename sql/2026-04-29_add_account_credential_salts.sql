ALTER TABLE `accounts`
  ADD COLUMN `passwordSalt` varchar(32) DEFAULT NULL AFTER `password`,
  ADD COLUMN `pinSalt` varchar(32) DEFAULT NULL AFTER `pin`;
