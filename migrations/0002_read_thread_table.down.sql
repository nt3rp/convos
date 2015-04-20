DROP TABLE read_status;
ALTER TABLE convos DROP CONSTRAINT "convos_recipient_id_fkey";
ALTER TABLE convos DROP CONSTRAINT "convos_sender_id_fkey";
DROP TABLE users;