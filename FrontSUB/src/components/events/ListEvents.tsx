import React, { useEffect, useState } from 'react';
import { Client } from '@stomp/stompjs';
import { type Match, type Event } from "../../types/match";
import MatchEventCard from './MatchEventCard';

interface Props {
  match: Match;
}

const ListEvents: React.FC<Props> = ({ match }) => {
  const [events, setEvents] = useState<Event[]>([]);
  const [eventTypes, setEventTypes] = useState<string[]>(["goal", "penalty"]);
  const localTeamId = match.homeTeam;

  useEffect(() => {
    const client = new Client({
        brokerURL: 'ws://localhost:15674/ws', // URL del WebSocket de RabbitMQ
        connectHeaders: {
            login: 'guest',
            passcode: 'guest',
        },
        debug: (str) => {
            console.log(str);
        },
        onConnect: () => {
          for (let eventType of eventTypes) {
            console.log('Conectado a RabbitMQ WebSocket');
            client.subscribe(`/exchange/match_events/match.${match.id}.event.${eventType}`, (message) => {
              console.log('Mensaje recibido:', message.body);
              setEvents((prevEvents) => [...prevEvents, JSON.parse(message.body)]);
            });
          }
        },
        onStompError: (frame) => {
            console.error('Error de STOMP:', frame.headers['message']);
        },
    });

    client.activate();

    return () => {
        client.deactivate();
    };
}, []);

  return (
    <div className="flex flex-col items-center relative w-100 bg-white px-20">
      {events.map((event) => (
        <MatchEventCard
          key={event.id}
          event={event}
          localTeamId= {localTeamId}
        />
      ))}
    </div>
  );
};



export default ListEvents;