package notifications

import "log/slog"

/*
Sync merges the local and remote notifications.

It applies the following rules:

	| remote \ local | Missing    | Exist      | Done       | Hidden   |
	| ---            | ---        | ---        | ---        | ---      |
	| Exist          | (1) Insert | (2) Update | (2) Update | (3) Keep |
	| Missing        | (3) Keep   | (3) Keep   | (4) Drop   | (4) Drop |

	(1) Insert: Add the notification ass is.
	(2) Update: Update the local notification with the remote data, keep the Meta
	    unchanged.
	(3) Keep: Keep the local notification unchanged.
	(4) Drop: Remove the notification from the local list.

Notes on (2) Update: Updating the notification will also reset the `Meta.Done`
state if the remote notification is newer than the local one.

TODO: refactor this to `func (n Notifications) Sync(remote Notifications) {}`
*/
func Sync(local, remote Notifications) Notifications {
	remoteMap := remote.Map()
	localMap := local.Map()

	n := Notifications{}

	// Add any new notifications to the list
	for remoteId, remote := range remoteMap {
		if _, ok := localMap[remoteId]; !ok {
			// (1) Insert
			slog.Debug("sync", "action", "insert", "id", remote.Id)

			remote.Meta.RemoteExists = true
			n = append(n, remote)
		}
	}

	for localId, local := range localMap {
		remote, remoteExist := remoteMap[localId]

		local.Meta.RemoteExists = remoteExist

		if remoteExist {
			// (3) Keep
			if local.Meta.Hidden {
				slog.Debug("sync", "action", "keeping hidden", "id", local.Id)
				n = append(n, local)
				continue
			}

			// (2) Update
			slog.Debug("sync", "action", "update", "id", remote.Id)

			if local.Meta.Done && remote.UpdatedAt.After(local.UpdatedAt) {
				slog.Debug("sync", "action", "reseting done", "id", local.Id)
				local.Meta.Done = false
			}

			remote.Meta = local.Meta
			n = append(n, remote)
		} else {
			if local.Meta.Done || local.Meta.Hidden {
				// (4) Drop
				slog.Debug("sync", "action", "drop", "id", local.Id)
				continue
			}

			// (3) Keep
			slog.Debug("sync", "action", "keep", "id", local.Id)
			n = append(n, local)
		}
	}

	n.Sort()

	// TODO: add uniq here

	return n
}
