import React, { useEffect, useState } from 'react';
import { type Match, type Event } from "../../types/match";
import MatchEventCard from './MatchEventCard';
import MatchEventService from '../../backend/matchEventService';

interface Props {
  match: Match;
}

const ListEvents: React.FC<Props> = ({ match }) => {
  // List of events to be displayed
  const [events, setEvents] = useState<Event[]>([]);
  const localTeamId = match.homeTeam;
  
  // MatchEventService instance. It will handle the connection to RabbitMQ
  const matchEventService = new MatchEventService(["#"], match, setEvents); // # is a wildcard to receive all events

  useEffect(() => {
    matchEventService.activate();
    return () => {
        matchEventService.deactivate();
    };
  }, []);

  return (
    <>
      <h2> Eventos  </h2>
      <div className="flex flex-col items-center relative w-100 bg-white px-20">
        {events.map((event) => (
          <MatchEventCard
          key={event.id}
          event={event}
          localTeamId= {localTeamId}
          />
        ))}
      </div>
    </>
  );
};



export default ListEvents;