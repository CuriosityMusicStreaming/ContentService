-- +migrate Up
CREATE TABLE content
(
    `content_id` binary(16) NOT NULL,
    `name` varchar(255) NOT NULL,
    `type` smallint(2) NOT NULL,
    PRIMARY KEY (`content_id`),
    INDEX `content_id_index` (`content_id`)
);
-- +migrate Down
DROP TABLE content;