package activitypub_extended

import (
	"encoding/json"

	"github.com/snoymy/activitypub"
)

const PropertyValueType activitypub.ActivityVocabularyType = "PropertyValue"

type PropertyValue struct {
	ID          activitypub.ID                      `json:"id,omitempty"`
	Type        activitypub.ActivityVocabularyType  `json:"type,omitempty"`
	Name        activitypub.NaturalLanguageValues   `json:"name,omitempty"`
	Value       activitypub.NaturalLanguageValues   `json:"value,omitempty"`
}

func PropertyValueNew(id activitypub.ID) *PropertyValue {
    return &PropertyValue{
        ID: id,
        Type: PropertyValueType,
    }
}

func (o PropertyValue) GetID() activitypub.ID {
    return o.ID
}

func (o PropertyValue) GetType() activitypub.ActivityVocabularyType {
    return o.Type
}

func (o PropertyValue) GetLink() activitypub.IRI {
    return activitypub.IRI(o.ID)
}

func (o PropertyValue) IsLink() bool {
    return false
}

func (o PropertyValue) IsObject() bool {
    return true
}

func (o PropertyValue) IsCollection() bool {
    return false
}

func (o PropertyValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(&map[string]string{
        "type": string(o.Type),
        "name": o.Name.String(),
        "value": o.Value.String(),
    })
}
