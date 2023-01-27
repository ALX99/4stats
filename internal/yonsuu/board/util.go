package board

func getHighestPostNo(page indexPage) int64 {
	var maxNo int64
	for _, threads := range page.Threads {
		for _, post := range threads.Posts {
			if maxNo < post.No {
				maxNo = post.No
			}
		}
	}
	return maxNo
}
