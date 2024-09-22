
import { Client } from '@stomp/stompjs';
import { type Match, type Event } from "../types/match";

export default class MatchEventService {
  client: Client;
  match: Match;

  constructor(eventTypes: string[], match: Match, RABBITMQ_HOST: string, setEvents: (events: (prevEvents: Event[]) => Event[]) => void) {
    this.match = match;
    this.client = new Client({
      brokerURL: `ws://${RABBITMQ_HOST}:15674/ws`, // URL del WebSocket de RabbitMQ
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

export const getEventTypes = (): {eventType: string, display: string}[] => {
  return [
    { eventType: 'goal', display: 'Gol' },
    { eventType: 'penalty', display: 'Penal' },
    { eventType: 'red card', display: 'Tarjeta Roja' },
    { eventType: 'yellow card', display: 'Tarjeta Amarilla' },
    { eventType: 'substitution', display: 'Sustitución' },
    { eventType: 'offside', display: 'Fuera de Juego' },
    { eventType: 'corner kick', display: 'Tiro de Esquina' },
    { eventType: 'free kick', display: 'Tiro Libre' },
    { eventType: 'start', display: 'Inicio' },
    { eventType: 'half-time', display: 'Medio Tiempo' },
    { eventType: 'end', display: 'Final' },
  ];
}

export const createEvent = async (event: Event, BACKEND_HOST:string): Promise<void> => {
  await fetch(`http://${BACKEND_HOST}:8080/events`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(event),
  });
};