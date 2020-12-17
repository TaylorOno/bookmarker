package service

type NewBookmarkRequest struct {
	UserId string `json:"userId" validate:"required"`
	Book   string `json:"book" validate:"required"`
	Series string `json:"series"`
	Status string `json:"status" validate:"oneof=IN_PROGRESS FINISHED"`
	Page   int    `json:"Page" validate:"gte=0"`

	//AdditionalProperties provide for extendable data model there are no guarantees on any fields provided.  Data will
	//be projected into secondary indexes so be cautious of field size.
	AdditionalProperties map[string]interface{} `json:"additionalProperties"`
}

type DeleteBookmarkRequest struct {
	UserId string
	Book   string
}

type BookmarkRequest struct {
	UserId string
	Book   string
}

type Bookmark struct {
	Book                 string                 `json:"book"`
	LastUpdated          string                 `json:"lastUpdated,omitempty"`
	Series               string                 `json:"series,omitempty"`
	Status               string                 `json:"status"`
	Page                 int                    `json:"page"`
	AdditionalProperties map[string]interface{} `json:"additionalProperties,omitempty"`
}

type BookmarkListRequest struct {
	UserId string `validate:"required"`
	Limit  int64  `validate:"required"`
	Filter string `validate:"oneof=NONE IN_PROGRESS FINISHED"`
}
