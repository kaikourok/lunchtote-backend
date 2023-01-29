package general

func (db *GeneralRepository) Initialize() error {
	db.MustExec(`
    DROP SCHEMA public CASCADE;
    CREATE SCHEMA public;
  `)

	db.MustExec(`
    CREATE EXTENSION pg_bigm;

    CREATE TYPE announcement_type AS ENUM (
      'UPDATE',
      'ANNOUNCE',
      'IMPORTANT'
    );

    CREATE TYPE notification_type AS ENUM (
      'FOLLOWED',
      'REPLIED'
    );

    CREATE TYPE prohibition_related_data_type AS ENUM (
      'BAN',
      'UNBAN',
      'ADOMONISH'
    );

    CREATE TYPE relation_permission AS ENUM (
      'DISALLOW',
      'FOLLOW',
      'FOLLOWED',
      'MUTUAL_FOLLOW',
      'ALL'
    );

    CREATE TYPE room_role_type AS ENUM (
      'VISITOR',
      'INVITED',
      'DEFAULT',
      'MEMBER',
      'MASTER'
    );

    CREATE TYPE room_member_type AS ENUM (
      'INVITED',
      'MEMBER',
      'MASTER'
    );

    CREATE TYPE room_message_type AS ENUM (
      'NORMAL',
      'IMAGE'
    );

    CREATE TYPE board AS ENUM (
      'COMMUNITY',
      'TRADE',
      'BUG'
    );

    CREATE TYPE thread_status AS ENUM (
      'OPEN',
      'CLOSED',
      'DELETED'
    );

    CREATE TYPE notice_type AS ENUM (
      'UPDATE',
      'NOTICE'
    );

    CREATE TYPE forum_post_type AS ENUM (
      'SIGNED_IN',
      'ANONYMOUS',
      'ADMINISTRATOR'
    );

    CREATE TYPE forum_topic_status AS ENUM (
      'OPEN',
      'CLOSE'
    );

    CREATE TYPE message_fetch_config_category AS ENUM (
      'all',
      'follow',
      'follow-other',
      'replied',
      'replied-other',
      'own',
      'conversation',
      'search',
      'list',
      'character',
      'character-replied'
    );

    CREATE TABLE game_status (
			nth INT NOT NULL DEFAULT 0
		);

    CREATE TABLE announcements (
      id           SERIAL    NOT NULL PRIMARY KEY,
      title        TEXT      NOT NULL CHECK (title    != ''),
      overview     TEXT      NOT NULL CHECK (overview != ''),
      content      TEXT      NOT NULL CHECK (content  != ''),
      announced_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
      updated_at   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
      type         announcement_type NOT NULL
    );
    CREATE INDEX ON announcements(announced_at);
    CREATE INDEX ON announcements(updated_at);

    CREATE TABLE characters (
      seq                          SERIAL    NOT NULL PRIMARY KEY,
      id                           INT                UNIQUE,
      administrator                BOOLEAN   NOT NULL DEFAULT false,
      password                     TEXT      NOT NULL CHECK (password != ''),
      name                         TEXT      NOT NULL CHECK (name     != ''),
      nickname                     TEXT      NOT NULL CHECK (nickname != ''),
      username                     TEXT      NOT NULL CHECK (username != ''),
      email                        TEXT,
      summary                      TEXT      NOT NULL DEFAULT '',
      profile                      TEXT      NOT NULL DEFAULT '',
      mainicon                     TEXT      NOT NULL DEFAULT '',
			ap_recover_at                TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
      diary                        TEXT,
      diary_title                  TEXT,
      deleted_at                   TIMESTAMP          DEFAULT NULL,
      banned_at                    TIMESTAMP          DEFAULT NULL,
      registered_at                TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
      webhook                      TEXT      NOT NULL DEFAULT '',
      webhook_followed             BOOLEAN   NOT NULL DEFAULT true,
      webhook_replied              BOOLEAN   NOT NULL DEFAULT true,
      webhook_subscribe            BOOLEAN   NOT NULL DEFAULT true,
      webhook_mail                 BOOLEAN   NOT NULL DEFAULT true,
      notification_followed        BOOLEAN   NOT NULL DEFAULT true,
      notification_replied         BOOLEAN   NOT NULL DEFAULT true,
      notification_subscribe       BOOLEAN   NOT NULL DEFAULT true,
      notification_mail            BOOLEAN   NOT NULL DEFAULT true,
      notification_last_checked_at TIMESTAMP NOT NULL DEFAULT '2000-01-01 00:00:00-00',
      notification_token           TEXT      NOT NULL,
      book_permission_relates      BOOLEAN   NOT NULL DEFAULT true,
      book_permission              relation_permission NOT NULL DEFAULT 'ALL'
    );
    CREATE INDEX ON characters(id);
    CREATE INDEX ON characters(deleted_at);
    CREATE INDEX ON characters(banned_at);
    CREATE INDEX ON characters(administrator);
    CREATE INDEX ON characters USING gin(name gin_bigm_ops);
    CREATE INDEX ON characters USING gin(nickname gin_bigm_ops);
    CREATE UNIQUE INDEX ON characters(username);
    CREATE UNIQUE INDEX ON characters(email);
    CREATE UNIQUE INDEX ON characters(notification_token);

    CREATE TABLE characters_tags (
      id        SERIAL NOT NULL PRIMARY KEY,
      character INT    NOT NULL REFERENCES characters(id),
      tag       TEXT   NOT NULL CHECK (tag != '')
    );
    CREATE INDEX ON characters_tags(character);
    CREATE INDEX ON characters_tags USING gin(tag gin_bigm_ops);
  
    CREATE TABLE characters_profile_images (
      id        SERIAL NOT NULL PRIMARY KEY,
      character INT    NOT NULL REFERENCES characters(id),
      path      TEXT   NOT NULL CHECK (path != '')
    );
    CREATE INDEX ON characters_profile_images(character);

    CREATE TABLE characters_icons (
      id        SERIAL NOT NULL PRIMARY KEY,
      character INT    NOT NULL REFERENCES characters(id),
      path      TEXT   NOT NULL CHECK (path != '')
    );
    CREATE INDEX ON characters_icons(character);

    CREATE TABLE characters_icon_layering_groups (
      id        SERIAL NOT NULL PRIMARY KEY,
      character INT    NOT NULL REFERENCES characters(id),
      name      TEXT   NOT NULL CHECK (name != '')
    );
    CREATE INDEX ON characters_icon_layering_groups(character);

    CREATE TABLE characters_icon_process_schemas (
      id             SERIAL           NOT NULL PRIMARY KEY,
      layering_group INT              NOT NULL REFERENCES characters_icon_layering_groups(id) ON DELETE CASCADE ON UPDATE CASCADE,
      name           TEXT             NOT NULL CHECK (name != ''),
      x              DOUBLE PRECISION NOT NULL,
      y              DOUBLE PRECISION NOT NULL,
      scale          DOUBLE PRECISION NOT NULL,
      rotate         DOUBLE PRECISION NOT NULL
    );
    CREATE INDEX ON characters_icon_process_schemas(layering_group);

    CREATE TABLE characters_icon_layer_groups (
      id             SERIAL NOT NULL PRIMARY KEY,
      layering_group INT    NOT NULL REFERENCES characters_icon_layering_groups(id) ON DELETE CASCADE ON UPDATE CASCADE,
      name           TEXT   NOT NULL CHECK (name != ''),
      layer_order    INT    NOT NULL
    );
    CREATE INDEX ON characters_icon_layer_groups(layering_group);

    CREATE TABLE characters_icon_layer_items (
      id          SERIAL NOT NULL PRIMARY KEY,
      layer_group INT    NOT NULL REFERENCES characters_icon_layer_groups(id) ON DELETE CASCADE ON UPDATE CASCADE,
      path        TEXT   NOT NULL CHECK (path != '')
    );
    CREATE INDEX ON characters_icon_layer_items(layer_group);

    CREATE TABLE characters_uploaded_images (
      id          SERIAL     NOT NULL PRIMARY KEY,
      character   INT        NOT NULL REFERENCES characters(id),
      path        TEXT       NOT NULL CHECK (path != ''),
      md5         TEXT       NOT NULL CHECK (md5  != ''),
      uploaded_at TIMESTAMP  NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX ON characters_uploaded_images(character);
    CREATE UNIQUE INDEX ON characters_uploaded_images(character, md5);
  
    CREATE TABLE characters_twitter (
      id         SERIAL NOT NULL PRIMARY KEY,
      character  INT    NOT NULL REFERENCES characters(id),
      twitter_id TEXT   NOT NULL
    );
    CREATE UNIQUE INDEX ON characters_twitter(character);
    CREATE UNIQUE INDEX ON characters_twitter(twitter_id);
  
    CREATE TABLE characters_google (
      id         SERIAL NOT NULL PRIMARY KEY,
      character  INT    NOT NULL REFERENCES characters(id),
      google_id  TEXT   NOT NULL
    );
    CREATE UNIQUE INDEX ON characters_google(character);
    CREATE UNIQUE INDEX ON characters_google(google_id);
  
    CREATE TABLE mail_confirm_codes (
      id        SERIAL    NOT NULL PRIMARY KEY,
      character INT       NOT NULL REFERENCES characters(id),
      email     TEXT      NOT NULL CHECK (email != ''),
      code      TEXT      NOT NULL CHECK (code != ''),
      expire    TIMESTAMP NOT NULL
    );
    CREATE UNIQUE INDEX ON mail_confirm_codes(code);
  
    CREATE TABLE password_reset_codes (
      id        SERIAL    NOT NULL PRIMARY KEY,
      character INT       NOT NULL REFERENCES characters(id),
      code      TEXT      NOT NULL CHECK (code != ''),
      expire    TIMESTAMP NOT NULL
    );
    CREATE UNIQUE INDEX ON password_reset_codes(code);
  
		CREATE TABLE follows (
			id          SERIAL    NOT NULL PRIMARY KEY,
			follower    INT       NOT NULL REFERENCES characters(id),
			followed    INT       NOT NULL REFERENCES characters(id),
      followed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			CHECK (follower != followed)
		);
		CREATE INDEX ON follows(follower);
		CREATE INDEX ON follows(followed);
    CREATE INDEX ON follows(followed_at);
  
		CREATE TABLE mutes (
			id       SERIAL    NOT NULL PRIMARY KEY,
			muter    INT       NOT NULL REFERENCES characters(id),
			muted    INT       NOT NULL REFERENCES characters(id),
      muted_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			CHECK (muter != muted)
		);
		CREATE INDEX ON mutes(muter);
		CREATE INDEX ON mutes(muted);
    CREATE INDEX ON mutes(muted_at);

		CREATE TABLE blocks (
			id         SERIAL    NOT NULL PRIMARY KEY,
			blocker    INT       NOT NULL REFERENCES characters(id),
			blocked    INT       NOT NULL REFERENCES characters(id),
      blocked_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			CHECK (blocker != blocked)
		);
		CREATE INDEX ON blocks(blocker);
		CREATE INDEX ON blocks(blocked);
    CREATE INDEX ON blocks(blocked_at);

    CREATE TABLE character_prohibition_related_data (
			id        SERIAL    NOT NULL PRIMARY KEY,
      character INT       NOT NULL REFERENCES characters(id),
      reason    TEXT      NOT NULL,
      timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
      type      prohibition_related_data_type NOT NULL
    );
		CREATE INDEX ON character_prohibition_related_data(character);
		CREATE INDEX ON character_prohibition_related_data(timestamp);

    CREATE TABLE rooms (
      seq                  SERIAL    NOT NULL PRIMARY KEY,
      id                   INT                UNIQUE,
      master               INT       NOT NULL REFERENCES characters(id),
      official             BOOLEAN   NOT NULL DEFAULT false,
      belong               INT                REFERENCES rooms(id),
      title                TEXT      NOT NULL CHECK (title != ''),
      summary              TEXT      NOT NULL,
      description          TEXT      NOT NULL,
      searchable           BOOLEAN   NOT NULL,
      allow_recommendation BOOLEAN   NOT NULL,
      children_referable   BOOLEAN   NOT NULL,
      messages_count       INT       NOT NULL DEFAULT 0,
      members_count        INT       NOT NULL DEFAULT 1,
      created_at           TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
      updated_at           TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
      deleted_at           TIMESTAMP
    );
    CREATE UNIQUE INDEX ON rooms(id);
    CREATE INDEX ON rooms(master);
    CREATE INDEX ON rooms(official);
    CREATE INDEX ON rooms(belong);
    CREATE INDEX ON rooms(searchable);
    CREATE INDEX ON rooms(allow_recommendation);
    CREATE INDEX ON rooms USING gin(title gin_bigm_ops);
    
    CREATE TABLE rooms_tags (
      id   SERIAL NOT NULL PRIMARY KEY,
      room INT    NOT NULL REFERENCES rooms(id),
      tag  TEXT   NOT NULL CHECK (tag != '')
    );
    CREATE INDEX ON rooms_tags(room);
    CREATE INDEX ON rooms_tags USING gin(tag gin_bigm_ops);
    
    CREATE TABLE rooms_members (
      id                   SERIAL  NOT NULL PRIMARY KEY,
      room                 INT     NOT NULL REFERENCES rooms(id),
      member               INT     NOT NULL REFERENCES characters(id),
      write                BOOLEAN NOT NULL,
      ban                  BOOLEAN NOT NULL,
      invite               BOOLEAN NOT NULL,
      use_reply            BOOLEAN NOT NULL,
      use_secret           BOOLEAN NOT NULL,
      delete_other_message BOOLEAN NOT NULL,
      create_children_room BOOLEAN NOT NULL,
      type                 room_member_type NOT NULL
    );
    CREATE UNIQUE INDEX ON rooms_members(room, member);
    CREATE INDEX ON rooms_members(room);
    CREATE INDEX ON rooms_members(member);

    CREATE TABLE rooms_invited_characters (
      id         SERIAL    NOT NULL PRIMARY KEY,
      room       INT       NOT NULL REFERENCES rooms(id),
      invited    INT       NOT NULL REFERENCES characters(id),
      inviter    INT       NOT NULL REFERENCES characters(id),
      invited_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE UNIQUE INDEX ON rooms_invited_characters(room, invited);
    CREATE INDEX ON rooms_invited_characters(room);
    CREATE INDEX ON rooms_invited_characters(invited);
    CREATE INDEX ON rooms_invited_characters(invited_at);
  
    CREATE TABLE rooms_banned_characters (
      id        SERIAL    NOT NULL PRIMARY KEY,
      room      INT       NOT NULL REFERENCES rooms(id),
      banned    INT       NOT NULL REFERENCES characters(id),
      banner    INT       NOT NULL REFERENCES characters(id),
      banned_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE UNIQUE INDEX ON rooms_banned_characters(room, banned);
    CREATE INDEX ON rooms_banned_characters(room);
    CREATE INDEX ON rooms_banned_characters(banned);
    CREATE INDEX ON rooms_banned_characters(banned_at);
  
    CREATE TABLE rooms_subscribers (
      id        SERIAL NOT NULL PRIMARY KEY,
      room      INT    NOT NULL REFERENCES rooms(id),
      character INT    NOT NULL REFERENCES characters(id)
    );
    CREATE UNIQUE INDEX ON rooms_subscribers(room, character);
    CREATE INDEX ON rooms_subscribers(room);
    CREATE INDEX ON rooms_subscribers(character);
  
    CREATE TABLE rooms_role_tables_advisory_locker (
      id           SERIAL  NOT NULL PRIMARY KEY,
      room         INT     NOT NULL REFERENCES rooms(id),
      role_version INT NOT NULL DEFAULT 0
    );
    CREATE UNIQUE INDEX ON rooms_role_tables_advisory_locker(role_version);
  
    CREATE TABLE rooms_roles (
      id                   SERIAL  NOT NULL PRIMARY KEY,
      room                 INT     NOT NULL REFERENCES rooms(id),
      priority             INT     NOT NULL,
      name                 TEXT,
      read                 BOOLEAN,
      write                BOOLEAN,
      ban                  BOOLEAN,
      invite               BOOLEAN,
      use_reply            BOOLEAN,
      use_secret           BOOLEAN,
      delete_other_message BOOLEAN,
      create_children_room BOOLEAN,
      type                 room_role_type DEFAULT 'MEMBER'
    );
    CREATE UNIQUE INDEX ON rooms_roles(room, name);
    CREATE UNIQUE INDEX ON rooms_roles(room, priority);
    CREATE INDEX ON rooms_roles(room);

    CREATE TABLE rooms_roles_members (
      id        SERIAL NOT NULL PRIMARY KEY,
      role      INT    NOT NULL REFERENCES rooms_roles(id),
      character INT    NOT NULL REFERENCES characters(id)
    );
    CREATE UNIQUE INDEX ON rooms_roles_members(role, character);
    CREATE INDEX ON rooms_roles_members(role);
    CREATE INDEX ON rooms_roles_members(character);

    CREATE TABLE rooms_messages (
      id            SERIAL    NOT NULL PRIMARY KEY,
      room          INT       NOT NULL REFERENCES rooms(id),
      character     INT       NOT NULL REFERENCES characters(id),
      refer         INT                REFERENCES rooms_messages(id),
      refer_root    INT                REFERENCES rooms_messages(id),
      secret        BOOLEAN   NOT NULL DEFAULT false,
      icon          TEXT,
      name          TEXT      NOT NULL,
      message       TEXT      NOT NULL,
      search_text   TEXT      NOT NULL,
      replied_count INT       NOT NULL DEFAULT 0,
      deleted_at    TIMESTAMP          DEFAULT NULL,
      posted_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
      single        BOOLEAN   NOT NULL,
      relates       INT[]     NOT NULL,
      reply_permission relation_permission NOT NULL DEFAULT 'ALL'
    );
    CREATE INDEX ON rooms_messages(room);
    CREATE INDEX ON rooms_messages(character);
    CREATE INDEX ON rooms_messages(single);
    CREATE INDEX ON rooms_messages USING gin(search_text gin_bigm_ops);
    CREATE INDEX ON rooms_messages USING gin(relates);

    CREATE TABLE rooms_messages_belongs (
      id      SERIAL NOT NULL PRIMARY KEY,
      message INT    NOT NULL REFERENCES rooms_messages(id),
      room    INT    NOT NULL REFERENCES rooms(id)
    );
    CREATE INDEX ON rooms_messages_belongs(message);
    CREATE INDEX ON rooms_messages_belongs(room);

    CREATE TABLE rooms_messages_recipients (
      id        SERIAL NOT NULL PRIMARY KEY,
      message   INT    NOT NULL REFERENCES rooms_messages(id),
      character INT    NOT NULL REFERENCES characters(id)
    );
    CREATE UNIQUE INDEX ON rooms_messages_recipients(message, character);

    CREATE TABLE message_fetch_configs (
      id            SERIAL NOT NULL PRIMARY KEY,
      master        INT    NOT NULL REFERENCES characters(id),
      config_order  INT    NOT NULL,
      name          TEXT   NOT NULL,
      room          INT,
      search        TEXT,
      refer_root    INT,
      list          INT,
      character     INT,
      relate_filter BOOLEAN,
      children      BOOLEAN,
      category      message_fetch_config_category NOT NULL
    );
    CREATE INDEX ON message_fetch_configs(master);
    CREATE INDEX ON message_fetch_configs(config_order);

    CREATE TABLE notifications (
      id             SERIAL    NOT NULL PRIMARY KEY,
      character      INT       NOT NULL REFERENCES characters(id),
      icon           TEXT,
      message        TEXT      NOT NULL CHECK (message != ''),
      detail         TEXT      NOT NULL DEFAULT '',
      value          TEXT      NOT NULL DEFAULT '',
      notificated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
      type           notification_type NOT NULL
    );
    CREATE INDEX ON notifications(character);
    CREATE INDEX ON notifications(notificated_at);

    CREATE TABLE mails (
      id           SERIAL    NOT NULL PRIMARY KEY,
      sender       INT                REFERENCES characters(id),
      receiver     INT       NOT NULL REFERENCES characters(id),
      name         TEXT      NOT NULL DEFAULT '',
      title        TEXT      NOT NULL DEFAULT '',
      message      TEXT      NOT NULL CHECK (message != ''),
      read         BOOLEAN   NOT NULL DEFAULT false,
      posted_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
      deleted_at   TIMESTAMP          DEFAULT NULL,
      CHECK (sender != receiver)
    );
    CREATE INDEX ON mails(sender);
    CREATE INDEX ON mails(receiver);
    CREATE INDEX ON mails(posted_at);

    CREATE TABLE mails_everyone (
      id        SERIAL    NOT NULL PRIMARY KEY,
      name      TEXT      NOT NULL CHECK (name != ''),
      title     TEXT      NOT NULL CHECK (title != ''),
      message   TEXT      NOT NULL CHECK (message != ''),
      posted_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX ON mails_everyone(posted_at);

    CREATE TABLE threads (
      id             SERIAL        NOT NULL PRIMARY KEY,
      title          TEXT          NOT NULL CHECK (title      != ''),
      name           TEXT          NOT NULL CHECK (name       != ''),
      identifier     TEXT          NOT NULL CHECK (identifier != ''),
      password       TEXT          NOT NULL CHECK (password   != ''),
      created_at     TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
      updated_at     TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
      last_posted_at TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
      administrator  BOOLEAN       NOT NULL DEFAULT false,
      board          board         NOT NULL,
      status         thread_status NOT NULL DEFAULT 'OPEN'
    );
    CREATE INDEX ON threads(board);
    CREATE INDEX ON threads(status);

    CREATE TABLE threads_responses (
      seq           SERIAL    NOT NULL PRIMARY KEY,
      id            INT,
      thread        INT       NOT NULL REFERENCES threads(id),
      name          TEXT      NOT NULL CHECK (name       != ''),
      identifier    TEXT      NOT NULL CHECK (identifier != ''),
      message       TEXT      NOT NULL CHECK (message    != ''),
      secret        TEXT      NOT NULL,
      password      TEXT      NOT NULL CHECK (password != ''),
      administrator BOOLEAN   NOT NULL DEFAULT false,
      deleted_at    TIMESTAMP           DEFAULT NULL,
      posted_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX ON threads_responses(thread);
    CREATE INDEX ON threads_responses(id);

    CREATE TABLE diaries (
      id          SERIAL    NOT NULL PRIMARY KEY,
      character   INT       NOT NULL REFERENCES characters(id),
      title       TEXT      NOT NULL,
      diary       TEXT      NOT NULL,
      nth         INT       NOT NULL
    );
    CREATE INDEX ON diaries(character);
    CREATE INDEX ON diaries(nth);
    CREATE UNIQUE INDEX ON diaries(character, nth);

		CREATE TABLE lists (
			id         SERIAL    NOT NULL PRIMARY KEY,
			master     INT       NOT NULL REFERENCES characters(id),
			name       TEXT      NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX ON lists(master);

		CREATE TABLE lists_characters (
			id        SERIAL NOT NULL PRIMARY KEY,
			list      INT    NOT NULL REFERENCES lists(id),
			character INT    NOT NULL REFERENCES characters(id)
		);
		CREATE INDEX ON lists_characters(list);
		CREATE INDEX ON lists_characters(character);
		CREATE UNIQUE INDEX ON lists_characters(list, character);

		CREATE TABLE inquiries (
			id          SERIAL    NOT NULL PRIMARY KEY,
			character   INT                REFERENCES characters(id),
			inquiry     TEXT      NOT NULL CHECK (inquiry != ''),
			posted_at   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
      resolved_at TIMESTAMP 
		);
		CREATE INDEX ON inquiries(posted_at);
    CREATE INDEX ON inquiries(resolved_at);

    CREATE TABLE notices (
      id         SERIAL    NOT NULL PRIMARY KEY,
      title      TEXT      NOT NULL CHECK (title   != ''),
      content    TEXT      NOT NULL CHECK (content != ''),
      noticed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
      updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
      type       notice_type NOT NULL
    );
    CREATE INDEX ON notices(noticed_at);

    CREATE TABLE forum_groups (
      id                 SERIAL NOT NULL PRIMARY KEY,
      title              TEXT   NOT NULL CHECK (title != ''),
      forum_groups_order INT    NOT NULL
    );
    
    CREATE TABLE forums (
      id              SERIAL NOT NULL PRIMARY KEY,
      forum_group     INT    NOT NULL REFERENCES forum_groups(id),
      title           TEXT   NOT NULL CHECK (title != ''),
      summary         TEXT   NOT NULL,
      guide           TEXT   NOT NULL,
      forum_order     INT    NOT NULL,
      force_post_type forum_post_type
    );
    CREATE INDEX ON forums(forum_group);
    CREATE INDEX ON forums(forum_order);

    CREATE TABLE forum_topics (
      id             SERIAL    NOT NULL PRIMARY KEY,
      forum          INT       NOT NULL REFERENCES forums(id),
      title          TEXT      NOT NULL CHECK (title   != ''),
      character      INT                REFERENCES characters(id),
      name           TEXT,
      edit_password  TEXT,
      identifier     TEXT,
      posts          INT       NOT NULL DEFAULT 1,
      last_posted_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
      created_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
      pinned_order   INT,
      status         forum_topic_status NOT NULL DEFAULT 'OPEN',
      post_type      forum_post_type    NOT NULL,
      CHECK (
        (
          post_type = 'ADMINISTRATOR' AND
          name          IS NULL AND
          edit_password IS NULL AND
          identifier    IS NULL
        ) OR
        (
          post_type = 'SIGNED_IN' AND
          name          IS NULL AND
          edit_password IS NULL AND
          identifier    IS NULL
        ) OR
        (
          post_type = 'ANONYMOUS' AND
          name          IS NOT NULL AND
          edit_password IS NOT NULL AND
          identifier    IS NOT NULL
        )
      )
    );
    CREATE INDEX ON forum_topics(forum);
    CREATE INDEX ON forum_topics(last_posted_at);
    CREATE INDEX ON forum_topics(pinned_order);

    CREATE TABLE forum_topics_posts (
      id            SERIAL    NOT NULL PRIMARY KEY,
      topic         INT       NOT NULL REFERENCES forum_topics(id),
      character     INT                REFERENCES characters(id),
      name          TEXT,
      edit_password TEXT,
      identifier    TEXT,
      content       TEXT      NOT NULL CHECK (content != ''),
      posted_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
      updated_at    TIMESTAMP,
      deleted_at    TIMESTAMP,
      post_type     forum_post_type NOT NULL,
      CHECK (
        (
          post_type = 'ADMINISTRATOR' AND
          name          IS NULL AND
          edit_password IS NULL AND
          identifier    IS NULL
        ) OR
        (
          post_type = 'SIGNED_IN' AND
          name          IS NULL AND
          edit_password IS NULL AND
          identifier    IS NULL
        ) OR
        (
          post_type = 'ANONYMOUS' AND
          name          IS NOT NULL AND
          edit_password IS NOT NULL AND
          identifier    IS NOT NULL
        )
      )
    );
    CREATE INDEX ON forum_topics_posts(topic);
    CREATE INDEX ON forum_topics_posts(posted_at);
    CREATE INDEX ON forum_topics_posts(deleted_at);

    CREATE TABLE forum_topics_posts_reactions (
      id         SERIAL    NOT NULL PRIMARY KEY,
      post       INT       NOT NULL REFERENCES forum_topics_posts(id),
      emoji      TEXT      NOT NULL,
      character  INT       NOT NULL REFERENCES characters(id),
      reacted_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
    CREATE INDEX ON forum_topics_posts_reactions(post);
    CREATE INDEX ON forum_topics_posts_reactions(emoji);
    CREATE INDEX ON forum_topics_posts_reactions(character);
    CREATE INDEX ON forum_topics_posts_reactions(reacted_at);
    CREATE UNIQUE INDEX ON forum_topics_posts_reactions(post, emoji, character);

    CREATE TABLE forum_topics_posts_revisions (
      id        SERIAL    NOT NULL PRIMARY KEY,
      post      INT       NOT NULL REFERENCES characters(id),
      content   TEXT      NOT NULL,
      posted_at TIMESTAMP NOT NULL
    );
    CREATE INDEX ON forum_topics_posts_revisions(post);
    CREATE INDEX ON forum_topics_posts_revisions(posted_at);

		INSERT INTO game_status (
			nth
		) VALUES (
			0
		);
  `)

	return nil
}
