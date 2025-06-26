package snowflake_proxy

import (
	"fmt"

	sdnotify "github.com/coreos/go-systemd/v22/daemon"
)

func sdnotifyReady() {
	sdnotify.SdNotify(false, sdnotify.SdNotifyReady)
}

func sdnotifyStopping() {
	sdnotify.SdNotify(false, sdnotify.SdNotifyStopping)
}

func sdnotifyStatus(status string) {
	sdnotify.SdNotify(false, fmt.Sprintf("STATUS=%s", status))
}

func sdnotifyWatchdog() {
	sdnotify.SdNotify(false, sdnotify.SdNotifyWatchdog)
}
