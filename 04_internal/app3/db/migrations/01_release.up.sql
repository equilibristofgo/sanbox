CREATE TABLE tasks (
   id          SERIAL NOT NULL PRIMARY KEY,
   description VARCHAR(500) NOT NULL,
   priority    INT DEFAULT 0,
   start_date  TIMESTAMP,
   due_date    TIMESTAMP,
   owner       VARCHAR(100) NOT NULL,
   is_done     BOOLEAN NOT NULL DEFAULT FALSE
);