package account

type Avatar string

func (a Avatar) String() string {
	return string(a)
}

func NewAvatar(avt string) (Avatar, error) {
	return Avatar(avt), nil
}
