package controller

import (
	"encoding/base64"
	"encoding/json"
	"personaapp/internal/server/controllers/vacancy/storage"
	"time"
)

func toCursorData(cursor *Cursor) (*cursorData, error) {
	if cursor == nil {
		return nil, nil
	}

	var cursorData cursorData

	decoded, err := base64.StdEncoding.DecodeString(string(*cursor))
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(decoded, &cursorData); err != nil {
		return nil, err
	}

	return &cursorData, nil
}

func toCursor(cursor *storage.Cursor, categoriesIDs []string) (*Cursor, error) {
	if cursor == nil {
		return nil, nil
	}

	cursorData := cursorData{
		PrevCreatedAt: cursor.PrevCreatedAt,
		PrevPosition:  cursor.PrevPosition,
		CategoriesIDs: categoriesIDs,
	}

	data, err := json.Marshal(cursorData)
	if err != nil {
		return nil, err
	}

	c := Cursor(base64.StdEncoding.EncodeToString(data))

	return &c, nil
}

func toStorageCursor(cursorData *cursorData) *storage.Cursor {
	if cursorData == nil {
		return nil
	}

	return &storage.Cursor{
		PrevCreatedAt: cursorData.PrevCreatedAt,
		PrevPosition:  cursorData.PrevPosition,
	}
}

// cursor data
type cursorData struct {
	PrevCreatedAt time.Time `json:"created_at,string"`
	PrevPosition  int       `json:"position"`
	CategoriesIDs []string  `json:"categories"`
}

type basicCursorData struct {
	PrevCreatedAt string   `json:"created_at"`
	PrevPosition  int      `json:"int,string"`
	CategoriesIDs []string `json:"categories"`
}

func (cd cursorData) MarshalJSON() ([]byte, error) {
	return json.Marshal(basicCursorData{
		PrevCreatedAt: cd.PrevCreatedAt.Format(time.RFC3339Nano),
		PrevPosition:  cd.PrevPosition,
		CategoriesIDs: cd.CategoriesIDs,
	})
}

func (cd *cursorData) UnmarshalJSON(j []byte) error {
	var data basicCursorData
	if err := json.Unmarshal(j, &data); err != nil {
		return err
	}

	createdAt, err := time.Parse(time.RFC3339Nano, data.PrevCreatedAt)
	if err != nil {
		return err
	}

	cd.PrevPosition = data.PrevPosition
	cd.PrevCreatedAt = createdAt
	cd.CategoriesIDs = data.CategoriesIDs

	return nil
}
