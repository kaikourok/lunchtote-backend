ALTER TYPE relation_permission RENAME TO relation_permission_old;
CREATE TYPE relation_permission AS ENUM (
    'DISALLOW',
    'FOLLOW',
    'FOLLOWED',
    'MUTUAL_FOLLOW',
    'ALL'
    );

UPDATE characters SET book_permission = 'DISALLOW' WHERE book_permission = 'REPLY';
UPDATE rooms_messages SET reply_permission = 'DISALLOW' WHERE reply_permission = 'REPLY';

ALTER TABLE characters
    ALTER COLUMN book_permission TYPE relation_permission_old USING relation_permission_old::text::relation_permission;
ALTER TABLE rooms_messages
    ALTER COLUMN reply_permission TYPE relation_permission_old USING relation_permission_old::text::relation_permission;

DROP TYPE relation_permission_old;
