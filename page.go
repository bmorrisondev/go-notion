package gonotion

type Page struct {
	Object         *string `json:"object,omitempty"`
	Id             *string `json:"id,omitempty"`
	CreatedTime    *string `json:"createdTime,omitempty"`
	LastEditedTime *string `json:"lastEditedTime,omitempty"`
	Archived       *bool   `json:"archived,omitempty"`
	Icon           *FileObject
	Cover          *FileObject
	Properties     *map[string]Property `json:"properties,omitempty"`
	Url            *string              `json:"url,omitempty"`
}

func (p *Page) GetTitle() string {
	props := *p.Properties
	titleProp := props["Name"]
	titleArr := *titleProp.Title
	agg := ""
	for _, el := range titleArr {
		agg += *el.Text.Content
	}
	return agg
}
