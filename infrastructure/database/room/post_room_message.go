package room

import (
	"errors"
	"golang.org/x/exp/slices"

	"github.com/jmoiron/sqlx"
	"github.com/kaikourok/lunchtote-backend/entity/model"
	"github.com/kaikourok/lunchtote-backend/entity/service"
	"github.com/lib/pq"
)

func (db *RoomRepository) PostRoomMessage(characterId int, message *model.RoomPostMessage, uploadPath string) (messageId int, err error) {
	err = db.ExecTx(func(tx *sqlx.Tx) error {
		relates := make([]int32, 0, 16)
		var referRoot *int
		var userName string
		if message.Refer != nil {
			row := tx.QueryRowx(`
				WITH message AS (
					SELECT
						refer_root,
						character,
						reply_permission,
						room,
						relates
					FROM
						rooms_messages
					WHERE
						id = $2
				),
				invalid_list AS (
					SELECT blocker FROM blocks WHERE blocked = $1 UNION ALL
					SELECT blocked FROM blocks WHERE blocker = $1
				)

				SELECT
					(SELECT refer_root FROM message),
					EXISTS (SELECT * FROM follows WHERE follower = $1 AND followed = (SELECT character FROM message)),
					EXISTS (SELECT * FROM follows WHERE followed = $1 AND follower = (SELECT character FROM message)),
					EXISTS (SELECT * FROM rooms_messages_recipients WHERE message = $2 AND character IN (SELECT * FROM invalid_list)),
					(SELECT reply_permission FROM message),
					(SELECT nickname FROM characters WHERE id = $1),
					(SELECT character FROM message),
					(SELECT relates FROM message);
			`, characterId, message.Refer)

			var isFollowing, isFollower, isInvalid bool
			var replyPermission string
			var targetCharacterId int

			err := row.Scan(
				&referRoot,
				&isFollowing,
				&isFollower,
				&isInvalid,
				&replyPermission,
				&userName,
				&targetCharacterId,
				pq.Array(&relates),
			)
			if err != nil {
				return err
			}

			if isInvalid {
				return errors.New("返信を送る権限がありません")
			}

			var replyable bool
			switch replyPermission {
			case "DISALLOW":
				replyable = false
			case "FOLLOW":
				replyable = isFollowing
			case "FOLLOWED":
				replyable = isFollower
			case "MUTUAL_FOLLOW":
				replyable = isFollowing && isFollower
			case "REPLY":
				replyable = slices.Contains(relates, int32(characterId))
			case "ALL":
				replyable = true
			default:
				return errors.New("対象のメッセージがありません")
			}

			if !replyable && targetCharacterId != characterId {
				tx.Rollback()
				return errors.New("返信を送る権限がありません")
			}

		} else if message.DirectReply != nil {
			row := tx.QueryRowx(`
				SELECT
					EXISTS (SELECT * FROM characters_blocks WHERE blocker  = $1 AND blocked  = $2) AS block,
					EXISTS (SELECT * FROM characters_blocks WHERE blocked  = $1 AND blocker  = $2) AS blocked,
					EXISTS (SELECT * FROM characters        WHERE character_id = $2 AND deleted_at IS NULL) AS targetExists;
			`, characterId, message.DirectReply)

			var isBlocking, isBlocked, isTargetExists bool
			err := row.Scan(&isBlocking, &isBlocked, &isTargetExists)
			if err != nil {
				tx.Rollback()
				return err
			}

			if isBlocking || isBlocked {
				tx.Rollback()
				return errors.New("返信を送る権限がありません")
			}

			if !isTargetExists {
				tx.Rollback()
				return errors.New("返信先が存在しません")
			}

			relates = append(relates, int32(*message.DirectReply))
		}

		found := false
		for i := range relates {
			if relates[i] == int32(characterId) {
				found = true
				break
			}
		}
		if !found {
			relates = append(relates, int32(characterId))
		}

		row := tx.QueryRowx(`
			INSERT INTO rooms_messages (
				character,
				room,
				refer,
				refer_root,
				icon,
				name,
				message,
				search_text,
				reply_permission,
				secret,
				single,
				relates
			)
			SELECT
				$1 AS character,
				$2 AS room,
				$3 AS refer,
				$4 AS refer_root,
				$5 AS icon,
				(
					CASE $6
						WHEN '' THEN (SELECT nickname FROM characters WHERE id = $1)
						ELSE $6
					END
				) AS name,
				$7 AS message,
				$8 AS search_text,
				$9 AS reply_permission,
				$10 AS secret,
				$11 AS single,
				$12 AS relates
			RETURNING
				id;
		`,
			characterId,
			message.Room,
			message.Refer,
			referRoot,
			message.Icon,
			message.Name,
			service.StylizeMessage(message.Message, uploadPath),
			service.ConvertMessageToSearchText(message.Message),
			message.ReplyPermission,
			message.Secret,
			len(relates) == 1,
			pq.Array(relates),
		)

		err := row.Scan(&messageId)
		if err != nil {
			return err
		}

		if !message.Secret {
			_, err = tx.Exec(`
				UPDATE
					rooms
				SET
					messages_count = messages_count + 1
				WHERE
					id = $1;
			`, message.Room)
			if err != nil {
				return err
			}
		}

		_, err = tx.Exec(`
			WITH RECURSIVE ancestor_rooms (id, belong) AS (
				SELECT rooms.id, rooms.belong FROM rooms WHERE id = $1
				UNION ALL
				SELECT rooms.id, rooms.belong FROM ancestor_rooms, rooms WHERE ancestor_rooms.belong = rooms.id 
			)

			INSERT INTO rooms_messages_belongs (
				room,
				message
			) 
			SELECT
				id,
				$2
			FROM
				ancestor_rooms;
		`, message.Room, messageId)
		if err != nil {
			return err
		}

		if message.Refer == nil {
			_, err = tx.Exec(`
				UPDATE
					rooms_messages
				SET
					refer_root = $1
				WHERE
					id = $1;
			`, messageId)
			if err != nil {
				return err
			}
		}

		for i := range relates {
			_, err = tx.Exec(`
				INSERT INTO rooms_messages_recipients (
					message,
					character
				) VALUES (
					$1,
					$2
				);
			`, messageId, relates[i])
			if err != nil {
				return err
			}
		}

		if !message.Secret && message.Refer != nil {
			_, err = tx.Exec(`
				UPDATE
					rooms_messages
				SET
					replied_count = replied_count + 1
				WHERE
					id = $1;
			`, message.Refer)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return
}
