import { type Match } from "../types/match.ts";

const BACKEND_HOST = process.env.BACKEND_HOST || 'localhost';
const RABBITMQ_HOST = process.env.RABBITMQ_HOST || 'localhost';

export const getMatches = async (): Promise<Match[]> => {
    const response = await fetch(`http://${BACKEND_HOST}:8080/matches`);
    const matches = await response.json();
    return matches;
}

export const getMatchById = async (id: number): Promise<Match | undefined> => {
    const matches = await getMatches();
    return matches.find(match => match.id === id);
}

export const createMatch = async (match: Match): Promise<void> => {
    await fetch(`http://${BACKEND_HOST}:8080/matches`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(match),
    });
}

export const getBackendHost = (): string => {
    return BACKEND_HOST;
}

export const getRabbitMQHost = (): string => {
    return RABBITMQ_HOST;
}


export const getTeams = () : {teamName: string, teamAbbr: string}[] => {
    return [
        { teamName: "FC Barcelona", teamAbbr: "FCB" },
        { teamName: "Real Madrid", teamAbbr: "RMD" },
        { teamName: "Athletic Bilbao", teamAbbr: "ATH" },
        { teamName: "Atlético de Madrid", teamAbbr: "ATM" },
        { teamName: "CA Osasuna", teamAbbr: "OSA" },
        { teamName: "CD Leganés", teamAbbr: "LEG" },
        { teamName: "Celta de Vigo", teamAbbr: "CEL" },
        { teamName: "Deportivo Alavés", teamAbbr: "ALA" },
        { teamName: "Getafe CF", teamAbbr: "GET" },
        { teamName: "Girona FC", teamAbbr: "GIR" },
        { teamName: "Rayo Vallecano", teamAbbr: "RAY" },
        { teamName: "RCD Espanyol Barcelona", teamAbbr: "ESP" },
        { teamName: "RCD Mallorca", teamAbbr: "MAL" },
        { teamName: "Real Betis Balompié", teamAbbr: "BET" },
        { teamName: "Real Sociedad", teamAbbr: "SOC" },
        { teamName: "Real Valladolid CF", teamAbbr: "VAD" },
        { teamName: "Sevilla FC", teamAbbr: "SEV" },
        { teamName: "UD Las Palmas", teamAbbr: "LPA" },
        { teamName: "Valencia CF", teamAbbr: "VAL" },
        { teamName: "Villarreal CF", teamAbbr: "VIL" },
      ];
}