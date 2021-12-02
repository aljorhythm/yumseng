package rooms

import (
	"fmt"
	"github.com/aljorhythm/yumseng/cheers"
	"github.com/asaskevich/EventBus"
	"log"
)

type RoomEvents struct {
	eventBus               EventBus.Bus
	eventsCallbacksManager *EventsCallbacksManager
}

type EventType struct {
	name string
}

var (
	EVENT_CHEER_ADDED = EventType{
		"EVENT_CHEER_ADDED",
	}
	EVENT_ROOM_CONNECTED = EventType{
		"EVENT_ROOM_CONNECTED",
	}
	EVENT_LAST_SECONDS_COUNT = EventType{
		"EVENT_LAST_SECONDS_COUNT",
	}
)

func (eventType *EventType) topicName(room *Room) string {
	return fmt.Sprintf("%s:%s", room.Name, eventType.name)
}

func NewRoomEvents() *RoomEvents {
	bus := EventBus.New()
	eventCallbacksManager := newEventsCallbackManager()

	return &RoomEvents{
		bus,
		eventCallbacksManager,
	}
}

type Callback func(args ...interface{})

type EventCallbacks struct {
	callbacks *map[string]Callback
	topic     string
}

func (callbacks *EventCallbacks) callbackAll(args ...interface{}) {
	for callbackId, callback := range *callbacks.callbacks {
		log.Printf("calling callback %s in topic %s", callbackId, callbacks.topic)
		callback(args...)
	}
}

func (callbacks *EventCallbacks) addCallback(callbackId string, callback func(...interface{})) {
	log.Printf("topic %s adding callback %s", callbackId, callbacks.topic)
	(*callbacks.callbacks)[callbackId] = callback
}

func (callbacks *EventCallbacks) removeCallback(callbackId string) {
	log.Printf("removing callback %s to %s", callbackId, callbacks.topic)
	delete(*callbacks.callbacks, callbackId)
}

type EventsCallbacksManager struct {
	eventCallbacksList map[string]*EventCallbacks
}

func (manager *EventsCallbacksManager) getEventCallbacks(topic string) *EventCallbacks {
	eventCallbacks, _ := manager.eventCallbacksList[topic]
	return eventCallbacks
}

func (manager *EventsCallbacksManager) addEventCallback(topic string, callbackId string, callback Callback) *EventCallbacks {
	eventCallbacks, ok := manager.eventCallbacksList[topic]

	if !ok {
		log.Printf("manager creating event callbacks topic: %s", topic)
		manager.eventCallbacksList[topic] = newEventCallbacks(topic)
	}

	eventCallbacks, ok = manager.eventCallbacksList[topic]

	if !ok {
		log.Panicf("wanted to create callbacks but failed topic: %s callbackId: %s", topic, callbackId)
	} else {
		log.Printf("manager adding callback %s to %s", callbackId, topic)
		eventCallbacks.addCallback(callbackId, callback)
	}

	return eventCallbacks
}

func (manager *EventsCallbacksManager) removeEventCallback(topic string, callbackId string) {
	eventCallbacks, ok := manager.eventCallbacksList[topic]

	if ok {
		log.Printf("removing callback %s from %s", callbackId, topic)
		eventCallbacks.removeCallback(callbackId)
	} else {
		log.Panicf("callbacks do not exist to remove: %s callbackId: %s", topic, callbackId)
	}
}

func newEventsCallbackManager() *EventsCallbacksManager {
	var eventCallbacksList = map[string]*EventCallbacks{}
	return &EventsCallbacksManager{eventCallbacksList: eventCallbacksList}
}

func newEventCallbacks(topic string) *EventCallbacks {
	callbacksMap := &map[string]Callback{}
	callbacks := &EventCallbacks{
		callbacksMap,
		topic,
	}
	return callbacks
}

func (roomEvents *RoomEvents) SubscribeCheerAdded(room *Room, cliendId string, cb Callback) {
	topic := EVENT_CHEER_ADDED.topicName(room)
	eventCallbacks := roomEvents.eventsCallbacksManager.addEventCallback(topic, cliendId, cb)

	if !roomEvents.eventBus.HasCallback(topic) {
		log.Printf("subcribing event callbacks to topic %s", topic)
		roomEvents.eventBus.Subscribe(topic, eventCallbacks.callbackAll)
	}

	log.Printf("subscribed event callback %s to topic %s", cliendId, topic)
}

func (roomEvents *RoomEvents) PublishCheerAdded(room *Room, cheer cheers.Cheer) {
	topic := EVENT_CHEER_ADDED.topicName(room)
	log.Printf("publishing event rooms %s | topic %s", room.Name, topic)
	roomEvents.eventBus.Publish(topic, cheer)
}

func (roomEvents *RoomEvents) UnsubscribeCheerAdded(room *Room, clientId string) {
	topic := EVENT_CHEER_ADDED.topicName(room)
	log.Printf("unsubcribing rooms %s to %s", room.Name, topic)
	roomEvents.eventsCallbacksManager.removeEventCallback(topic, clientId)
}
