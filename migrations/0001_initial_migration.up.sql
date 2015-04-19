CREATE TABLE convos (
    id            SERIAL         PRIMARY KEY,
    parent_id     INTEGER        NOT NULL,
    sender_id     INTEGER        NOT NULL,
    recipient_id  INTEGER        NOT NULL,
    subject       VARCHAR(140)   NOT NULL,
    body          VARCHAR(64000) NOT NULL,
    FOREIGN KEY (parent_id) REFERENCES convos(id) ON DELETE CASCADE
);