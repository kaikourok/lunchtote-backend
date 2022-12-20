package room

import (
	"errors"
	"sort"

	"github.com/jmoiron/sqlx"
	"github.com/kaikourok/lunchtote-backend/entity/model"
	"github.com/kaikourok/lunchtote-backend/library/slice"
)

func (db *RoomRepository) UpdateRolePriorities(characterId, roomId int, priorities *[]model.RoomRolePriority) error {
	return db.ExecTx(func(tx *sqlx.Tx) error {
		err := db.roleTablesAdvisoryLock(roomId, tx)
		if err != nil {
			return err
		}

		rows, err := tx.Queryx(`
			SELECT
				rooms.id,
				rooms.master,
				rooms_roles.id,
				rooms_roles.priority
			FROM
				rooms
			JOIN
				rooms_roles ON (rooms.id = rooms_roles.room AND rooms.id = $1)
			WHERE
				rooms_roles.type = 'MEMBER';
		`, roomId)

		if err != nil {
			return err
		}
		defer rows.Close()

		roleIds := make([]int, 0, len(*priorities))
		maxPriority := 0

		for rows.Next() {
			var roleRoomId, masterId, roleId, prioriry int
			err = rows.Scan(
				&roleRoomId,
				&masterId,
				&roleId,
				&prioriry,
			)

			if characterId != masterId {
				return errors.New("ルームのロールを変更する権限がありません")
			}

			if maxPriority < prioriry {
				maxPriority = prioriry
			}

			if err != nil {
				return err
			}

			roleIds = append(roleIds, roleId)
		}

		if len(roleIds) != len(*priorities) {
			return errors.New("指定数が不足しています")
		}

		for _, priority := range *priorities {
			if !slice.Contains(priority.Role, &roleIds) {
				return errors.New("ルーム外のIDが指定されています")
			}
		}

		priorities := *priorities
		sort.Slice(priorities, func(a, b int) bool {
			return priorities[a].Priority > priorities[b].Priority
		})

		for i := range priorities {
			priorities[i].Priority = i + maxPriority + 1
		}

		for _, priority := range priorities {
			_, err = tx.NamedExec(`
				UPDATE
					rooms_roles 
				SET
					priority = :priority
				WHERE
					id = :role;
			`, priority)
			if err != nil {
				return err
			}
		}

		if err != nil {
			return err
		}

		return db.attachPermissions(roomId, tx)
	})
}
