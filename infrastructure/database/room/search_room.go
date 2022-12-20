package room

type SearchRoomOptions struct {
	Order          string
	Sort           string
	Offset         int
	Limit          int
	Participating  int
	Participatable bool
}

type SearchRoomOption func(*SearchRoomOptions)

func SearchRoomOptionOrder(order string) SearchRoomOption {
	return func(args *SearchRoomOptions) {
		if order == "asc" || order == "desc" {
			args.Order = order
		}
	}
}

func SearchRoomOptionSort(sort string) SearchRoomOption {
	return func(args *SearchRoomOptions) {
		if sort == "rno" || sort == "update-time" || sort == "responses" {
			args.Sort = sort
		}
	}
}

func SearchRoomOptionOffset(offset int) SearchRoomOption {
	return func(args *SearchRoomOptions) {
		args.Offset = offset
	}
}

func SearchRoomOptionLimit(limit int) SearchRoomOption {
	return func(args *SearchRoomOptions) {
		args.Offset = limit
	}
}

func SearchRoomOptionParticipating(state int) SearchRoomOption {
	return func(args *SearchRoomOptions) {
		args.Participating = state
		// -1: 非参加中のもののみ, 0: 条件指定なし, 1: 参加中のもののみ
	}
}

func SearchRoomOptionParticipatable(state bool) SearchRoomOption {
	return func(args *SearchRoomOptions) {
		args.Participatable = state
	}
}
