import { type Event } from "../types/match";

export const getEvents = async () => {
    const eventos: Event[] = [
        { id: 1, matchId: 1, team: "FC Barcelona", player: "Gavi", type: "goal", minute: 15 },
        { id: 2, matchId: 1, team: "Real Madrid", player: "Bellingham", type: "goal", minute: 30 },
        { id: 3, matchId: 1, team: "FC Barcelona", player: "Lewandowski", type: "penalty", minute: 35 },
        { id: 4, matchId: 1, team: "Real Madrid", player: "Vinicius Jr.", type: "red card", minute: 40 },
        { id: 5, matchId: 1, team: "FC Barcelona", player: "Gundogan", type: "yellow card", minute: 50 },
        { id: 6, matchId: 1, team: "Real Madrid", player: "Joselu", type: "substitution", minute: 60 },
        { id: 7, matchId: 1, team: "FC Barcelona", player: "Araujo", type: "offside", minute: 65 },
        { id: 8, matchId: 1, team: "Real Madrid", player: "Modric", type: "corner kick", minute: 75 },
        { id: 9, matchId: 1, team: "FC Barcelona", player: "Ferran Torres", type: "free kick", minute: 80 },
        { id: 10, matchId: 1, team: "Real Madrid", player: "Bellingham", type: "goal", minute: 85 },
        { id: 11, matchId: 1, team: "FC Barcelona", player: "Ter Stegen", type: "start", minute: 0 },
        { id: 12, matchId: 1, team: "Real Madrid", player: "Courtois", type: "half-time", minute: 45 },
        { id: 13, matchId: 1, team: "Real Madrid", player: "Modric", type: "end", minute: 90 }
    ];
    // Ordenar eventos por minuto
    eventos.sort((a, b) => a.minute - b.minute);
    return eventos;
};
