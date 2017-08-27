CREATE TABLE `t_item` (
    `item_id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `item_name` VARCHAR(128),
    `status` VARCHAR(128),
    `pri` INT(5),
    `type_id` INTEGER,
    CONSTRAINT fk_type FOREIGN KEY(`type_id`) REFERENCES t_type(`type_id`)
);
