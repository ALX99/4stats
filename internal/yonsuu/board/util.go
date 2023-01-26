package board

func getHighestPostNo(page indexPage) int {
	maxNo := 0
	for _, prevThreads := range page.Threads {
		for _, prevPost := range prevThreads.Posts {
			if maxNo < prevPost.No {
				maxNo = prevPost.No
			}
		}
	}
	return maxNo
}
