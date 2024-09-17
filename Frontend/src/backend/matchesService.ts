import { type Match } from "../types/match.ts";

const BACKEND_HOST = process.env.BACKEND_HOST || 'localhost';
const RABBITMQ_HOST = process.env.RABBITMQ_HOST || 'localhost';

export const getMatches = async (): Promise<Match[]> => {
    const response = await fetch(`http://${BACKEND_HOST}:8081/matches`);
    const matches = await response.json();
    return matches;
}

export const getMatchById = async (id: number): Promise<Match | undefined> => {
    const matches = await getMatches();
    return matches.find(match => match.id === id);
}

export const getRabbitMQHost = (): string => {
    return RABBITMQ_HOST;
}