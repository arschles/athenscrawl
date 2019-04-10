package queue

type TagCrawler struct {
	ch chan ModuleAndTag
}

func NewTagCrawler() *TagCrawler {
	return &TagCrawler{}
}

func (tc *TagCrawler) Add(module, tag string) {
	go func() {
		tc.ch <- ModuleAndTag{Module: module, Name: tag}
	}()
}

func (tc *TagCrawler) Ch() <-chan ModuleAndTag {
	return tc.ch
}
