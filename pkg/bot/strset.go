package bot

type strFilter func(string) bool

func strSet(items []string) strFilter {
	var set = make(map[string]struct{}, len(items)/4)
	for _, item := range items {
		set[item] = struct{}{}
	}
	return func(str string) bool {
		var _, inSet = set[str]
		return inSet
	}
}

func (filter strFilter) FilterList(list []string) []string {
	var filtered = make([]string, 0, len(list))
	for _, item := range list {
		if filter(item) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func (filter strFilter) FilterStream(in <-chan string) <-chan string {
	var out = make(chan string)
	go func() {
		defer close(out)
		for item := range in {
			if filter(item) {
				out <- item
			}
		}
	}()
	return out
}

func andFilter(filters ...strFilter) strFilter {
	return func(str string) bool {
		var ok = true
		for _, filter := range filters {
			ok = ok && filter(str)
		}
		return ok
	}
}

func orFilter(filters ...strFilter) strFilter {
	return func(str string) bool {
		var ok = false
		for _, filter := range filters {
			ok = ok || filter(str)
		}
		return ok
	}
}
