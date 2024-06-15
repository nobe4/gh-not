package notifications

// Sync merges the local and remote notifications.
//
// It applies the following rules:
// | remote \ local | Missing   | Exist     | ToDelete  |
// | ---            | ---       | ---       | ---       |
// | Exist          | (1)Insert | (2)Update | (2)Update |
// | Missing        | (3)Noop   | (3)Noop   | (4)Drop   |
//
//  1. Insert: Add the notification ass is.
//  2. Update: Update the local notification with the remote data, keep the Meta
//     unchanged.
//  3. Noop: Do nothing.
//  4. Drop: Remove the notification from the local list.
func Sync(local, remote Notifications) Notifications {
	remoteMap := remote.Map()
	localMap := local.Map()

	n := Notifications{}

	// Add any new notifications to the list
	for remoteId, remote := range remoteMap {
		if _, ok := localMap[remoteId]; !ok {
			// (1)Insert
			n = append(n, remote)
		}
	}

	for localId, local := range localMap {
		remote, remoteExist := remoteMap[localId]

		if remoteExist {
			// (2)Update
			remote.Meta = local.Meta
			n = append(n, remote)
		} else {
			if local.Meta.ToDelete {
				// (4)Drop
				continue
			}

			// (3)Noop
			n = append(n, local)
		}
	}

	n.Sort()

	return n
}
