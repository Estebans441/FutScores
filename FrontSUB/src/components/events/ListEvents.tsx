import React, { useEffect, useState } from 'react';
import { fetchEvents } from "../../backend/eventsService";
import { type Match, type Event } from "../../types/match";
import './ListEvents.css';

interface Props {
  match: Match;
}

const ListEvents: React.FC<Props> = ({ match }) => {
  const [events, setEvents] = useState<Event[]>([]);
  const localTeamId = match.homeTeam;


  const eventIcons: { [key: string]: string } = {
    "goal": '/event_icons/gol.png',
    "penalty": '/event_icons/penal.png',
    "red card": '/event_icons/roja.png',
    "yellow card": '/event_icons/amarilla.png',
    "substitution": '/event_icons/cambio.png',
    "offside": '/event_icons/fuera.png',
    "corner kick": '/event_icons/esquina.png',
    "free kick": '/event_icons/libre.png',
    "start": '/event_icons/tiempo.png',
    "half-time": '/event_icons/tiempo.png',
    "end": '/event_icons/tiempo.png'
  };

  const eventDescriptions: { [key: string]: string } = {
    "goal": 'Gol de',
    "penalty": 'Penalti de',
    "red card": 'Tarjeta roja para',
    "yellow card": 'Tarjeta amarilla para',
    "substitution": 'Cambio de',
    "offside": 'Fuera de juego de',
    "corner kick": 'Tiro de esquina de',
    "free kick": 'Tiro libre de',
    "start": 'Inicio del tiempo',
    "half-time": 'Medio tiempo',
    "end": 'Fin del tiempo'
  };

  useEffect(() => {
    fetchEvents((evento) => {
        console.log(evento);
        setEvents((prevEvents) => [...prevEvents, evento]);
    });
  }, []);

  return (
    <div className="timeline">
      {events.map((event) => (
        <div className={`event ${event.team === localTeamId ? 'left' : 'right'}`} key={event.id}>
          {event.team === localTeamId && (
            <div className="event-content">
              <span className="description">{eventDescriptions[event.type]} {event.player}</span>
              <img src={eventIcons[event.type]} alt={event.type} className="icon" />
            </div>
          )}
          <span className="minute">{event.minute}'</span>
          {event.team !== localTeamId && (
            <div className="event-content">
              <img src={eventIcons[event.type]} alt={event.type} className="icon" />
              <span className="description">{eventDescriptions[event.type]} {event.player}</span>
            </div>
          )}
        </div>
      ))}
    </div>
  );
};

export default ListEvents;