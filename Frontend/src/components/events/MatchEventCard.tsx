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
  "goal": 'Gol',
  "penalty": 'Penalti',
  "red card": 'Tarjeta roja',
  "yellow card": 'Tarjeta amarilla',
  "substitution": 'Cambio',
  "offside": 'Fuera de juego',
  "corner kick": 'Tiro de esquina',
  "free kick": 'Tiro libre',
  "start": 'Inicio del tiempo',
  "half-time": 'Medio tiempo',
  "end": 'Fin del tiempo'
};

const MatchEventCard: React.FC<Props> = ({ event, localTeamId }) => {
  const isLocalTeam = event.team === localTeamId;

  return (
    <div className="event-row">
      {/* Div para el equipo local */}
      <div className={`event-content local-team ${isLocalTeam ? 'active' : ''}`}>
      {isLocalTeam && (
  <>
    <div className="event-description">
      <span className="player">{event.player}</span>
      <span className="event-type">{eventDescriptions[event.type]}</span>
    </div>
    <img src={eventIcons[event.type]} alt={event.type} className="icon" />
  </>
)}
      </div>

      {/* Div para el minuto */}
      <div className="event-minute">
        <span className="minute">{event.minute}'</span>
      </div>

      {/* Div para el equipo visitante */}
      <div className={`event-content visitor-team ${!isLocalTeam ? 'active' : ''}`}>
      {!isLocalTeam && (
  <>
    <img src={eventIcons[event.type]} alt={event.type} className="icon" />
    <div className="event-description">
      <span className="player">{event.player}</span>
      <span className="event-type">{eventDescriptions[event.type]}</span>
    </div>
  </>
)}
      </div>
    </div>
  );
};

export default MatchEventCard;
