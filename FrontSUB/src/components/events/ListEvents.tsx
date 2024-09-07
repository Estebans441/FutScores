import React, { useEffect, useState } from 'react';
import { fetchEvents } from "../../backend/eventsService";
import { type Match, type Event } from "../../types/match";
import MatchEventCard from './MatchEventCard';

interface Props {
  match: Match;
}

const ListEvents: React.FC<Props> = ({ match }) => {
  const [events, setEvents] = useState<Event[]>([]);
  const localTeamId = match.homeTeam;

  useEffect(() => {
    fetchEvents((evento) => {
        console.log(evento);
        setEvents((prevEvents) => [...prevEvents, evento]);
    });
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