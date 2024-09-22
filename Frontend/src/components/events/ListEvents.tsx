import React, { useEffect, useState } from 'react';
import { type Match, type Event } from "../../types/match";
import MatchEventCard from './MatchEventCard';
import MatchEventService from '../../backend/matchEventService';

interface Props {
  match: Match;
  RABBITMQ_HOST: string;
}

const ListEvents: React.FC<Props> = ({ match, RABBITMQ_HOST }) => {
  // State to store the list of events to be displayed
  const [events, setEvents] = useState<Event[]>([]);
  const localTeamId = match.homeTeam;

  // Create an instance of the event service
  const matchEventService = new MatchEventService(["#"], match, RABBITMQ_HOST, setEvents); // '#' is a wildcard to receive all events

  // Hook to fetch initial events from the REST API
  useEffect(() => {
    const fetchEvents = async () => {
      try {
        const response = await fetch(`http://localhost:8080/matches/${match.id}/events`);
        if (!response.ok) {
          throw new Error('Error fetching events');
        }
        const fetchedEvents: Event[] = await response.json();

        // Sort the events by minute before storing them in state
        const sortedEvents = fetchedEvents.sort((a, b) => a.minute - b.minute);
        setEvents(sortedEvents);
      } catch (error) {
        console.error('Failed to fetch events:', error);
      }
    };
    fetchEvents();
  }, [match.id]); // Dependency array ensures the hook runs when match.id changes

  // Hook to handle subscription to RabbitMQ
  useEffect(() => {
    matchEventService.activate(); // Activate subscription when the component mounts
    return () => {
      matchEventService.deactivate(); // Deactivate subscription when the component unmounts
    };
  }, [matchEventService]); // Dependency array ensures the hook runs only when the matchEventService changes

  return (
    <>
      <h2> Events </h2>
      <div className="flex flex-col items-center relative w-100 bg-white px-20">
        {events.map((event) => (
          <MatchEventCard
            key={event.id}
            event={event}
            localTeamId={localTeamId}
          />
        ))}
      </div>
    </>
  );
};

export default ListEvents;