CREATE TABLE `consent_template` (
    `template_name`    varchar(100) NOT NULL,
    `version`          varchar(100) NOT NULL,
    `id`               bigint       NOT NULL AUTO_INCREMENT,
    `create_timestamp` timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
);

CREATE TABLE `last_updated_template` (
    `template_name` varchar(100) NOT NULL,
    `version` varchar(100) NOT NULL,
    `id` bigint NOT NULL AUTO_INCREMENT,
    `create_timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `consent_template_id` bigint NOT NULL,
    PRIMARY KEY (`id`)
);