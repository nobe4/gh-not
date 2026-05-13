package notifications

// MergeUpdatedNotification merges n into o, preserving enrichment if still fresh.
func (n *Notification) MergeUpdatedNotification(o *Notification) *Notification {
	meta := n.Meta
	meta.RemoteExists = true

	if o.UpdatedAt.After(n.UpdatedAt) {
		meta.Done = false
		meta.Enriched = false

		o.clearEnrichment()
	} else if meta.Enriched {
		o.mergeEnrichment(n)
	}

	o.Meta = meta

	return o
}

func (n *Notification) mergeEnrichment(o *Notification) {
	n.Author = o.Author
	n.LatestCommentor = o.LatestCommentor
	n.Assignees = o.Assignees
	n.Reviewers = o.Reviewers
	n.ReviewersTeams = o.ReviewersTeams
	n.MergedBy = o.MergedBy
	n.Subject.State = o.Subject.State
	n.Subject.HTMLURL = o.Subject.HTMLURL
}

func (n *Notification) clearEnrichment() {
	n.mergeEnrichment(&Notification{})
}
