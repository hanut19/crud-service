package constants

const (
	DB_USERNAME = "amandeepbaghoria"
	DB_PASSWORD = "gm4iEWef7jywSfGc"
	SECRETKEY   = "secretkeyjwt"
	SERVER_PORT = "8081"
	SERVER_URL  = "http://localhost:"
	ADMIN       = "admin"
	USER        = "user"
	CREATE      = "create"
	READ        = "read"
	EDIT        = "edit"
	DELETE      = "delete"
)

func GetRole() []string {
	return []string{ADMIN, USER}
}

func GetPermission(role string) map[string]bool {
	var retmap map[string]bool
	permission := make(map[string]map[string]bool)

	admin_map := map[string]bool{CREATE: true, READ: true, EDIT: true, DELETE: true}
	user_map := map[string]bool{CREATE: false, READ: true, EDIT: false, DELETE: false}

	permission[ADMIN] = admin_map
	permission[USER] = user_map
	if val, ok := permission[role]; ok {
		return val
	}
	return retmap
}

func IsAccess(role string, permissionType string) bool {
	resp := GetPermission(role)
	if permissionStatus, ok := resp[permissionType]; ok {
		return permissionStatus
	}
	return false
}
