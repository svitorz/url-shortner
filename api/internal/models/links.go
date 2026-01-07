package models

type Link struct {
	slug       string
	target_url string
	owner      User
	is_active  bool
}
