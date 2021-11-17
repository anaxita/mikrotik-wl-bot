package robot

func (b *Robot) roleByUsername(userName string) int {
	for name := range b.store.Admins {
		if name == userName {
			return roleAdmin
		}
	}

	return roleUser
}

func (b *Robot) isChatAllow(chatID int64) bool {
	for _, id := range b.store.AllowChatIDs {
		if id == chatID {
			return true
		}
	}

	return false
}

func (b *Robot) isAdmin(username string) bool {
	for admin := range b.store.Admins {
		if username == admin {
			return true
		}
	}

	return false
}
