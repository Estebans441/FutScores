import { type Match } from "../types/match.ts";

export const getMatches = async () : Promise<Match[]> => {
    const response = await fetch('http://localhost:8080/matches');
    const matches = await response.json();
    return matches;
}

export const getMatchById = async (id: number): Promise<Match | undefined> => {
    const matches = await getMatches();
    return matches.find(match => match.id === id);
}