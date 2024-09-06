import { type Event } from "../types/match"

export const getEvents = async () => {
    const eventos: Event[] = [
        { id: 1, matchId: 1, team: "FC Barcelona", player: "Messi", type: "goal", minute: 10 },
        { id: 2, matchId: 1, team: "Real Madrid", player: "Ronaldo", type: "goal", minute: 20 },
        { id: 3, matchId: 1, team: "FC Barcelona", player: "Messi", type: "penalty", minute: 30 },
        { id: 4, matchId: 1, team: "Real Madrid", player: "Ronaldo", type: "red card", minute: 40 },
        { id: 5, matchId: 1, team: "FC Barcelona", player: "Messi", type: "yellow card", minute: 50 },
        { id: 6, matchId: 1, team: "Real Madrid", player: "Ronaldo", type: "substitution", minute: 60 },
        { id: 7, matchId: 1, team: "FC Barcelona", player: "Messi", type: "offside", minute: 70 },
        { id: 8, matchId: 1, team: "Real Madrid", player: "Ronaldo", type: "corner kick", minute: 80 },
        { id: 9, matchId: 1, team: "FC Barcelona", player: "Messi", type: "free kick", minute: 90 },
        { id: 10, matchId: 1, team: "Real Madrid", player: "Ronaldo", type: "goal", minute: 100 },
        { id: 11, matchId: 1, team: "FC Barcelona", player: "Messi", type: "start", minute: 0 },
        { id: 12, matchId: 1, team: "Real Madrid", player: "Ronaldo", type: "half-time", minute: 45 },
        { id: 13, matchId: 1, team: "FC Barcelona", player: "Messi", type: "end", minute: 90 }
    ]
    // ordenar eventos por minuto
    eventos.sort((a, b) => a.minute - b.minute)
    return eventos
}