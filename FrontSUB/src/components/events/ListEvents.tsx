import React, { useEffect, useState } from 'react';
import { type Match, type Event } from "../../types/match";
import MatchEventCard from './MatchEventCard';
import MatchEventService from '../../backend/matchEventService';

interface Props {
  match: Match;
}

const ListEvents: React.FC<Props> = ({ match }) => {
  const [events, setEvents] = useState<Event[]>([]);
  const [eventTypes, setEventTypes] = useState<string[]>(["#"]);
  const localTeamId = match.homeTeam;
  const matchEventService = new MatchEventService(eventTypes, match, setEvents);

  useEffect(() => {
    matchEventService.activate();
    matchEventService.client.onConnect = () => {
      console.log('Conectado a RabbitMQ WebSocket');
      for (let eventType of eventTypes) {
        matchEventService.subscribe(eventType, setEvents);
      }
    };
    return () => {
        matchEventService.deactivate();
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