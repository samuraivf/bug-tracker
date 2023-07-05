package handler

const (
	id    = "/:id"
	empty = "/"

	auth     = "/auth"
	signUp   = "/sign-up"
	signIn   = "/sign-in"
	logout   = "/logout"
	refresh  = "/refresh"
	verify   = "/verify-email"
	setEmail = "/set-email"

	project      = "/project"
	create       = "/create"
	update       = "/update"
	addMember    = "/add-member"
	deleteMember = "/member"
	leave        = "/leave/:id"
	setAdmin     = "/set-admin"
	withTasks    = "/with-tasks" + id

	task           = "/task"
	workOnTask     = "/work-on-task"
	stopWorkOnTask = "/stop-work-on-task"
)
