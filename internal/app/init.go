package app

import "strconv"

func init() {
	if isDev, err := strconv.ParseBool(IsDevelopment); err == nil && isDev {
		IsDevBuild = true
	} else {
		IsDevBuild = false
	}
}
