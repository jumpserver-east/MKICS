package router

func commonGroups() []CommonRouter {
	return []CommonRouter{
		&BaseRouter{},
		&WecomRouter{},
		&MaxkbRouter{},
		&KFRouter{},
		&StaffRouter{},
		&PolicyRouter{},
	}
}
