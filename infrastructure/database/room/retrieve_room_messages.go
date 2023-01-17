package room

import (
	"encoding/json"
	"errors"
	"log"
	"math"

	"github.com/kaikourok/lunchtote-backend/entity/model"
	"github.com/kaikourok/lunchtote-backend/library/slice"
)

func (db *RoomRepository) RetrieveRoomMessages(characterId int, options *model.RoomMessageRetrieveSettings) (messages *[]model.RoomMessage, isContinueFollowing, isContinuePrevious *bool, err error) {
	basePoint := options.BasePoint
	if options.RangeType == "latest" {
		basePoint = math.MaxInt32
	}

	type RetrieveMessagesResult struct {
		messages   *[]model.RoomMessage
		isContinue bool
		err        error
	}
	var previousResultChannel chan RetrieveMessagesResult
	var followingResultChannel chan RetrieveMessagesResult

	generateSql := func(fetchMode string) (string, *map[string]interface{}) {
		args := make(map[string]interface{}, 16)
		args["user_id"] = characterId
		args["base_point"] = basePoint
		args["fetch_number"] = options.FetchNumber + 1

		if options.Room != nil {
			args["room"] = options.Room
		}
		if slice.Contains(options.Category, &[]string{"character-replied", "character"}) {
			args["target_id"] = options.TargetCharacterId
		}
		if options.Category == "conversation" {
			args["refer_root"] = options.ReferRoot
		}
		if options.Category == "search" {
			args["search"] = options.Search
		}
		if options.Category == "list" {
			args["list"] = options.ListId
		}

		sql := `
			WITH
				follow_list   AS (SELECT followed AS character FROM follows WHERE follower = :user_id),
				follower_list AS (SELECT follower AS character FROM follows WHERE followed = :user_id),
				mute_list     AS (SELECT muted    AS character FROM mutes   WHERE muter    = :user_id),
				block_list    AS (SELECT blocked  AS character FROM blocks  WHERE blocker  = :user_id UNION ALL
				                  SELECT blocker  AS character FROM blocks  WHERE blocked  = :user_id),
				banned_rooms  AS (SELECT room FROM rooms_banned_characters WHERE banned = :user_id)
		`
		if options.Category == "list" {
			sql += `
				, list_members AS (SELECT lists_characters.character FROM lists JOIN lists_characters ON (lists.id = lists_characters.list) WHERE lists.id = :list AND lists.master = :user_id)
			`
		}

		sql += `
			SELECT
				messages.id,
				messages.character,
				messages.refer,
				messages.refer_root,
				messages.secret,
				messages.icon,
				messages.name,
				messages.message,
				messages.posted_at,
				messages.reply_permission,
				messages.replied_count,
				rooms.id,
				rooms.title,
				(
          CASE messages.reply_permission
						WHEN 'DISALLOW'      THEN messages.character = :user_id
						WHEN 'FOLLOW'        THEN messages.character = :user_id OR  EXISTS (SELECT * FROM follower_list WHERE character = messages.character)
						WHEN 'FOLLOWED'      THEN messages.character = :user_id OR  EXISTS (SELECT * FROM follow_list   WHERE character = messages.character)
						WHEN 'MUTUAL_FOLLOW' THEN messages.character = :user_id OR (EXISTS (SELECT * FROM follower_list WHERE character = messages.character) AND EXISTS (SELECT * FROM follow_list WHERE character = messages.character))
						WHEN 'ALL'           THEN true
          END
				),
				JSON_AGG(JSON_BUILD_OBJECT(
					'id',   recipient.id,
					'name', recipient.nickname	
				))
			FROM
				rooms_messages AS messages
			JOIN
				rooms_messages_recipients AS messages_recipients ON (messages.id = messages_recipients.message)
			JOIN
				characters AS recipient ON (messages_recipients.character = recipient.id)
		`

		if options.Room == nil {
			sql += `
				JOIN
					rooms ON (messages.room = rooms.id)
			`
		} else {
			if options.Children {
				sql += `
					JOIN
						rooms_messages_belongs ON (messages.id = rooms_messages_belongs.message AND rooms_messages_belongs.room = :room)
					JOIN
						rooms ON (messages.room = rooms.id)
				`
			} else {
				sql += `
					JOIN
						rooms ON (messages.room = rooms.id AND rooms.id = :room)
				`
			}
		}

		sql += `
			WHERE
		`

		switch fetchMode {
		case "previous":
			sql += `
				messages.id <= :base_point AND
			`
		case "following":
			sql += `
				messages.id >= :base_point AND
			`
		}

		switch options.Category {
		case "follow":
			sql += `
				messages.character IN (SELECT * FROM follow_list UNION ALL SELECT :user_id) AND
			`
		case "follow-other":
			sql += `
				messages.character IN (SELECT * FROM follow_list) AND
			`
		case "replied":
			sql += `
				:user_id = ANY(messages.relates)  AND
			`
		case "character":
			sql += `
				messages.character = :target_id AND
			`
		case "character-replied":
			sql += `
				:target_id = ANY(messages.relates)  AND
			`
		case "replied-other":
			sql += `
				:user_id = ANY(messages.relates) AND
				messages.character != :user_id AND
			`
		case "own":
			sql += `
				messages.character = :user_id AND
			`
		case "conversation":
			sql += `
				messages.refer_root = :refer_root AND
			`
		case "search":
			sql += `
				messages.search_text LIKE likequery(:search) AND
			`
		case "list":
			sql += `
				messages.character IN (SELECT * FROM list_members) AND
			`
		}

		if options.RelateFilter {
			sql += `
				(
					messages.single = true OR
					:user_id = ANY(messages.relates) OR
					messages.relates <@ ARRAY(SELECT * FROM follow_list)
				) AND
			`
		}

		sql += `
				messages.deleted_at IS NULL AND
				NOT(messages.relates && ARRAY(SELECT * FROM block_list)) AND
				messages.character NOT IN (SELECT * FROM mute_list)
		`

		sql += `
			GROUP BY
				messages.id,
				messages.character,
				messages.refer,
				messages.refer_root,
				messages.secret,
				messages.icon,
				messages.name,
				messages.message,
				messages.replied_count,
				messages.posted_at,
				messages.reply_permission,
				rooms.id,
				rooms.title
			`

		if fetchMode == "previous" {
			sql += `
				ORDER BY
					messages.id DESC
			`
		} else {
			sql += `
				ORDER BY
					messages.id
			`
		}

		sql += `
			LIMIT
				:fetch_number;
		`

		return sql, &args
	}

	retrieveRangeMessages := func(rangeType string, number int) (*[]model.RoomMessage, bool, error) {
		sql, args := generateSql(rangeType)
		rows, err := db.NamedQuery(sql, *args)
		if err != nil {
			return nil, false, err
		}
		defer rows.Close()

		messages := make([]model.RoomMessage, 0, number+1)
		for rows.Next() {
			var message model.RoomMessage
			var recipientsJson string

			err = rows.Scan(
				&message.Id,
				&message.Character,
				&message.Refer,
				&message.ReferRoot,
				&message.Secret,
				&message.Icon,
				&message.Name,
				&message.Message,
				&message.PostedAt,
				&message.ReplyPermission,
				&message.RepliedCount,
				&message.Room.Id,
				&message.Room.Title,
				&message.Replyable,
				&recipientsJson,
			)

			if err != nil {
				return nil, false, err
			}

			err = json.Unmarshal([]byte(recipientsJson), &message.Recipients)
			if err != nil {
				return nil, false, err
			}

			if rangeType == "following" {
				messages = append([]model.RoomMessage{message}, messages...)
			} else {
				messages = append(messages, message)
			}
		}

		isContinue := len(messages) == number+1
		if isContinue {
			if rangeType == "following" {
				messages = messages[1:]
			} else {
				messages = messages[:len(messages)-1]
			}
		}

		return &messages, isContinue, nil
	}

	retrievePreviousMessages := func(number int) {
		messages, isContinue, err := retrieveRangeMessages("previous", number)
		previousResultChannel <- RetrieveMessagesResult{messages, isContinue, err}
	}

	retrieveFollowingMessages := func(number int) {
		messages, isContinue, err := retrieveRangeMessages("following", number)
		followingResultChannel <- RetrieveMessagesResult{messages, isContinue, err}
	}

	if options.RangeType == "initial" {
		go retrievePreviousMessages(options.FetchNumber + 1)
		go retrieveFollowingMessages(options.FetchNumber)

		previousResult := <-previousResultChannel
		followingResult := <-followingResultChannel

		if previousResult.err != nil || followingResult.err != nil {
			log.Println(previousResult.err)
			log.Println(followingResult.err)
			return nil, nil, nil, errors.New("データの取得中にエラーが発生しました")
		}

		messages := make([]model.RoomMessage, 0, len(*previousResult.messages)+len(*followingResult.messages))
		messages = append(messages, *followingResult.messages...)
		for i := range *previousResult.messages {
			if basePoint != (*previousResult.messages)[i].Id {
				messages = append(messages, (*previousResult.messages)[i])
			}
		}

		return &messages, &followingResult.isContinue, &previousResult.isContinue, nil
	} else if options.RangeType == "previous" || options.RangeType == "latest" {
		previousResultChannel = make(chan RetrieveMessagesResult)

		go retrievePreviousMessages(options.FetchNumber)

		previousResult := <-previousResultChannel

		if previousResult.err != nil {
			log.Println(previousResult.err)
			return nil, nil, nil, errors.New("データの取得中にエラーが発生しました")
		}

		return previousResult.messages, nil, &previousResult.isContinue, nil
	} else {
		followingResultChannel = make(chan RetrieveMessagesResult)

		go retrieveFollowingMessages(options.FetchNumber)

		followingResult := <-followingResultChannel

		if followingResult.err != nil {
			log.Println(followingResult.err)
			return nil, nil, nil, errors.New("データの取得中にエラーが発生しました")
		}

		return followingResult.messages, &followingResult.isContinue, nil, nil
	}
}
