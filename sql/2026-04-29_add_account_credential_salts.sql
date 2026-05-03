ALTER TABLE `accounts`
  ADD COLUMN IF NOT EXISTS `passwordSalt` varchar(32) DEFAULT NULL AFTER `password`,
  ADD COLUMN IF NOT EXISTS `pinSalt` varchar(32) DEFAULT NULL AFTER `pin`;
