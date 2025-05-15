package notifications

import "log/slog"

/*
Sync merges the local and remote notifications.

It applies the following rules:

	| remote \ local | Missing    | Exist      | Done       | Hidden   |
	| ---            | ---        | ---        | ---        | ---      |
	| Exist          | (1) Insert | (2) Update | (2) Update | (3) Keep |
	| Missing        |            | (3) Keep   | (4) Drop   | (4) Drop |

	(1) Insert: Add the notification ass is.
	(2) Update: Update the local notification with the remote data, keep the Meta
	    unchanged.
	(3) Keep: Keep the local notification unchanged.
	(4) Drop: Remove the notification from the local list.

Notes on (2) Update: Updating the notification will also reset the `Meta.Done`
state if the remote notification is newer than the local one.

TODO: refactor this to `func (n Notifications) Sync(remote Notifications) {}`.
*/
//revive:disable:cognitive-complexity // There's enough comments/details to keep
// it all here.
func Sync(local, remote Notifications) Notifications {
	// TODO: do we need to have the whole map?
	remoteMap := remote.Map()
	localMap := local.Map()

	n := Notifications{}

	// Add any new notifications to the list
	for i := range remote {
		if _, ok := localMap[remote[i].ID]; !ok {
			// (1) Insert
			slog.Debug("sync", "action", "insert", "id", remote[i].ID)

			remote[i].Meta.RemoteExists = true
			n = append(n, remote[i])
		}
	}

	for i := range local {
		remote, remoteExist := remoteMap[local[i].ID]

		local[i].Meta.RemoteExists = remoteExist

		if remoteExist {
			// (3) Keep
			if local[i].Meta.Hidden {
				slog.Debug("sync", "action", "keeping hidden", "id", local[i].ID)
				n = append(n, local[i])

				continue
			}

			// (2) Update
			slog.Debug("sync", "action", "update", "id", remote.ID)

			if local[i].Meta.Done && remote.UpdatedAt.After(local[i].UpdatedAt) {
				slog.Debug("sync", "action", "resetting done", "id", local[i].ID)
				local[i].Meta.Done = false
			}

			remote.Meta = local[i].Meta
			n = append(n, remote)
		} else {
			if local[i].Meta.Done || local[i].Meta.Hidden {
				// (4) Drop
				slog.Debug("sync", "action", "drop", "id", local[i].ID)

				continue
			}

			// (3) Keep
			slog.Debug("sync", "action", "keep", "id", local[i].ID)
			n = append(n, local[i])
		}
	}

	n.Sort()

	// TODO: add uniq here

	return n
}
