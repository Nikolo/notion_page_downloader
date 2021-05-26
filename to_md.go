package main

import (
	"path/filepath"

	"github.com/Nikolo/notionapi"
	"github.com/Nikolo/notion_page_downloader/tomarkdown"
	"github.com/kjk/u"
)

func mdPath(page *notionapi.Page) string {
	//pageID = notionapi.ToNoDashID(pageID)
	//name := fmt.Sprintf("%s.%d.page.md", pageID, n)
	name := tomarkdown.MarkdownFileNameForPage(page) //fmt.Sprintf("%s.md", pageID)
	return filepath.Join(dataDir, name)
}

func toMd(pageID string) *notionapi.Page {
	client := makeNotionClient()
	page, err := downloadPage(client, pageID)
	if err != nil {
		logf("toMd: downloadPage() failed with '%s'\n", err)
		path := filepath.Join(dataDir, tomarkdown.MarkdownFileNameForPageID(pageID))
		u.WriteFileMust(path, []byte("Unable to download"))
		return nil
	}
	if page == nil {
		logf("toMd: page is nil\n")
		return nil
	}

	notionapi.PanicOnFailures = true

	c := tomarkdown.NewConverter(page)
	md := c.ToMarkdown()
	path := mdPath(page)
	u.WriteFileMust(path, md)
        if page.IsSubPage(page.Root()) {
		queueDownload = append(queueDownload, page.GetSubPages()...)
        }
	for _, r := range page.SpaceRecords {
		queueDownload = append(queueDownload, r.Space.Pages...)
		logf("Add to queue %d pages from space %s\n", len(r.Space.Pages), r.Space.ID)
	}
	logf("queue size: %d\n", len(queueDownload))
/*
	if !flgNoOpen {
		path, err := filepath.Abs(path)
		must(err)
		uri := "file://" + path
		logf("Opening browser with %s\n", uri)
		u.OpenBrowser(uri)
	}
*/
	return page
}
