ALTER TABLE `items`
    ADD COLUMN IF NOT EXISTS `ringID` INT(11) DEFAULT NULL AFTER `cashSN`;

CREATE TABLE IF NOT EXISTS `marriages` (
    `id` INT(11) NOT NULL AUTO_INCREMENT,
    `husbandID` INT(11) NOT NULL DEFAULT '0',
    `wifeID` INT(11) NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `rings` (
    `id` INT(11) NOT NULL AUTO_INCREMENT,
    `itemID` INT(11) NOT NULL DEFAULT '0',
    `ownerCharacterID` INT(11) NOT NULL DEFAULT '0',
    `partnerRingID` INT(11) NOT NULL DEFAULT '0',
    `partnerCharacterID` INT(11) NOT NULL DEFAULT '0',
    `partnerName` VARCHAR(255) NOT NULL,
    `ringType` TINYINT(4) NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`),
    KEY `idx_ownerCharacterID` (`ownerCharacterID`),
    KEY `idx_partnerCharacterID` (`partnerCharacterID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

ALTER TABLE `account_storage_items`
    ADD COLUMN IF NOT EXISTS `ringID` INT(11) DEFAULT NULL AFTER `creatorName`;

ALTER TABLE `account_cashshop_storage_items`
    ADD COLUMN IF NOT EXISTS `ringID` INT(11) DEFAULT NULL AFTER `sn`;
