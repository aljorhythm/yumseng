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
	// todo rename to intensity
	EVENT_LAST_SECONDS_COUNT = EventType{
		"EVENT_LAST_SECONDS_COUNT",
	}
)

func (eventType *EventType) topicName(room *Room) string {
	return fmt.Sprintf("%s:%s", room.Name, eventType.name)
}

func (eventType *EventType) GetName() string {
	return eventType.name
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
	log.Printf("EventSocket: %s Adding client callback to callback list of topic %s", callbackId, callbacks.topic)
	(*callbacks.callbacks)[callbackId] = callback
}

func (callbacks *EventCallbacks) removeCallback(callbackId string) {
	log.Printf("Removing callback %s to %s", callbackId, callbacks.topic)
	delete(*callbacks.callbacks, callbackId)
}

type EventsCallbacksManager struct {
	eventCallbacksList map[string]*EventCallbacks
}

func (manager *EventsCallbacksManager) getEventCallbacks(topic string) *EventCallbacks {
	eventCallbacks, _ := manager.eventCallbacksList[topic]
	return eventCallbacks
}

func (manager *EventsCallbacksManager) addEventCallback(topic string, clientId string, callback Callback) *EventCallbacks {

	callbackId := clientId

	if _, found := manager.eventCallbacksList[topic]; !found {
		log.Printf("EventsSocketId: %s Event callbacks list of topic not found EventsCallbacksManager creating list for topic: %s", callbackId, topic)
		manager.eventCallbacksList[topic] = newEventCallbacks(topic)
	}

	eventCallbacks, ok := manager.eventCallbacksList[topic]

	if !ok {
		log.Printf("EventsSocketId: %s Event callbacks list of topic not found. client not added to topic: %s", callbackId, topic)
	} else {
		log.Printf("EventsSocketId: %s Event callbacks list of topic found. client added to topic: %s", callbackId, topic)
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

func (roomEvents *RoomEvents) SubscribeCheerAdded(room *Room, clientId string, cb Callback) {
	topic := EVENT_CHEER_ADDED.topicName(room)
	eventCallbacks := roomEvents.eventsCallbacksManager.addEventCallback(topic, clientId, cb)

	if !roomEvents.eventBus.HasCallback(topic) {
		log.Printf("subcribing event callbacks to topic %s", topic)
		err := roomEvents.eventBus.Subscribe(topic, eventCallbacks.callbackAll)
		if err != nil {
			//todo more info
			log.Panicf("something wrong when subscribing")
		}
	}
	log.Printf("EventsSocketId[%s]: Subscribed callback to topic[%s]", clientId, topic)
}

func (roomEvents *RoomEvents) PublishCheerAdded(room *Room, cheer cheers.Cheer) {
	topic := EVENT_CHEER_ADDED.topicName(room)
	roomEvents.eventBus.Publish(topic, cheer)
	log.Printf("Publishing event room[%s] | topic %s", room.Name, topic)
}

func (roomEvents *RoomEvents) UnsubscribeCheerAdded(room *Room, clientId string) {
	topic := EVENT_CHEER_ADDED.topicName(room)
	log.Printf("unsubcribing rooms %s to %s", room.Name, topic)
	roomEvents.eventsCallbacksManager.removeEventCallback(topic, clientId)
}
