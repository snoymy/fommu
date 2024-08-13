package mapper

import (
	"encoding/json"

	"github.com/snoymy/activitypub"
)

func JsonToActivity(body string) (*activitypub.Activity, error) {
    activity := &activitypub.Activity{}
    err := json.Unmarshal([]byte(body), &activity)
    if err != nil {
        return nil, err
    }

    var nestedScheme struct {
        Tag []*activitypub.Object `json:"tag"`
        Attachment []*activitypub.Object `json:"attachment"`
        Audience []*activitypub.Object `json:"audience"`
        To []*activitypub.Object `json:"to"`
        Bto []*activitypub.Object `json:"bto"`
        CC []*activitypub.Object `json:"cc"`
        BCC []*activitypub.Object `json:"bcc"`
    }

    err = json.Unmarshal([]byte(body), &nestedScheme)
    if err != nil {
        return nil, err
    }

    for _, item := range nestedScheme.Tag {
        activity.Tag.Append(item)
    }

    for _, item := range nestedScheme.Attachment {
        activity.Attachment.Append(item)
    }

    for _, item := range nestedScheme.Audience {
        activity.Audience.Append(item)
    }

    for _, item := range nestedScheme.To {
        activity.To.Append(item)
    }

    for _, item := range nestedScheme.Bto {
        activity.Bto.Append(item)
    }

    for _, item := range nestedScheme.CC {
        activity.CC.Append(item)
    }

    for _, item := range nestedScheme.BCC {
        activity.BCC.Append(item)
    }

    return activity, nil
}
