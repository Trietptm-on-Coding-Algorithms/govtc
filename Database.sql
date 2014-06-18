CREATE TABLE `hashes` (
	`id`	INTEGER PRIMARY KEY AUTOINCREMENT,
	`md5`	TEXT,
	`sha256`	TEXT,
	`positives`	INTEGER,
	`total`	INTEGER,
	`permalink`	TEXT,
	`responsecode`	INTEGER,
	`scans`	TEXT,
	`scandate`	NUMERIC,
	`updatedate`	NUMERIC
);

CREATE  INDEX `idx_md5` ON `hashes` (`md5` ASC);
CREATE  INDEX `idx_sha256` ON `hashes` (`sha256` ASC);