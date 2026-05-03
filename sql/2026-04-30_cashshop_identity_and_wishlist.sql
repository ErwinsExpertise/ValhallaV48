ALTER TABLE `account_storage_items`
    ADD COLUMN IF NOT EXISTS `cashID` BIGINT(20) DEFAULT NULL AFTER `creatorName`,
    ADD COLUMN IF NOT EXISTS `cashSN` INT(11) DEFAULT NULL AFTER `cashID`;

CREATE TABLE IF NOT EXISTS `cashshop_wishlist` (
    `characterID` INT(11) NOT NULL,
    `slot` TINYINT(3) UNSIGNED NOT NULL,
    `sn` INT(11) NOT NULL,
    PRIMARY KEY (`characterID`, `slot`),
    CONSTRAINT `cashshop_wishlist_fk_character`
        FOREIGN KEY (`characterID`) REFERENCES `characters` (`id`)
        ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
