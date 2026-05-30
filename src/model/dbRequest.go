package model

type CardRequestHeader struct {
	Id int `json:"id"`
}

type CardCollectionByKeywordHeader struct {
	Keyword string `json:"keyword"`
}

type CardsCollectionBody struct {
	Cards []CardAPI `json:"Cards"`
	Total int       `json:"total"`
}

// REQUESTS
//
//	GET, DELETE
type GetDeleteCardRequest struct {
	Header CardRequestHeader
}

type GetCardsByKeywordRequest struct {
	Header CardCollectionByKeywordHeader
}

// POST, PUT
type PostPutCardRequest struct {
	Body CardAPI
}
type FilterCardRequest struct {
	CardAPI
}

// RESPONSE
//
//	GET
type GetCardResponse struct {
	Body CardAPI
	Err  error
}

type NoContentCardResponse struct {
	Err error
}

type ExampleCardResponse struct {
	Body []CardAPI
	Err  error
}

type CardsCollectionResponse struct {
	Body CardsCollectionBody
	Err  error
}
