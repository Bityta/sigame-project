package player

type Role string

const (
	RoleHost   Role = "host"
	RolePlayer Role = "player"
)

func (r Role) String() string {
	return string(r)
}

func (r Role) IsHost() bool {
	return r == RoleHost
}

func (r Role) IsPlayer() bool {
	return r == RolePlayer
}

