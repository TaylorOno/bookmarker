package repository

type UserBookmark struct {
	UserId      string
	LastUpdated string
	Book        string
	Series      string
	Status      string
	Page        int

	//AdditionalProperties provide for extendable data model there are no guarantees on any fields provided.  Data will
	//be projected into secondary indexes so be cautious of field size.
	AdditionalProperties map[string]interface{}
}
