package rabbitmq

import (
	"fmt"

	pb "github.com/ncfex/dcart-auth/internal/adapters/secondary/messaging/proto"
	"github.com/ncfex/dcart-auth/internal/domain/shared"
	"github.com/ncfex/dcart-auth/internal/domain/user"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func SerializeEvent(event shared.Event) (*pb.EventMessage, error) {
	var payload []byte
	var err error

	switch e := event.(type) {
	case *user.UserRegisteredEvent:
		protoEvent := &pb.UserRegisteredEvent{
			Base: &pb.BaseEvent{
				AggregateId:   e.GetAggregateID(),
				AggregateType: e.GetAggregateType(),
				EventType:     e.GetEventType(),
				Version:       int32(e.GetVersion()),
				Timestamp:     timestamppb.New(e.GetTimestamp()),
			},
			Username:     e.Username,
			PasswordHash: e.PasswordHash,
		}
		payload, err = proto.Marshal(protoEvent)
	case *user.UserPasswordChangedEvent:
		protoEvent := &pb.UserPasswordChangedEvent{
			Base: &pb.BaseEvent{
				AggregateId:   e.GetAggregateID(),
				AggregateType: e.GetAggregateType(),
				EventType:     e.GetEventType(),
				Version:       int32(e.GetVersion()),
				Timestamp:     timestamppb.New(e.GetTimestamp()),
			},
			NewPasswordHash: e.NewPasswordHash,
		}
		payload, err = proto.Marshal(protoEvent)
	default:
		return nil, fmt.Errorf("unknown event type: %T", event)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to marshal event: %w", err)
	}

	return &pb.EventMessage{
		AggregateId:   event.GetAggregateID(),
		AggregateType: event.GetAggregateType(),
		EventType:     event.GetEventType(),
		Version:       int32(event.GetVersion()),
		Timestamp:     timestamppb.New(event.GetTimestamp()),
		Payload:       payload,
	}, nil
}

func DeserializeEvent(msg *pb.EventMessage, registry shared.EventRegistry) (shared.Event, error) {
	baseEvent := shared.BaseEvent{
		AggregateID:   msg.AggregateId,
		AggregateType: msg.AggregateType,
		EventType:     msg.EventType,
		Version:       int(msg.Version),
		Timestamp:     msg.Timestamp.AsTime(),
	}

	switch shared.EventType(msg.EventType) {
	case user.EventTypeUserRegistered:
		var protoEvent pb.UserRegisteredEvent
		if err := proto.Unmarshal(msg.Payload, &protoEvent); err != nil {
			return nil, fmt.Errorf("failed to unmarshal UserRegisteredEvent: %w", err)
		}
		return &user.UserRegisteredEvent{
			BaseEvent:    baseEvent,
			Username:     protoEvent.Username,
			PasswordHash: protoEvent.PasswordHash,
		}, nil
	case user.EventTypeUserPasswordChanged:
		var protoEvent pb.UserPasswordChangedEvent
		if err := proto.Unmarshal(msg.Payload, &protoEvent); err != nil {
			return nil, fmt.Errorf("failed to unmarshal UserPasswordChangedEvent: %w", err)
		}
		return &user.UserPasswordChangedEvent{
			BaseEvent:       baseEvent,
			NewPasswordHash: protoEvent.NewPasswordHash,
		}, nil
	default:
		return nil, fmt.Errorf("unknown event type: %s", msg.EventType)
	}
}
