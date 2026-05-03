CREATE TABLE IF NOT EXISTS `pending_migrations` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `accountID` int(10) unsigned NOT NULL,
  `characterID` int NOT NULL,
  `worldID` tinyint unsigned NOT NULL DEFAULT '0',
  `destinationType` varchar(16) NOT NULL,
  `destinationID` int NOT NULL,
  `clientIP` varchar(45) NOT NULL DEFAULT '',
  `nonce` varchar(32) NOT NULL,
  `createdAt` bigint NOT NULL,
  `expiresAt` bigint NOT NULL,
  `consumedAt` bigint NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `idx_pending_migration_lookup` (`characterID`,`destinationType`,`destinationID`,`consumedAt`,`expiresAt`),
  KEY `idx_pending_migration_account` (`accountID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
