ALTER TABLE `quickslot_keymap`
  ADD COLUMN `petConsumeItem` INT(11) NOT NULL DEFAULT '0' AFTER `key2`,
  ADD COLUMN `petConsumeMPItem` INT(11) NOT NULL DEFAULT '0' AFTER `petConsumeItem`;
