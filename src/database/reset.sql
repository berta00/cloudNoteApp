USE noteapp;

DELETE FROM users     WHERE id >= 0;
DELETE FROM emailConf WHERE id >= 0;
DELETE FROM folder    WHERE id >= 0;
DELETE FROM basicNote WHERE id >= 0;
