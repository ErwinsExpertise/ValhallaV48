DELETE b1
FROM `buddy` b1
JOIN `buddy` b2
  ON b1.`characterID` = b2.`characterID`
 AND b1.`friendID` = b2.`friendID`
 AND b1.`id` > b2.`id`;

ALTER TABLE `buddy`
  MODIFY `accepted` tinyint(4) NOT NULL DEFAULT '0' COMMENT '0 is request pending, 1 is accepted',
  ADD UNIQUE KEY `idx_buddy_pair` (`characterID`, `friendID`);
