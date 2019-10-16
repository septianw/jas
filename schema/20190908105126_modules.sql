-- +migrate Up
CREATE TABLE `modules` (
  `mid` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(191) NULL,
  `version` VARCHAR(255) NULL,
  `status` VARCHAR(25) NULL,
  `sopath` TEXT NULL,
  PRIMARY KEY (`mid`),
  CONSTRAINT uc_modules UNIQUE (mid,name))
ENGINE = InnoDB;

-- +migrate Down
DROP TABLE IF EXIST `modules`;