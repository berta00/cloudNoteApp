USE auth;

DELETE FROM users     WHERE id > 0;
DELETE FROM emailConf WHERE id > 0;
