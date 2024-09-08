
import { Client } from '@stomp/stompjs';
import { type Match, type Event } from "../types/match";

export default class MatchEventService {
  client: Client;
  match: Match;

  constructor(eventTypes: string[], match: Match, setEvents: (events: (prevEvents: Event[]) => Event[]) => void) {
    this.match = match;  
    this.client = new Client({
        brokerURL: 'ws://localhost:15674/ws', // URL del WebSocket de RabbitMQ
        connectHeaders: {
            login: 'guest',
            passcode: 'guest',
        },
        debug: (str) => {
            console.log(str);
        },
        onConnect: () => {
            console.log('Conectado a RabbitMQ WebSocket');
            for (let eventType of eventTypes) {
                this.subscribe(eventType, setEvents);
            }
        },
        onStompError: (frame) => {
            console.error('Error de STOMP:', frame.headers['message']);
        },
    });
  }

  subscribe(eventType: string, setEvents: (events: (prevEvents: Event[]) => Event[]) => void) {
    this.client.subscribe(`/exchange/match_events/match.${this.match.id}.event.${eventType}`, (message) => {
      console.log('Mensaje recibido:', message.body);
      setEvents((prevEvents: Event[]) => [...prevEvents, JSON.parse(message.body)]);
    });
  }

  subscribeToNewTopic(eventType: string, setEvents: (events: (prevEvents: Event[]) => Event[]) => void) {
    if (this.client.connected) {
      this.subscribe(eventType, setEvents);
    } else {
      console.error('No se puede suscribir, el cliente no está conectado.');
    }
  }

  unsubscribe(eventType: string) {
    this.client.unsubscribe(`/exchange/match_events/match.${this.match.id}.event.${eventType}`);
  }

  activate() {
    this.client.activate();
  }

  deactivate() {
    this.client.deactivate();
  }
}