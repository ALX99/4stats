package board

func getHighestPostNo(c catalog) int64 {
	var maxNo int64
	for _, page := range c {
		for _, thread := range page.Threads {
			if maxNo < thread.No {
				maxNo = thread.No
			}
			for _, reply := range thread.LastReplies {
				if maxNo < reply.No {
					maxNo = reply.No
				}
			}
		}
	}
	return maxNo
}
