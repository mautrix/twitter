package types

import (
	"encoding/json"
	"maps"
	"slices"

	"github.com/rs/zerolog"
)

type TwitterEvent interface {
	isTwitterEvent()
}

func (*Message) isTwitterEvent()                    {}
func (*MessageEdit) isTwitterEvent()                {}
func (*MessageDelete) isTwitterEvent()              {}
func (*MessageReactionCreate) isTwitterEvent()      {}
func (*MessageReactionDelete) isTwitterEvent()      {}
func (*ConversationCreate) isTwitterEvent()         {}
func (*ConversationDelete) isTwitterEvent()         {}
func (*ParticipantsJoin) isTwitterEvent()           {}
func (*ConversationMetadataUpdate) isTwitterEvent() {}
func (*ConversationNameUpdate) isTwitterEvent()     {}
func (*ConversationRead) isTwitterEvent()           {}
func (*TrustConversation) isTwitterEvent()          {}
func (*EndAVBroadcast) isTwitterEvent()             {}
func (*PollingError) isTwitterEvent()               {}

type PollingError struct {
	Error error
}

type RawTwitterEvent []byte

func (rte *RawTwitterEvent) UnmarshalJSON(data []byte) error {
	*rte = data
	return nil
}

type twitterEventContainer struct {
	Message                    *Message                    `json:"message,omitempty"`
	MessageDelete              *MessageDelete              `json:"message_delete,omitempty"`
	MessageEdit                *MessageEdit                `json:"message_edit,omitempty"`
	ReactionCreate             *MessageReactionCreate      `json:"reaction_create,omitempty"`
	ReactionDelete             *MessageReactionDelete      `json:"reaction_delete,omitempty"`
	ConversationCreate         *ConversationCreate         `json:"conversation_create,omitempty"`
	ConversationDelete         *ConversationDelete         `json:"remove_conversation,omitempty"`
	ParticipantsJoin           *ParticipantsJoin           `json:"participants_join,omitempty"`
	ConversationMetadataUpdate *ConversationMetadataUpdate `json:"conversation_metadata_update,omitempty"`
	ConversationNameUpdate     *ConversationNameUpdate     `json:"conversation_name_update,omitempty"`
	ConversationRead           *ConversationRead           `json:"conversation_read,omitempty"`
	TrustConversation          *TrustConversation          `json:"trust_conversation,omitempty"`
	EndAVBroadcast             *EndAVBroadcast             `json:"end_av_broadcast,omitempty"`
	// DisableNotifications       *types.DisableNotifications       `json:"disable_notifications,omitempty"`
}

func (rte *RawTwitterEvent) Parse() (TwitterEvent, map[string]any, error) {
	var tec twitterEventContainer
	if err := json.Unmarshal(*rte, &tec); err != nil {
		return nil, nil, err
	}
	switch {
	case tec.Message != nil:
		return tec.Message, nil, nil
	case tec.MessageDelete != nil:
		return tec.MessageDelete, nil, nil
	case tec.MessageEdit != nil:
		return tec.MessageEdit, nil, nil
	case tec.ReactionCreate != nil:
		return tec.ReactionCreate, nil, nil
	case tec.ReactionDelete != nil:
		return tec.ReactionDelete, nil, nil
	case tec.ConversationCreate != nil:
		return tec.ConversationCreate, nil, nil
	case tec.ConversationDelete != nil:
		return tec.ConversationDelete, nil, nil
	case tec.ParticipantsJoin != nil:
		return tec.ParticipantsJoin, nil, nil
	case tec.ConversationMetadataUpdate != nil:
		return tec.ConversationMetadataUpdate, nil, nil
	case tec.ConversationNameUpdate != nil:
		return tec.ConversationNameUpdate, nil, nil
	case tec.ConversationRead != nil:
		return tec.ConversationRead, nil, nil
	case tec.TrustConversation != nil:
		return tec.TrustConversation, nil, nil
	case tec.EndAVBroadcast != nil:
		return tec.EndAVBroadcast, nil, nil
	default:
		var unrecognized map[string]any
		if err := json.Unmarshal(*rte, &unrecognized); err != nil {
			return nil, nil, err
		}
		return nil, unrecognized, nil
	}
}

func (rte *RawTwitterEvent) ParseWithErrorLog(log *zerolog.Logger) TwitterEvent {
	evt, unrecognized, err := rte.Parse()
	if err != nil {
		logEvt := log.Err(err)
		if log.GetLevel() == zerolog.TraceLevel {
			logEvt.RawJSON("entry_data", *rte)
		}
		logEvt.Msg("Failed to parse entry")
	} else if unrecognized != nil {
		logEvt := log.Warn().Strs("type_keys", slices.Collect(maps.Keys(unrecognized)))
		if log.GetLevel() == zerolog.TraceLevel {
			logEvt.Any("entry_data", unrecognized)
		}
		logEvt.Msg("Unrecognized entry type")
	} else {
		log.Trace().RawJSON("entry_data", *rte).Msg("Parsed entry")
		return evt
	}
	return nil
}
