package app

import (
	"github.com/rs/zerolog/log"
	"os/user"
	"strconv"
	"strings"
	"syscall"
)

func dropPrivileges(username string) {
	if len(strings.TrimSpace(username)) > 0 {
		if syscall.Getuid() == 0 {
			log.Info().Msgf("Running as root, downgrading to currentUser %s", username)
			currentUser, err := user.Lookup(username)
			if err != nil {
				log.Fatal().Msgf("User not found or other error: %v", err)
			}
			// TODO: Write error handling for int from string parsing
			uid, _ := strconv.ParseInt(currentUser.Uid, 10, 32)
			gid, _ := strconv.ParseInt(currentUser.Gid, 10, 32)

			//cerr, errno := C.setgid(C.__gid_t(gid))
			if err := syscall.Setgid(int(gid)); err != nil {
				log.Fatal().Msgf("Unable to set GID due to error: %d", err)
			}

			//cerr, errno = C.setuid(C.__uid_t(uid))
			if err := syscall.Setuid(int(uid)); err != nil {
				log.Fatal().Msgf("Unable to set UID due to error: %d", err)
			}
		}
	}
}
