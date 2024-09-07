import { type Event } from '../../types/match';
import './MatchEvent.css';

interface Props {
  event: Event
  localTeamId: string
}

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

const MatchEventCard: React.FC<Props> = ({ event, localTeamId }) => {
  return (
    <div className={`event ${event.team === localTeamId ? 'left' : 'right'}`}>
      <span className="minute">{event.minute}'</span>

      {event.team === localTeamId && (
        <div className="event-content">
          <span className="description">{eventDescriptions[event.type]} {event.player}</span>
          <img src={eventIcons[event.type]} alt={event.type} className="icon" />
        </div>
      )}


      {event.team !== localTeamId && (
        <div className="event-content">
          <img src={eventIcons[event.type]} alt={event.type} className="icon" />
          <span className="description">{eventDescriptions[event.type]} {event.player}</span>
        </div>
      )}

    </div>
  );
};

export default MatchEventCard;